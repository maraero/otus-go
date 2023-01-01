package eventservice

import (
	"context"
	"time"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/events"
)

type Repository interface {
	CreateEvent(ctx context.Context, e events.Event) (int64, error)
	UpdateEvent(ctx context.Context, id int64, e events.Event) error
	DeleteEvent(ctx context.Context, id int64) error
	GetEventListByDate(ctx context.Context, date time.Time) ([]events.Event, error)
	GetEventListByWeek(ctx context.Context, date time.Time) ([]events.Event, error)
	GetEventListByMonth(ctx context.Context, date time.Time) ([]events.Event, error)
	GetEventByID(ctx context.Context, id int64) (events.Event, error)
}

type EventService struct {
	repository Repository
}
