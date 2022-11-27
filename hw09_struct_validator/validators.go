package hw09structvalidator

func validateInt(val int64, rList []Rule) error {
	for _, rule := range rList {
		switch rule.Name {
		case ruleIntMin:
			err := validateIntMin(val, rule.Val)
			if err != nil {
				return err
			}
		case ruleIntMax:
			err := validateIntMax(val, rule.Val)
			if err != nil {
				return err
			}
		case ruleIntIn:
			err := validateIntIn(val, rule.Val)
			if err != nil {
				return err
			}
		default:
			return ErrUnknownIntRule
		}
	}
	return nil
}

func validateString(val string, rList []Rule) error {
	for _, rule := range rList {
		switch rule.Name {
		case ruleStringLen:
			err := validateStringLen(val, rule.Val)
			if err != nil {
				return err
			}
		case ruleStringRegexp:
			err := validateStringRegexp(val, rule.Val)
			if err != nil {
				return err
			}
		case ruleStringIn:
			err := validateStringIn(val, rule.Val)
			if err != nil {
				return err
			}
		default:
			return ErrUnknownStringRule
		}
	}
	return nil
}

func validateSliceString(vals []string, rList []Rule) error {
	for _, val := range vals {
		err := validateString(val, rList)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateSliceInt(vals []int64, rList []Rule) error {
	for _, val := range vals {
		err := validateInt(val, rList)
		if err != nil {
			return err
		}
	}
	return nil
}
