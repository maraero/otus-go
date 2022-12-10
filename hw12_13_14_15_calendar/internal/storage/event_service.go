package storage

import (
	"time"

	"github.com/google/uuid"
)

func CreateEvent(e Event) (Event, error) {
	err := validateEvent(e)
	if err != nil {
		return Event{}, err
	}
	e.ID = uuid.New().String()
	return e, nil
}

func UpdateEvent(id string, e Event) (Event, error) {
	err := validateEvent(e)
	if err != nil {
		return Event{}, err
	}
	return e, nil
}

func DeleteEvent(id string) error {
	return nil
}

func GetEventListByDate(date time.Time) []Event {
	return []Event{}
}

func GetEventListByWeek(date time.Time) []Event {
	return []Event{}
}

func GetEventListByMonth(date time.Time) []Event {
	return []Event{}
}

func validateEvent(e Event) error {
	if e.Title == "" {
		return ErrEmptyTitle
	}

	if e.DateStart.IsZero() {
		return ErrEmptyDateStart
	}

	if e.DateEnd.IsZero() {
		return ErrEmptyDateEnd
	}

	if e.DateStart.After(e.DateEnd) {
		return ErrInvalidDates
	}

	if e.DateEnd.Before(time.Now()) {
		return ErrEndInThePast
	}

	return nil
}
