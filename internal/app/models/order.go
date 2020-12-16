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
	GoodTillCancelled                      = "good-till-cancelled"
	ImmediateOrCancel                      = "immediate-or-cancel"
	FillOrKill                             = "fill-or-kill"
)

// We also assume that there will be only sell/buy operations with securities on the market for the sake of simplicity
type TimeLimitedOrder struct {
	TradeCode    uuid.UUID
	ValidUntil   time.Time
	Type         TimeLimitedOrderType
	Price        decimal.Decimal
	Quantity     uint
	Operation    Operation
	CounterParty string
}
