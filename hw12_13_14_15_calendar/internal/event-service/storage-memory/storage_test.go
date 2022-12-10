package memorystorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	evt "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/domain"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestCreateEvent(t *testing.T) {
	storage := New()
	newEvent := evt.Event{
		Title:     "create test event",
		DateStart: time.Now(),
		DateEnd:   time.Now().Add(2 * time.Hour),
		UserID:    "create test user id",
	}
	id, err := storage.CreateEvent(context.Background(), newEvent)
	require.NoError(t, err)
	require.Equal(t, id, int64(1))
	fmt.Println(storage.events)
	require.Equal(t, newEvent.Title, storage.events[id].Title)
}

func TestUpdateEvent(t *testing.T) {
	storage := New()
	initialEvent := evt.Event{
		Title:     "update test event",
		DateStart: time.Now(),
		DateEnd:   time.Now().Add(2 * time.Hour),
		UserID:    "update test user id",
	}
	err := storage.UpdateEvent(context.Background(), int64(5), initialEvent)
	require.Error(t, err)
	require.Equal(t, err, evt.ErrNotFound)
	id, err := storage.CreateEvent(context.Background(), initialEvent)
	require.NoError(t, err)
	require.Equal(t, id, int64(1))
	updatedInitialEvent := initialEvent
	newTitle := "very updated"
	updatedInitialEvent.Title = newTitle
	err = storage.UpdateEvent(context.Background(), id, updatedInitialEvent)
	require.NoError(t, err)
	require.Equal(t, newTitle, storage.events[id].Title)
}

func TestDeleteEvent(t *testing.T) {
	storage := New()
	initialEvent := evt.Event{
		Title:     "test event",
		DateStart: time.Now(),
		DateEnd:   time.Now().Add(2 * time.Hour),
		UserID:    "test user id",
	}
	err := storage.DeleteEvent(context.Background(), int64(5))
	require.Error(t, err)
	require.Equal(t, err, evt.ErrNotFound)
	id, err := storage.CreateEvent(context.Background(), initialEvent)
	require.NoError(t, err)
	err = storage.DeleteEvent(context.Background(), id)
	require.NoError(t, err)
	require.True(t, storage.events[id].Deleted)
}

type MemoryStorageSuite struct {
	suite.Suite
	storage *Storage
}

var dayDuration = 24 * time.Hour
var weekDuration = 7 * dayDuration
var monthDuration = 30 * dayDuration // not precisely

func (m *MemoryStorageSuite) SetupTest() {
	m.storage = New()

	m.storage.events = map[int64]evt.Event{
		1: {
			ID:        int64(1),
			Title:     "title 1",
			DateStart: time.Now(),
			DateEnd:   time.Now().Add(2 * time.Hour),
			UserID:    "user id 1",
		},
		2: {
			ID:        int64(2),
			Title:     "title 2",
			DateStart: time.Now().Add(2 * weekDuration),
			DateEnd:   time.Now().Add(2*weekDuration + 2*time.Hour),
			UserID:    "user id 1",
		},
		3: {
			ID:        int64(3),
			Title:     "title 3",
			DateStart: time.Now().Add(2 * monthDuration),
			DateEnd:   time.Now().Add(2*monthDuration + 2*time.Hour),
			UserID:    "user id 1",
		},
	}
}

func (m *MemoryStorageSuite) TestGetEventListByDate() {
	m.Run("success event list by date", func() {
		res, err := m.storage.GetEventListByDate(context.Background(), time.Now())
		require.NoError(m.T(), err)
		require.Equal(m.T(), 1, len(res))
		require.Equal(m.T(), int64(1), res[0].ID)
	})

	m.Run("empty event list by date", func() {
		res, err := m.storage.GetEventListByDate(context.Background(), time.Now().Add(dayDuration))
		require.NoError(m.T(), err)
		require.Equal(m.T(), 0, len(res))
	})
}

func (m *MemoryStorageSuite) TestGetEventListByWeek() {
	m.Run("success event list by week", func() {
		res, err := m.storage.GetEventListByWeek(context.Background(), time.Now().Add(2*weekDuration))
		require.NoError(m.T(), err)
		require.Equal(m.T(), 1, len(res))
		require.Equal(m.T(), int64(2), res[0].ID)
	})

	m.Run("empty event list by week", func() {
		res, err := m.storage.GetEventListByWeek(context.Background(), time.Now().Add(4*weekDuration))
		require.NoError(m.T(), err)
		require.Equal(m.T(), 0, len(res))
	})
}

func (m *MemoryStorageSuite) TestGetEventListByMonth() {
	m.Run("success event list by month", func() {
		res, err := m.storage.GetEventListByMonth(context.Background(), time.Now().Add(2*monthDuration))
		require.NoError(m.T(), err)
		require.Equal(m.T(), 1, len(res))
		require.Equal(m.T(), int64(3), res[0].ID)
	})

	m.Run("empty event list by month", func() {
		res, err := m.storage.GetEventListByMonth(context.Background(), time.Now().Add(4*monthDuration))
		require.NoError(m.T(), err)
		require.Equal(m.T(), 0, len(res))
	})
}

func TestMemoryStorageSuite(t *testing.T) {
	suite.Run(t, new(MemoryStorageSuite))
}
