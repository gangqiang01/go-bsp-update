package global

import (
	"errors"
)

var (
	ErrInvalidParms = errors.New("Invalid parms")
	ErrUnknown      = errors.New("unknown error.")
)
