package event_service

import (
	"context"
	"time"
)

type Event struct {
	ID               string
	Title            string
	DateStart        time.Time
	DateEnd          time.Time
	Descripion       string
	UserId           string
	DateNotification time.Time
	Deleted          bool
}

type Storager interface {
	Connect(ctx context.Context, dsn string) error
	Close(ctx context.Context) error
	CreateEvent(ctx context.Context, e Event) (Event, error)
	UpdateEvent(ctx context.Context, id string, e Event) (Event, error)
	DeleteEvent(ctx context.Context, id string) error
	GetEventListByDate(ctx context.Context, date time.Time) []Event
	GetEventListByWeek(ctx context.Context, date time.Time) []Event
	GetEventListByMonth(ctx context.Context, date time.Time) []Event
}

type EventService struct {
	storage *Storager
}
