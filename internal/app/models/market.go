package models

type MarketDataSnapshot struct {
	Asks []OrderSnapshot
	Bids []OrderSnapshot
}
