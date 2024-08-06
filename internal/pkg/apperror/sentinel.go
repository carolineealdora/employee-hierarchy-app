package apperror

import (
	"errors"
)

var (
	ErrValidator = errors.New("invalid request. please fill all field(s) with the correct data")
)
