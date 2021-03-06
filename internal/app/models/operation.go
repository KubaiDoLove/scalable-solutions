package models

type MarketOperation string

// We assume that we will have more than 2 operations in future,
// so we define an enum instead of IsSelling or IsBuying flag
const (
	Ask MarketOperation = "ask"
	Bid MarketOperation = "bid"
)
