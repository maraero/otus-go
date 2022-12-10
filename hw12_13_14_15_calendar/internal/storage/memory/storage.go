package memorystorage

import (
	"context"
	"sync"
	"time"

	es "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/storage/event-service"
)

type Storage struct {
	sync.RWMutex
	events map[string]es.Event
}

func New() *Storage {
	return &Storage{events: make(map[string]es.Event)}
}

func (s *Storage) Connect(_ context.Context, _ string) error {
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}

func (s *Storage) CreateEvent(_ context.Context, e es.Event) (es.Event, error) {
	s.Lock()
	defer s.Unlock()
	s.events[e.ID] = e
	return s.events[e.ID], nil
}

func (s *Storage) UpdateEvent(_ context.Context, id string, e es.Event) (es.Event, error) {
	s.Lock()
	defer s.Unlock()
	_, ok := s.events[id]
	if !ok {
		return es.Event{}, es.ErrNotFound
	}
	s.events[id] = e
	return s.events[id], nil
}

func (s *Storage) DeleteEvent(_ context.Context, id string) error {
	s.Lock()
	defer s.Unlock()
	event, ok := s.events[id]
	if !ok {
		return es.ErrNotFound
	}
	event.Deleted = true
	s.events[id] = event
	return nil
}

func (s *Storage) GetEventListByDate(_ context.Context, date time.Time) []es.Event {
	s.Lock()
	defer s.Unlock()
	var res []es.Event
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

	return res
}

func (s *Storage) GetEventListByWeek(_ context.Context, date time.Time) []es.Event {
	s.Lock()
	defer s.Unlock()
	var res []es.Event
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

	return res
}

func (s *Storage) GetEventListByMonth(_ context.Context, date time.Time) []es.Event {
	s.Lock()
	defer s.Unlock()
	var res []es.Event
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

	return res
}
