package eventservice

import (
	"context"
	"time"

	config "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
	event "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/domain"
	memorystorage "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/storage-memory"
	sqlstorage "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/storage-sql"
)

func New(ctx context.Context, storageType string, sqlDriver string, DSN string) (*EventService, error) {
	var storage Storage
	if storageType == config.StorageInMemory {
		storage = memorystorage.New()
	} else {
		storage = sqlstorage.New()
	}

	err := storage.Connect(ctx, sqlDriver, DSN)
	if err != nil {
		return nil, err
	}

	return &EventService{storage: storage}, nil
}

func (es *EventService) CreateEvent(ctx context.Context, e event.Event) (id int64, err error) {
	err = es.validateEvent(e)
	if err != nil {
		return 0, err
	}

	id, err = es.storage.CreateEvent(ctx, e)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (es *EventService) UpdateEvent(ctx context.Context, id int64, e event.Event) error {
	err := es.validateEvent(e)
	if err != nil {
		return err
	}
	return es.storage.UpdateEvent(ctx, id, e)
}

func (es *EventService) DeleteEvent(ctx context.Context, id int64) error {
	return es.storage.DeleteEvent(ctx, id)
}

func (es *EventService) GetEventListByDate(ctx context.Context, date time.Time) ([]event.Event, error) {
	return es.storage.GetEventListByDate(ctx, date)
}

func (es *EventService) GetEventListByWeek(ctx context.Context, date time.Time) ([]event.Event, error) {
	return es.storage.GetEventListByWeek(ctx, date)
}

func (es *EventService) GetEventListByMonth(ctx context.Context, date time.Time) ([]event.Event, error) {
	return es.storage.GetEventListByMonth(ctx, date)
}

func (es *EventService) validateEvent(e event.Event) error {
	if e.Title == "" {
		return event.ErrEmptyTitle
	}

	if e.DateStart.IsZero() {
		return event.ErrEmptyDateStart
	}

	if e.DateEnd.IsZero() {
		return event.ErrEmptyDateEnd
	}

	if e.DateStart.After(e.DateEnd) {
		return event.ErrInvalidDates
	}

	if e.DateEnd.Before(time.Now()) {
		return event.ErrEndInThePast
	}

	return nil
}