package models

type Operation string

// We assume that we will have more than 2 operations in future,
// so we define an enum instead of IsSelling or IsBuying flag
const (
	Sell Operation = "sell"
	Buy            = "buy"
)
