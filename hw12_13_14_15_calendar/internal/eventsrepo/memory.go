package eventsrepo

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/events"
)

type MemoryRepo struct {
	sync.RWMutex
	events map[int64]events.Event
	last   int64
}

func NewMemoryRepository() Repository {
	return &MemoryRepo{events: make(map[int64]events.Event)}
}

func (r *MemoryRepo) CreateEvent(_ context.Context, e events.Event) (int64, error) {
	r.Lock()
	defer r.Unlock()
	id := r.next()
	e.ID = id
	r.events[id] = e
	return id, nil
}

func (r *MemoryRepo) UpdateEvent(_ context.Context, id int64, e events.Event) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.events[id]; !ok {
		return events.ErrNotFound
	}
	r.events[id] = e
	return nil
}

func (r *MemoryRepo) DeleteEvent(_ context.Context, id int64) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.events[id]; !ok {
		return events.ErrNotFound
	}
	delete(r.events, id)
	return nil
}

func (r *MemoryRepo) GetEventListByDate(_ context.Context, date time.Time) ([]events.Event, error) {
	r.Lock()
	defer r.Unlock()
	var res []events.Event
	year, month, day := date.Date()

	for _, event := range r.events {
		eYear, eMonth, eDay := event.DateStart.Date()
		if eDay == day && eMonth == month && eYear == year {
			res = append(res, event)
		}
	}

	return order(res), nil
}

func (r *MemoryRepo) GetEventListByWeek(_ context.Context, date time.Time) ([]events.Event, error) {
	r.Lock()
	defer r.Unlock()
	var res []events.Event
	year, week := date.ISOWeek()

	for _, event := range r.events {
		eYear, eWeek := event.DateStart.ISOWeek()
		if eWeek == week && eYear == year {
			res = append(res, event)
		}
	}

	return order(res), nil
}

func (r *MemoryRepo) GetEventListByMonth(_ context.Context, date time.Time) ([]events.Event, error) {
	r.Lock()
	defer r.Unlock()
	var res []events.Event
	year, month, _ := date.Date()

	for _, event := range r.events {
		eYear, eMonth, _ := event.DateStart.Date()
		if eMonth == month && eYear == year {
			res = append(res, event)
		}
	}

	return order(res), nil
}

func (r *MemoryRepo) GetEventByID(_ context.Context, id int64) (events.Event, error) {
	res, ok := r.events[id]
	if ok {
		return res, nil
	}

	return events.Event{}, events.ErrNotFound
}

func (r *MemoryRepo) next() int64 {
	r.last++
	return r.last
}

func order(events []events.Event) []events.Event {
	sort.Slice(events, func(i, j int) bool {
		return events[i].DateStart.Before(events[j].DateStart)
	})
	return events
}
