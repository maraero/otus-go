package events

import "errors"

var (
	ErrEmptyTitle     = errors.New("emtpty title")
	ErrEmptyDateStart = errors.New("empty start date")
	ErrEmptyDateEnd   = errors.New("empty end date")
	ErrInvalidDates   = errors.New("stard date can not be after end date")
	ErrEndInThePast   = errors.New("end date already passed")
	ErrNotFound       = errors.New("event not found")
)
