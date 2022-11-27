package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrorNotStruct = errors.New("input is not struct")
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
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Struct {
		return ErrorNotStruct
	}

	// Place your code here.
	return nil
}
