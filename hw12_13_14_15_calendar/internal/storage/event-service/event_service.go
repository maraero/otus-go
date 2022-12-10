package eventservice

import (
	"time"

	"github.com/google/uuid"
)

func NewEventService() *EventService {
	return &EventService{}
}

func (es EventService) CreateEvent(e Event) (Event, error) {
	err := es.validateEvent(e)
	if err != nil {
		return Event{}, err
	}
	e.ID = uuid.NewString()
	return e, nil
}

func (es EventService) UpdateEvent(id string, e Event) (Event, error) {
	err := es.validateEvent(e)
	if err != nil {
		return Event{}, err
	}
	return e, nil
}

func (es EventService) DeleteEvent(id string) error {
	return nil
}

func (es EventService) GetEventListByDate(date time.Time) []Event {
	return []Event{}
}

func (es EventService) GetEventListByWeek(date time.Time) []Event {
	return []Event{}
}

func (es EventService) GetEventListByMonth(date time.Time) []Event {
	return []Event{}
}

func (es EventService) validateEvent(e Event) error {
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
