package eventsrepository

import (
	"context"
	"time"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/events"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type EventsRepository interface {
	CreateEvent(ctx context.Context, e events.Event) (int64, error)
	UpdateEvent(ctx context.Context, id int64, e events.Event) error
	DeleteEvent(ctx context.Context, id int64) error
	GetEventListByDate(ctx context.Context, date time.Time) ([]events.Event, error)
	GetEventListByWeek(ctx context.Context, date time.Time) ([]events.Event, error)
	GetEventListByMonth(ctx context.Context, date time.Time) ([]events.Event, error)
	GetEventByID(ctx context.Context, id int64) (events.Event, error)
}

func New(strg *storage.Storage) EventsRepository {
	switch strg.Source {
	case storage.StorageInMemory:
		return newMemoryRepository()
	case storage.StorageSQL:
		return newSQLRepository(strg.Connection)
	default:
		return newMemoryRepository()
	}
}
