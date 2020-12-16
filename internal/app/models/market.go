package models

// MarketDataSnapshot to get actual market data
type MarketDataSnapshot struct {
	Asks []OrderSnapshot
	Bids []OrderSnapshot
}
