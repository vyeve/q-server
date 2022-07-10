package models

import (
	"errors"
)

// common errors
var (
	ErrIncorrectAmount = errors.New("incorrect amount format")
	ErrUnsupportedFiat = errors.New("unsupported fiat")
)
