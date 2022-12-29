package eventservice

import (
	"context"
	"time"

	event "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/domain"
	repository_memory "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/repository-memory"
	repository_sql "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/repository-sql"
	storage "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/storage"
)

func New(strg *storage.Storage) *EventService {
	if strg.Source == storage.StorageSQL {
		return &EventService{repository: repository_sql.New(strg.Connection)}
	}
	return &EventService{repository: repository_memory.New()}
}

func (es *EventService) CreateEvent(ctx context.Context, e event.Event) (id int64, err error) {
	if err = e.Validate(); err != nil {
		return 0, err
	}

	id, err = es.repository.CreateEvent(ctx, e)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (es *EventService) UpdateEvent(ctx context.Context, id int64, e event.Event) error {
	if err := e.Validate(); err != nil {
		return err
	}
	return es.repository.UpdateEvent(ctx, id, e)
}

func (es *EventService) DeleteEvent(ctx context.Context, id int64) error {
	return es.repository.DeleteEvent(ctx, id)
}

func (es *EventService) GetEventListByDate(ctx context.Context, date time.Time) ([]event.Event, error) {
	return es.repository.GetEventListByDate(ctx, date)
}

func (es *EventService) GetEventListByWeek(ctx context.Context, date time.Time) ([]event.Event, error) {
	return es.repository.GetEventListByWeek(ctx, date)
}

func (es *EventService) GetEventListByMonth(ctx context.Context, date time.Time) ([]event.Event, error) {
	return es.repository.GetEventListByMonth(ctx, date)
}

func (es *EventService) GetEventByID(ctx context.Context, id int64) (event.Event, error) {
	return es.repository.GetEventByID(ctx, id)
}
