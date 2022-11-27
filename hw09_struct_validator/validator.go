package hw09structvalidator

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var res strings.Builder
	for _, err := range v {
		res.WriteString(fmt.Sprintf("field \"%v\" error: \"%v\"", err.Field, err.Err))
	}
	return res.String()
}

func Validate(v interface{}) error {
	// Place your code here.
	return nil
}
