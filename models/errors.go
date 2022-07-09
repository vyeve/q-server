package models

import (
	"errors"
)

var (
	ErrIncorrectAmount = errors.New("incorrect amount format")
	ErrUnsupportedFiat = errors.New("unsupported fiat")
)
