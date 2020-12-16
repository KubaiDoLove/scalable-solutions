package models

import (
	"errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

var (
	ErrNoEmptyGeneralInfo = errors.New("need general info to process further")
)

// Since we need "good-till-cancelled" order type for a task accomplishment
// let's define only time limited orders
type TimeLimitedOrderType string

// Enum for time limited orders only
const (
	OneDay            TimeLimitedOrderType = "one-day"
	GoodTillCancelled TimeLimitedOrderType = "good-till-cancelled"
	ImmediateOrCancel TimeLimitedOrderType = "immediate-or-cancel"
	FillOrKill        TimeLimitedOrderType = "fill-or-kill"
)

// OrderGeneralInfo consists "must have" data for any order
type OrderGeneralInfo struct {
	ID           uuid.UUID
	TradeCode    uuid.UUID
	ValidUntil   *time.Time
	Price        decimal.Decimal
	Quantity     uint
	Operation    MarketOperation
	CounterParty string
	// We never delete orders, only disable
	IsEnabled bool
}

// OrderSnapshot for market data snapshots
type OrderSnapshot struct {
	Price    decimal.Decimal
	Quantity uint
}

// We also assume that there will be only sell/buy operations with securities on the market for the sake of simplicity
type Order struct {
	*OrderGeneralInfo
	Type TimeLimitedOrderType
}

func (o Order) IsProcessable() bool {
	isNotExpired := o.ValidUntil != nil && time.Now().UTC().Before(*o.ValidUntil)
	return o.IsEnabled && isNotExpired
}

func (o Order) Snapshot() *OrderSnapshot {
	return &OrderSnapshot{
		Price:    o.Price,
		Quantity: o.Quantity,
	}
}

// NewGoodTillCancelledOrder creates
func NewGoodTillCancelledOrder(info *OrderGeneralInfo) (*Order, error) {
	if info == nil {
		return nil, ErrNoEmptyGeneralInfo
	}

	if info.ValidUntil == nil {
		defaultExpireTime := time.Now().UTC().Add(time.Hour * 24 * 90)
		info.ValidUntil = &defaultExpireTime
	}

	info.ID = uuid.New()
	info.IsEnabled = true
	return &Order{
		OrderGeneralInfo: info,
		Type:             GoodTillCancelled,
	}, nil
}
