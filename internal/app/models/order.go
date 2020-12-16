package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

// Since we need "good-till-cancelled" order type for a task accomplishment
// let's define only time limited orders
type TimeLimitedOrderType string

// Enum for time limited orders
const (
	OneDay            TimeLimitedOrderType = "one-day"
	GoodTillCancelled TimeLimitedOrderType = "good-till-cancelled"
	ImmediateOrCancel TimeLimitedOrderType = "immediate-or-cancel"
	FillOrKill        TimeLimitedOrderType = "fill-or-kill"
)

type OrderGeneralInfo struct {
	ValidUntil   *time.Time
	Price        decimal.Decimal
	Quantity     uint
	Operation    MarketOperation
	CounterParty string
}

// We also assume that there will be only sell/buy operations with securities on the market for the sake of simplicity
type TimeLimitedOrder struct {
	OrderGeneralInfo
	TradeCode uuid.UUID
	Type      TimeLimitedOrderType
}

func NewGoodTillCancelledOrder(info OrderGeneralInfo) *TimeLimitedOrder {
	if info.ValidUntil == nil {
		defaultExpireTime := time.Now().UTC().Add(time.Hour * 24 * 90)
		info.ValidUntil = &defaultExpireTime
	}

	return &TimeLimitedOrder{
		OrderGeneralInfo: info,
		TradeCode:        uuid.New(),
		Type:             GoodTillCancelled,
	}
}
