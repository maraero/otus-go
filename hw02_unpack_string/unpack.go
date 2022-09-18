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
	var startBackslash bool

	for _, char := range input {
		isBackslash := string(char) == "\\"
		isDigit := unicode.IsDigit(char)

		switch {
		case isBackslash && startBackslash:
			prevChar = "\\"
			startBackslash = false
		case isBackslash && prevChar == "":
			startBackslash = true
		case isBackslash && prevChar != "":
			result.WriteString(prevChar)
			prevChar = ""
			startBackslash = true

		case isDigit && startBackslash:
			prevChar = string(char)
			startBackslash = false
		case isDigit && prevChar == "":
			return "", ErrInvalidString
		case isDigit && prevChar != "":
			num, e := strconv.Atoi(string(char))
			if e != nil {
				return "", e
			}
			result.WriteString(strings.Repeat(prevChar, num))
			prevChar = ""

		case !isDigit && startBackslash:
			return "", ErrInvalidString
		case !isDigit && prevChar == "":
			prevChar = string(char)
		case !isDigit && prevChar != "":
			result.WriteString(prevChar)
			prevChar = string(char)
		}
	}

	if startBackslash {
		return "", ErrInvalidString
	}

	if prevChar != "" {
		result.WriteString(prevChar)
	}

	return result.String(), nil
}
