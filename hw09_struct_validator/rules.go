package hw09structvalidator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Rule struct {
	Name string
	Val  string
}

const ruleListSep = "|"

const (
	ruleSep = ":"
	listSep = ","
)

const ruleStringLen = "len"

const (
	ruleStringRegexp = "regexp"
	ruleStringIn     = "in"
	ruleIntMin       = "min"
	ruleIntMax       = "max"
	ruleIntIn        = "in"
)

func getRuleListFromTag(tag string) []Rule {
	var rules []Rule
	for _, rule := range strings.Split(tag, ruleListSep) {
		ruleParts := strings.SplitN(rule, ruleSep, 2)

		if len(ruleParts) == 2 {
			rules = append(rules, Rule{Name: ruleParts[0], Val: ruleParts[1]})
		}
	}
	return rules
}

func validateIntMin(val int64, min string) error {
	minimum, err := strconv.ParseInt(min, 10, 64)
	if err != nil {
		return fmt.Errorf("%w: can not convert %v to int", ErrInvalidRuleParam, min)
	}
	if val < minimum {
		return fmt.Errorf("%w: must be not less than %v", ErrValidation, minimum)
	}
	return nil
}

func validateIntMax(val int64, max string) error {
	maximum, err := strconv.ParseInt(max, 10, 64)
	if err != nil {
		return fmt.Errorf("%w: can not convert %v to int", ErrInvalidRuleParam, max)
	}
	if val > maximum {
		return fmt.Errorf("%w: must be not more than %v", ErrValidation, maximum)
	}
	return nil
}

func validateIntIn(val int64, in string) error {
	strSet := strings.Split(in, listSep)
	var intSet []int64
	for _, s := range strSet {
		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("%w: can not conver list element %v to int", ErrInvalidRuleParam, s)
		}
		intSet = append(intSet, val)
	}
	for _, v := range intSet {
		if val == v {
			return nil
		}
	}
	return fmt.Errorf("%w: does not match any value in list %v", ErrValidation, in)
}

func validateStringLen(val string, ruleLen string) error {
	length, err := strconv.Atoi(ruleLen)
	if err != nil {
		return fmt.Errorf("%w: can not convert %v to int", ErrInvalidRuleParam, ruleLen)
	}
	if len(val) == length {
		return nil
	}
	return fmt.Errorf("%w: length must be %v", ErrValidation, ruleLen)
}

func validateStringRegexp(val string, regStr string) error {
	reg, err := regexp.Compile(regStr)
	if err != nil {
		return fmt.Errorf("%w: can not conver %v to valid regexp", ErrInvalidRuleParam, regStr)
	}
	if reg.MatchString(val) {
		return nil
	}
	return fmt.Errorf("%w: must match regexp \"%v\"", ErrValidation, regStr)
}

func validateStringIn(val string, in string) error {
	strSet := strings.Split(in, listSep)
	for _, s := range strSet {
		if s == val {
			return nil
		}
	}
	return fmt.Errorf("%w: does not match any value in list %v", ErrValidation, in)
}
