package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	eS "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/storage/event-service"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestCreateEvent(t *testing.T) {
	storage := New()
	newEvent := eS.Event{
		ID:        uuid.NewString(),
		Title:     "create test event",
		DateStart: time.Now(),
		DateEnd:   time.Now().Add(2 * time.Hour),
		UserId:    "create test user id",
	}
	event, err := storage.CreateEvent(context.Background(), newEvent)
	require.NoError(t, err)
	require.Equal(t, newEvent, event)
}

func TestUpdateEvent(t *testing.T) {
	storage := New()
	initialEvent := eS.Event{
		ID:        uuid.NewString(),
		Title:     "update test event",
		DateStart: time.Now(),
		DateEnd:   time.Now().Add(2 * time.Hour),
		UserId:    "update test user id",
	}
	_, err := storage.UpdateEvent(context.Background(), "unknown id", initialEvent)
	require.Error(t, err)
	require.Equal(t, err, eS.ErrNotFound)
	event, err := storage.CreateEvent(context.Background(), initialEvent)
	require.NoError(t, err)
	require.Equal(t, event, initialEvent)
	updatedInitialEvent := initialEvent
	newTitle := "very updated"
	updatedInitialEvent.Title = newTitle
	updatedEvent, err := storage.UpdateEvent(context.Background(), initialEvent.ID, updatedInitialEvent)
	require.NoError(t, err)
	require.NotEqual(t, initialEvent, updatedEvent)
	require.Equal(t, newTitle, updatedEvent.Title)
}

func TestDeleteEvent(t *testing.T) {
	storage := New()
	initialEvent := eS.Event{
		ID:        uuid.NewString(),
		Title:     "test event",
		DateStart: time.Now(),
		DateEnd:   time.Now().Add(2 * time.Hour),
		UserId:    "test user id",
	}
	err := storage.DeleteEvent(context.Background(), "unknown id")
	require.Error(t, err)
	require.Equal(t, err, eS.ErrNotFound)
	event, err := storage.CreateEvent(context.Background(), initialEvent)
	require.NoError(t, err)
	require.Equal(t, event, initialEvent)
	err = storage.DeleteEvent(context.Background(), initialEvent.ID)
	require.NoError(t, err)
	require.True(t, storage.events[initialEvent.ID].Deleted)
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

	m.storage.events = map[string]eS.Event{
		"1": {
			ID:        "1",
			Title:     "title 1",
			DateStart: time.Now(),
			DateEnd:   time.Now().Add(2 * time.Hour),
			UserId:    "user id 1",
		},
		"2": {
			ID:        "2",
			Title:     "title 2",
			DateStart: time.Now().Add(2 * weekDuration),
			DateEnd:   time.Now().Add(2*weekDuration + 2*time.Hour),
			UserId:    "user id 1",
		},
		"3": {
			ID:        "3",
			Title:     "title 3",
			DateStart: time.Now().Add(2 * monthDuration),
			DateEnd:   time.Now().Add(2*monthDuration + 2*time.Hour),
			UserId:    "user id 1",
		},
	}
}

func (m *MemoryStorageSuite) TestGetEventListByDate() {
	m.Run("success event list by date", func() {
		res := m.storage.GetEventListByDate(context.Background(), time.Now())
		require.Equal(m.T(), 1, len(res))
		require.Equal(m.T(), "1", res[0].ID)
	})

	m.Run("empty event list by date", func() {
		res := m.storage.GetEventListByDate(context.Background(), time.Now().Add(dayDuration))
		require.Equal(m.T(), 0, len(res))
	})
}

func (m *MemoryStorageSuite) TestGetEventListByWeek() {
	m.Run("success event list by week", func() {
		res := m.storage.GetEventListByWeek(context.Background(), time.Now().Add(2*weekDuration))
		require.Equal(m.T(), 1, len(res))
		require.Equal(m.T(), "2", res[0].ID)
	})

	m.Run("empty event list by week", func() {
		res := m.storage.GetEventListByWeek(context.Background(), time.Now().Add(4*weekDuration))
		require.Equal(m.T(), 0, len(res))
	})
}

func (m *MemoryStorageSuite) TestGetEventListByMonth() {
	m.Run("success event list by month", func() {
		res := m.storage.GetEventListByMonth(context.Background(), time.Now().Add(2*monthDuration))
		require.Equal(m.T(), 1, len(res))
		require.Equal(m.T(), "3", res[0].ID)
	})

	m.Run("empty event list by month", func() {
		res := m.storage.GetEventListByMonth(context.Background(), time.Now().Add(4*monthDuration))
		require.Equal(m.T(), 0, len(res))
	})
}

func TestMemoryStorageSuite(t *testing.T) {
	suite.Run(t, new(MemoryStorageSuite))
}
