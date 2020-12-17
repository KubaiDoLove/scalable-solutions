package clickhouse

import (
	"context"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/KubaiDoLove/scalable-solutions/internal/app/datastore"
	"github.com/KubaiDoLove/scalable-solutions/internal/app/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type OrderBook struct {
	db *sqlx.DB
}

// Usually you want to get db connection from the function parameter,
// but for our example we will establish it right here
func New() (datastore.DataStore, error) {
	db, err := sqlx.Open("clickhouse", "tcp://127.0.0.1:9000?compress=true&debug=true")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS orders (
        	id UUID,
        	tradeCode UUID,
        	validUntil Date,
        	price String,
        	quantity UInt32,
        	operation String,
        	counterParty String,
        	isEnabled UInt8,
        	type String
        ) engine=Memory
    `)
	if err != nil {
		return nil, err
	}

	return &OrderBook{
		db: db,
	}, nil
}

func (o *OrderBook) Close() error {
	return o.db.Close()
}

func (o *OrderBook) CreateOrder(ctx context.Context, order *models.Order) error {
	if order == nil {
		return datastore.ErrEmptyStruct
	}

	if order.ID == uuid.Nil {
		return datastore.ErrZeroID
	}

	tx, err := o.db.BeginTx(ctx, nil)
	stmt, err := tx.Prepare(
		`INSERT INTO orders 
					(id, tradeCode, validUntil, price, quantity, operation, counterParty, isEnabled, type)
					VALUES
					(?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}

	isEnabled := uint8(0)
	if order.IsEnabled {
		isEnabled = uint8(1)
	}

	if _, err := stmt.ExecContext(
		ctx,
		order.ID,
		order.TradeCode,
		order.ValidUntil,
		order.Price.String(),
		uint32(order.Quantity),
		order.Operation,
		order.CounterParty,
		isEnabled,
		order.Type,
	); err != nil {
		return err
	}

	return tx.Commit()
}

func (o *OrderBook) DisableOrder(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return datastore.ErrZeroID
	}

	stmt, err := o.db.Prepare(`ALTER TABLE orders UPDATE isEnabled = 0 WHERE id = ?`)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(ctx, id); err != nil {
		return err
	}

	return nil
}

func (o *OrderBook) OrderByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	if id == uuid.Nil {
		return nil, datastore.ErrZeroID
	}

	stmt, err := o.db.Preparex(`SELECT * FROM orders WHERE id = ?`)
	if err != nil {
		return nil, err
	}

	order := new(models.Order)
	if err := stmt.Get(order, id); err != nil {
		return nil, err
	}

	if !order.IsProcessable() {
		return nil, datastore.ErrOrderDoesNotExist
	}

	return order, nil
}

func (o *OrderBook) MatchOrder(ctx context.Context, order *models.Order) ([]models.Order, error) {
	if order == nil {
		return nil, datastore.ErrEmptyStruct
	}

	return o.matchOrderByOperation(order.Operation, order.Price)
}

func (o *OrderBook) matchOrderByOperation(operation models.MarketOperation, price decimal.Decimal) ([]models.Order, error) {
	query := `
		SELECT *
		FROM orders
		WHERE 
			operation = 'ask'
				AND isEnabled = 1
				AND now('Europe/London') <= validUntil
				AND toFloat64(price) <= toFloat64(?)
				AND quantity > 0
	`

	if operation == models.Ask {
		query = `
			SELECT *
			FROM orders
			WHERE operation = 'bid'
				AND isEnabled = 1
				AND now('Europe/London') <= validUntil
				AND toFloat64(price) >= toFloat64(?)
				AND quantity > 0
		`
	}

	matchingOrders := make([]models.Order, 0)

	stmt, err := o.db.Preparex(query)
	if err != nil {
		return nil, err
	}

	err = stmt.Select(&matchingOrders, price.String())
	if err != nil {
		return nil, err
	}

	return matchingOrders, nil
}

func (o *OrderBook) MarketDataSnapshot(ctx context.Context) (*models.MarketDataSnapshot, error) {
	stmt, err := o.db.Preparex(`
		SELECT price, quantity
		FROM orders
		WHERE operation = ? AND isEnabled = 1 AND now('Europe/London') <= validUntil
		ORDER BY toFloat64(price)
	`)
	if err != nil {
		return nil, err
	}

	asks := make([]models.OrderSnapshot, 0)
	err = stmt.Select(&asks, "ask")
	if err != nil {
		return nil, err
	}

	bids := make([]models.OrderSnapshot, 0)
	err = stmt.Select(&bids, "bid")
	if err != nil {
		return nil, err
	}

	return &models.MarketDataSnapshot{
		Asks: asks,
		Bids: bids,
	}, nil
}
