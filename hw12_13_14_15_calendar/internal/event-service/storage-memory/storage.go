package memorystorage

import (
	"context"
	"sort"
	"sync"
	"time"

	evt "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/domain"
)

type Storage struct {
	sync.RWMutex
	events map[int64]evt.Event
	last   int64
}

func New() *Storage {
	return &Storage{events: make(map[int64]evt.Event)}
}

func (s *Storage) Connect(_ context.Context, _ string, _ string) error {
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}

func (s *Storage) CreateEvent(_ context.Context, e evt.Event) (int64, error) {
	s.Lock()
	defer s.Unlock()
	id := s.next()
	e.ID = id
	s.events[id] = e
	return id, nil
}

func (s *Storage) UpdateEvent(_ context.Context, id int64, e evt.Event) error {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.events[id]; !ok {
		return evt.ErrNotFound
	}
	s.events[id] = e
	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id int64) error {
	s.Lock()
	defer s.Unlock()
	event, ok := s.events[id]
	if !ok {
		return evt.ErrNotFound
	}
	event.Deleted = true
	s.events[id] = event
	return nil
}

func (s *Storage) GetEventListByDate(_ context.Context, date time.Time) ([]evt.Event, error) {
	s.Lock()
	defer s.Unlock()
	var res []evt.Event
	year, month, day := date.Date()

	for _, event := range s.events {
		if event.Deleted {
			continue
		}
		eYear, eMonth, eDay := event.DateStart.Date()
		if eDay == day && eMonth == month && eYear == year {
			res = append(res, event)
		}
	}

	return order(res), nil
}

func (s *Storage) GetEventListByWeek(_ context.Context, date time.Time) ([]evt.Event, error) {
	s.Lock()
	defer s.Unlock()
	var res []evt.Event
	year, week := date.ISOWeek()

	for _, event := range s.events {
		if event.Deleted {
			continue
		}
		eYear, eWeek := event.DateStart.ISOWeek()
		if eWeek == week && eYear == year {
			res = append(res, event)
		}
	}

	return order(res), nil
}

func (s *Storage) GetEventListByMonth(_ context.Context, date time.Time) ([]evt.Event, error) {
	s.Lock()
	defer s.Unlock()
	var res []evt.Event
	year, month, _ := date.Date()

	for _, event := range s.events {
		if event.Deleted {
			continue
		}
		eYear, eMonth, _ := event.DateStart.Date()
		if eMonth == month && eYear == year {
			res = append(res, event)
		}
	}

	return order(res), nil
}

func (s *Storage) next() int64 {
	s.last++
	return s.last
}

func order(events []evt.Event) []evt.Event {
	sort.Slice(events, func(i, j int) bool {
		return events[i].DateStart.Before(events[j].DateStart)
	})
	return events
}
