package event_service

import (
	"context"
	"time"
)

type Event struct {
	ID               int64     `db:"id"`
	Title            string    `db:"title"`
	DateStart        time.Time `db:"date_start"`
	DateEnd          time.Time `db:"date_end"`
	Descripion       string    `db:"description"`
	UserId           string    `db:"user_id"`
	DateNotification time.Time `db:"date_notification"`
	Deleted          bool      `db:"deleted"`
}

type Storager interface {
	Connect(ctx context.Context, dsn string) error
	Close(ctx context.Context) error
	CreateEvent(ctx context.Context, e Event) (int64, error)
	UpdateEvent(ctx context.Context, id int64, e Event) error
	DeleteEvent(ctx context.Context, id int64) error
	GetEventListByDate(ctx context.Context, date time.Time) ([]Event, error)
	GetEventListByWeek(ctx context.Context, date time.Time) ([]Event, error)
	GetEventListByMonth(ctx context.Context, date time.Time) ([]Event, error)
}

type EventService struct {
	storage *Storager
}
