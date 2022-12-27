package eventservice

import (
	"context"
	"time"

	evt "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/domain"
)

type Storage interface {
	CreateEvent(ctx context.Context, e evt.Event) (int64, error)
	UpdateEvent(ctx context.Context, id int64, e evt.Event) error
	DeleteEvent(ctx context.Context, id int64) error
	GetEventListByDate(ctx context.Context, date time.Time) ([]evt.Event, error)
	GetEventListByWeek(ctx context.Context, date time.Time) ([]evt.Event, error)
	GetEventListByMonth(ctx context.Context, date time.Time) ([]evt.Event, error)
	GetEventById(ctx context.Context, id int64) (evt.Event, error)
}

type EventService struct {
	storage Storage
}
