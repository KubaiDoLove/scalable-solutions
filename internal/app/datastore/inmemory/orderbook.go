package inmemory

import (
	"context"
	"github.com/KubaiDoLove/scalable-solutions/internal/app/datastore"
	"github.com/KubaiDoLove/scalable-solutions/internal/app/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"sort"
	"sync"
)

// Thread safe "in memory" store
type OrderBook struct {
	mu   sync.Mutex
	asks map[uuid.UUID]models.Order
	bids map[uuid.UUID]models.Order
}

func New() datastore.DataStore {
	return &OrderBook{
		asks: make(map[uuid.UUID]models.Order, 0),
		bids: make(map[uuid.UUID]models.Order, 0),
	}
}

// CreateOrder validates order on a very basic level and saves it
func (o *OrderBook) CreateOrder(ctx context.Context, order *models.Order) error {
	if order == nil {
		return datastore.ErrEmptyStruct
	}

	if order.ID == uuid.Nil {
		return datastore.ErrZeroID
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	if order.Operation == models.Ask {
		o.asks[order.ID] = *order
		return nil
	}

	o.bids[order.ID] = *order
	return nil
}

// DisableOrder to remove it from market snapshot and other reads
func (o *OrderBook) DisableOrder(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return datastore.ErrZeroID
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	if order, orderInAsks := o.asks[id]; orderInAsks {
		if order.IsProcessable() {
			o.asks[id].OrderGeneralInfo.IsEnabled = false
		}
		return nil
	}

	if order, orderInBids := o.bids[id]; orderInBids {
		if order.IsProcessable() {
			o.bids[id].OrderGeneralInfo.IsEnabled = false
		}
		return nil
	}

	return datastore.ErrOrderDoesNotExist
}

// OrderByID returns only enabled and not expired order
func (o *OrderBook) OrderByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	if id == uuid.Nil {
		return nil, datastore.ErrZeroID
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	if order, orderInAsks := o.asks[id]; orderInAsks {
		if order.IsProcessable() {
			return &order, nil
		}
	}

	if order, orderInBids := o.bids[id]; orderInBids {
		if order.IsProcessable() {
			return &order, nil
		}
	}

	return nil, datastore.ErrOrderDoesNotExist
}

// MatchOrder to get available bids/asks for a given order
func (o *OrderBook) MatchOrder(ctx context.Context, order *models.Order) ([]models.Order, error) {
	if order == nil {
		return nil, datastore.ErrEmptyStruct
	}

	if order.Operation == models.Bid {
		return o.matchBid(order.Price)
	}

	return o.matchAsk(order.Price)
}

func (o *OrderBook) matchBid(bidPrice decimal.Decimal) ([]models.Order, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	matchingAsks := make([]models.Order, 0)

	for _, ask := range o.asks {
		if ask.IsProcessable() && ask.Price.LessThanOrEqual(bidPrice) && ask.Quantity > 0 {
			matchingAsks = append(matchingAsks, ask)
		}
	}

	return matchingAsks, nil
}

func (o *OrderBook) matchAsk(askPrice decimal.Decimal) ([]models.Order, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	matchingBids := make([]models.Order, 0)

	for _, bid := range o.bids {
		if bid.IsProcessable() && bid.Price.GreaterThanOrEqual(askPrice) && bid.Quantity > 0 {
			matchingBids = append(matchingBids, bid)
		}
	}

	return matchingBids, nil
}

// MarketDataSnapshot to get actual market data ordered by price
func (o *OrderBook) MarketDataSnapshot(ctx context.Context) (*models.MarketDataSnapshot, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	asks := make([]models.OrderSnapshot, 0)
	for _, ask := range o.asks {
		if ask.IsProcessable() {
			asks = append(asks, *ask.Snapshot())
		}
	}

	bids := make([]models.OrderSnapshot, 0)
	for _, bid := range o.bids {
		if bid.IsProcessable() {
			bids = append(bids, *bid.Snapshot())
		}
	}

	sort.SliceStable(asks, func(i, j int) bool { return asks[i].Price.LessThan(asks[j].Price) })
	sort.SliceStable(bids, func(i, j int) bool { return bids[i].Price.LessThan(bids[j].Price) })

	return &models.MarketDataSnapshot{
		Asks: asks,
		Bids: bids,
	}, nil
}
