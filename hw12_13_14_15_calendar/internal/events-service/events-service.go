package eventsservice

import (
	"context"
	"time"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/events"
	er "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/events-repository"
)

type EventsService struct {
	repository er.EventsRepository
}

func New(repository er.EventsRepository) *EventsService {
	return &EventsService{repository: repository}
}

func (es *EventsService) CreateEvent(ctx context.Context, e events.Event) (id int64, err error) {
	if err = e.Validate(); err != nil {
		return 0, err
	}

	id, err = es.repository.CreateEvent(ctx, e)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (es *EventsService) UpdateEvent(ctx context.Context, id int64, e events.Event) error {
	if err := e.Validate(); err != nil {
		return err
	}
	return es.repository.UpdateEvent(ctx, id, e)
}

func (es *EventsService) DeleteEvent(ctx context.Context, id int64) error {
	return es.repository.DeleteEvent(ctx, id)
}

func (es *EventsService) GetEventListByDate(ctx context.Context, date time.Time) ([]events.Event, error) {
	return es.repository.GetEventListByDate(ctx, date)
}

func (es *EventsService) GetEventListByWeek(ctx context.Context, date time.Time) ([]events.Event, error) {
	return es.repository.GetEventListByWeek(ctx, date)
}

func (es *EventsService) GetEventListByMonth(ctx context.Context, date time.Time) ([]events.Event, error) {
	return es.repository.GetEventListByMonth(ctx, date)
}

func (es *EventsService) GetEventByID(ctx context.Context, id int64) (events.Event, error) {
	return es.repository.GetEventByID(ctx, id)
}
