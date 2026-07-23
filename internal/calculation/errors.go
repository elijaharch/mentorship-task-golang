package calculation

import "errors"

var (
	ErrNotFound         = errors.New("calculation not found")
	ErrInvalidOperation = errors.New("invalid operation")
	ErrDivisionByZero   = errors.New("division by zero")
	ErrInvalidNumber    = errors.New("invalid number")
)
