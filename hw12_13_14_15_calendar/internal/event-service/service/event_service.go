package eventservice

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	event "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/domain"
	memory_storage "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/storage-memory"
	sql_storage "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/storage-sql"
)

func New(dbConn *sqlx.DB) *EventService {
	var storage Storage
	if dbConn == nil {
		storage = memory_storage.New()
	} else {
		storage = sql_storage.New(dbConn)
	}
	return &EventService{storage: storage}
}

func (es *EventService) CreateEvent(ctx context.Context, e event.Event) (id int64, err error) {
	if err = e.Validate(); err != nil {
		return 0, err
	}

	id, err = es.storage.CreateEvent(ctx, e)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (es *EventService) UpdateEvent(ctx context.Context, id int64, e event.Event) error {
	if err := e.Validate(); err != nil {
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

func (es *EventService) GetEventById(ctx context.Context, id int64) (event.Event, error) {
	return es.storage.GetEventById(ctx, id)
}
