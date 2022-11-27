package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	testsNegative := []struct {
		in          interface{}
		expectedErr error
		contains    []string
	}{
		{in: 5, expectedErr: ErrNotStruct, contains: []string{}}, // not struct
		{in: App{ // validation error
			Version: "123456",
		}, expectedErr: ErrValidation, contains: []string{"field \"Version\" error: \"validation error: length must be 5\""}},
		{
			in: Response{ // validation error
				Code: 504,
				Body: "any",
			},
			expectedErr: ErrValidation,
			contains: []string{
				"field \"Code\" error: \"validation error: does not match any value in list 200,404,500\"",
			},
		},
		{in: User{ // validation errors
			ID:     "wrong id",
			Name:   "any name",
			Age:    95,
			Email:  "test@@test",
			Role:   "user",
			Phones: []string{"9990009900", "495"},
			meta:   json.RawMessage{},
		}, expectedErr: ErrValidation, contains: []string{
			"field \"ID\" error: \"validation error: length must be 36\"",
			"field \"Age\" error: \"validation error: must be not more than 50\"",
			"field \"Email\" error: \"validation error: must match regexp",
			"field \"Role\" error: \"validation error: does not match any value in list admin,stuff\"",
			"field \"Phones\" error: \"validation error: length must be 11\"",
		}},
		{in: struct { // invalid rule param
			N int `validate:"min:a"`
		}{N: 5}, expectedErr: ErrInvalidRuleParam, contains: []string{"can not convert a to int"}},
	}

	testsPositive := []struct {
		in          interface{}
		expectedErr error
	}{
		{in: struct{}{}}, // emty struct
		{
			in: struct { // fields withou validate tag
				A string
				B int
				C []string
				D []int
			}{A: "test", B: 1, C: []string{"a", "b", "c"}, D: []int{1, 2, 3}},
		},
		{in: struct { // wrong tag
			N int `validation:"min:a"`
		}{N: 5}},
		{in: App{Version: "12345"}},            // corrrect
		{in: Response{Code: 200, Body: "any"}}, // correct
		{in: User{ // correct
			ID:     "123456789012345678901234567890123456",
			Name:   "any name",
			Age:    20,
			Email:  "test@test.com",
			Role:   "admin",
			Phones: []string{"79990009900", "70009990099"},
			meta:   json.RawMessage{},
		}},
	}

	for i, tt := range testsNegative {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			fmt.Printf("%v", err)
			require.ErrorIs(t, err, tt.expectedErr)
			for _, c := range tt.contains {
				require.ErrorContains(t, err, c)
			}
		})
	}

	for i, tt := range testsPositive {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			res := Validate(tt.in)
			require.Equal(t, nil, res)
			_ = tt
		})
	}
}
