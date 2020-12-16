package datastore

import "errors"

var (
	ErrEmptyStruct       = errors.New("no empty struct")
	ErrZeroID            = errors.New("no zero id")
	ErrOrderDoesNotExist = errors.New("order does not exist")
)
