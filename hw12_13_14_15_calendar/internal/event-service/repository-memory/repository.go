package eventrepositorymemory

import (
	"context"
	"sort"
	"time"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/events"
)

func New() *Repository {
	return &Repository{events: make(map[int64]events.Event)}
}

func (s *Repository) CreateEvent(_ context.Context, e events.Event) (int64, error) {
	s.Lock()
	defer s.Unlock()
	id := s.next()
	e.ID = id
	s.events[id] = e
	return id, nil
}

func (s *Repository) UpdateEvent(_ context.Context, id int64, e events.Event) error {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.events[id]; !ok {
		return events.ErrNotFound
	}
	s.events[id] = e
	return nil
}

func (s *Repository) DeleteEvent(_ context.Context, id int64) error {
	s.Lock()
	defer s.Unlock()
	event, ok := s.events[id]
	if !ok {
		return events.ErrNotFound
	}
	event.Deleted = true
	s.events[id] = event
	return nil
}

func (s *Repository) GetEventListByDate(_ context.Context, date time.Time) ([]events.Event, error) {
	s.Lock()
	defer s.Unlock()
	var res []events.Event
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

func (s *Repository) GetEventListByWeek(_ context.Context, date time.Time) ([]events.Event, error) {
	s.Lock()
	defer s.Unlock()
	var res []events.Event
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

func (s *Repository) GetEventListByMonth(_ context.Context, date time.Time) ([]events.Event, error) {
	s.Lock()
	defer s.Unlock()
	var res []events.Event
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

func (s *Repository) GetEventByID(_ context.Context, id int64) (events.Event, error) {
	res, ok := s.events[id]
	if ok {
		return res, nil
	}

	return events.Event{}, events.ErrNotFound
}

func (s *Repository) next() int64 {
	s.last++
	return s.last
}

func order(events []events.Event) []events.Event {
	sort.Slice(events, func(i, j int) bool {
		return events[i].DateStart.Before(events[j].DateStart)
	})
	return events
}
