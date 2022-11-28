package hw09structvalidator

import (
	"fmt"
	"reflect"
)

const validateTag = "validate"

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("%w. Input is %v", ErrNotStruct, val.Kind())
	}

	var errList ValidationErrors
	fields := reflect.VisibleFields(val.Type())

	for _, field := range fields {
		rVal := val.FieldByName(field.Name)
		err := validateField(field, rVal)
		if err != nil {
			errList = append(errList, ValidationError{Field: field.Name, Err: err})
		}
	}
	if len(errList) == 0 {
		return nil
	}
	return errList
}

func validateField(field reflect.StructField, rVal reflect.Value) error {
	tag := field.Tag.Get(validateTag)

	if len(tag) == 0 {
		return nil
	}

	ruleList := getRuleListFromTag(tag)

	switch field.Type.Kind() { //nolint:exhaustive
	case reflect.String:
		val := rVal.String()
		err := validateString(val, ruleList)
		if err != nil {
			return err
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := rVal.Int()
		err := validateInt(val, ruleList)
		if err != nil {
			return err
		}
	case reflect.Slice:
		switch rVal.Type().Elem().Kind() { //nolint:exhaustive
		case reflect.String:
			vals := rVal.Interface().([]string)
			err := validateSliceString(vals, ruleList)
			if err != nil {
				return err
			}
		case reflect.Int:
			vals := rVal.Interface().([]int64)
			err := validateSliceInt(vals, ruleList)
			if err != nil {
				return err
			}
		default:
			return nil
		}
	default:
		return nil
	}
	return nil
}
