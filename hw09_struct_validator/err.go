package hw09structvalidator

import (
	"errors"
	"fmt"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrNotStruct         = errors.New("input is not struct")
	ErrInvalidRuleParam  = errors.New("invalid rule param")
	ErrValidation        = errors.New("validation error")
	ErrUnknownIntRule    = errors.New("unknown int rule")
	ErrUnknownStringRule = errors.New("unknown string rule")
)

func (v ValidationErrors) Error() string {
	var res strings.Builder
	for _, err := range v {
		res.WriteString(fmt.Sprintf("field \"%v\" error: \"%v\"", err.Field, err.Err))
	}
	return res.String()
}

func (v ValidationErrors) Is(tgt error) bool {
	for _, err := range v {
		if errors.Is(err.Err, tgt) {
			return true
		}
	}
	return false
}
