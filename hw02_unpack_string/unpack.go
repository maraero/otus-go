package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var result strings.Builder

	var prevChar string

	for _, char := range input {
		isDigit := unicode.IsDigit(char)
		isPrevEmpty := prevChar == ""

		switch {
		case isDigit && isPrevEmpty:
			return "", ErrInvalidString
		case isDigit && !isPrevEmpty:
			num, e := strconv.Atoi(string(char))
			if e != nil {
				return "", e
			}
			for i := 0; i < num; i++ {
				result.WriteString(prevChar)
			}
			prevChar = ""
		case !isDigit && isPrevEmpty:
			prevChar = string(char)
		case !isDigit && !isPrevEmpty:
			result.WriteString(prevChar)
			prevChar = string(char)
		}
	}

	if prevChar != "" {
		result.WriteString(prevChar)
	}

	return result.String(), nil
}
