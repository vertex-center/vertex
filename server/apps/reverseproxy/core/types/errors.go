package types

import (
	"errors"
)

var (
	ErrSourceInvalid = errors.New("source is empty")
	ErrTargetInvalid = errors.New("target is empty")
)
