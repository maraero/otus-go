package serverhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
	event "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/domain"
	es "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/service"
	l "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/stretchr/testify/suite"
)

type SuiteTest struct {
	suite.Suite
	ts *httptest.Server
}

func buildTestConfig() config.Config {
	loggerConfig := config.Logger{
		Level:            "info",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	serverConfig := config.Server{
		Host:     "localhost",
		HTTPPort: "8000",
		GrpcPort: "8001",
	}
	storageConfig := config.Storage{
		Type:   "in-memory",
		Driver: "postgres",
		DSN:    "postgresql://admin:admin@localhost:5432/calendar?sslmode=disable",
	}
	return config.Config{
		Logger:  loggerConfig,
		Server:  serverConfig,
		Storage: storageConfig,
	}
}

func (s *SuiteTest) SetupTest() {
	config := buildTestConfig()

	logger, err := l.New(config.Logger)
	s.Require().NoError(err)

	var dbConnection *sqlx.DB
	eventService := es.New(dbConnection)
	calendar := app.New(eventService, logger)
	s.ts = httptest.NewServer(New(calendar, config.Server).router)
}

func (s *SuiteTest) TeardownTest() {
	s.ts.Close()
}

func TestHttpServer(t *testing.T) {
	suite.Run(t, new(SuiteTest))
}

func (s *SuiteTest) TestRoot() {
	client := &http.Client{}
	rootURL := s.ts.URL + "/"
	req, err := http.NewRequest(http.MethodGet, rootURL, nil)
	s.Require().NoError(err)
	res, err := client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	response, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	s.Require().NoError(err)
	s.Require().Equal([]byte("Hello, world\n"), response)
}

func (s *SuiteTest) TestHelloWorld() {
	client := &http.Client{}
	rootURL := s.ts.URL + "/hello-world"
	req, err := http.NewRequest(http.MethodGet, rootURL, nil) //nolint:noctx
	s.Require().NoError(err)
	res, err := client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	response, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	s.Require().NoError(err)
	s.Require().Equal([]byte("Hello, world\n"), response)
}

func (s *SuiteTest) TestCreateEvent() {
	s.Run("successful", func() {
		client := &http.Client{}

		newEvent := event.Event{
			Title:       "test",
			DateStart:   time.Now().Add(1 * time.Hour),
			DateEnd:     time.Now().Add(3 * time.Hour),
			Description: "test event",
			UserID:      "test user id",
		}
		createEventURL := s.ts.URL + "/events"
		reqBody, err := json.Marshal(newEvent)
		s.Require().NoError(err)

		req, err := http.NewRequest(http.MethodPost, createEventURL, bytes.NewBuffer(reqBody))
		s.Require().NoError(err)
		res, err := client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, res.StatusCode)

		response, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		s.Require().NoError(err)

		responseJSON := CreatedEvent{}
		err = json.Unmarshal(response, &responseJSON)
		s.Require().NoError(err)
		s.Require().NotZero(responseJSON.ID)

		getEventByIDUrl := s.ts.URL + "/events/" + fmt.Sprint(responseJSON.ID)
		req, err = http.NewRequest(http.MethodGet, getEventByIDUrl, nil) //nolint:noctx
		s.Require().NoError(err)
		res, err = client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, res.StatusCode)

		response, err = io.ReadAll(res.Body)
		defer res.Body.Close()
		s.Require().NoError(err)

		evtResponseJSON := event.Event{}
		err = json.Unmarshal(response, &evtResponseJSON)
		s.Require().NoError(err)
		s.Require().Equal(newEvent.Title, evtResponseJSON.Title)
		s.Require().Equal(newEvent.Description, evtResponseJSON.Description)
	})

	s.Run("bad request", func() {
		client := &http.Client{}

		newEvent := struct { // wrong date_start format + missing user_id
			Title       string
			DateStart   string // wrong format
			DateEnd     time.Time
			Description string
		}{
			Title:       "test",
			DateStart:   "2022-10-12",
			DateEnd:     time.Now().Add(3 * time.Hour),
			Description: "test event",
		}
		createEventURL := s.ts.URL + "/events"
		reqBody, err := json.Marshal(newEvent)
		s.Require().NoError(err)

		req, err := http.NewRequest(http.MethodPost, createEventURL, bytes.NewBuffer(reqBody)) //nolint:noctx
		s.Require().NoError(err)
		res, err := client.Do(req)
		s.Require().NoError(err)
		defer res.Body.Close()
		s.Require().Equal(http.StatusBadRequest, res.StatusCode)
	})
}

func (s *SuiteTest) TestUpdateEvent() {
	s.Run("successful", func() {
		client := &http.Client{}

		newEvent := event.Event{
			Title:       "new event",
			DateStart:   time.Now().Add(1 * time.Hour),
			DateEnd:     time.Now().Add(3 * time.Hour),
			Description: "test new event",
			UserID:      "test user id",
		}
		reqBody, err := json.Marshal(newEvent)
		s.Require().NoError(err)
		createURL := s.ts.URL + "/events"
		req, err := http.NewRequest(http.MethodPost, createURL, bytes.NewBuffer(reqBody))
		s.Require().NoError(err)
		res, err := client.Do(req)
		s.Require().NoError(err)
		response, err := io.ReadAll(res.Body)
		s.Require().NoError(err)
		defer res.Body.Close()
		responseJSON := CreatedEvent{}
		json.Unmarshal(response, &responseJSON)
		s.Require().NotZero(responseJSON.ID)

		updatedEvent := newEvent
		updatedEvent.Title = "updated event"
		updateURL := s.ts.URL + "/events/" + fmt.Sprint(responseJSON.ID)
		reqBody, err = json.Marshal(updatedEvent)
		s.Require().NoError(err)
		req, err = http.NewRequest(http.MethodPut, updateURL, bytes.NewBuffer(reqBody)) //nolint:noctx
		s.Require().NoError(err)
		res, err = client.Do(req)
		s.Require().NoError(err)
		defer res.Body.Close()
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, res.StatusCode)

		getEventByIDUrl := s.ts.URL + "/events/" + fmt.Sprint(responseJSON.ID)
		req, err = http.NewRequest(http.MethodGet, getEventByIDUrl, nil) //nolint:noctx
		s.Require().NoError(err)
		res, err = client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, res.StatusCode)

		response, err = io.ReadAll(res.Body)
		defer res.Body.Close()
		s.Require().NoError(err)

		evtResponseJSON := event.Event{}
		err = json.Unmarshal(response, &evtResponseJSON)
		s.Require().NoError(err)
		s.Require().Equal(updatedEvent.Title, evtResponseJSON.Title)
		s.Require().Equal(updatedEvent.Description, evtResponseJSON.Description)
	})

	s.Run("bad request", func() {
		client := &http.Client{}

		newEvent := event.Event{
			Title:       "new event",
			DateStart:   time.Now().Add(1 * time.Hour),
			DateEnd:     time.Now().Add(3 * time.Hour),
			Description: "test new event",
			UserID:      "test user id",
		}
		reqBody, err := json.Marshal(newEvent)
		s.Require().NoError(err)
		createURL := s.ts.URL + "/events"
		req, err := http.NewRequest(http.MethodPost, createURL, bytes.NewBuffer(reqBody))
		s.Require().NoError(err)
		res, err := client.Do(req)
		s.Require().NoError(err)
		response, err := io.ReadAll(res.Body)
		s.Require().NoError(err)
		defer res.Body.Close()
		responseJSON := CreatedEvent{}
		json.Unmarshal(response, &responseJSON)
		s.Require().NotZero(responseJSON.ID)

		updatedEvent := newEvent
		updatedEvent.Title = "updated event"
		updateURL := s.ts.URL + "/events/0" // wrong id
		reqBody, err = json.Marshal(updatedEvent)
		s.Require().NoError(err)
		req, err = http.NewRequest(http.MethodPut, updateURL, bytes.NewBuffer(reqBody)) //nolint:noctx
		s.Require().NoError(err)
		res, err = client.Do(req)
		s.Require().NoError(err)
		defer res.Body.Close()
		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, res.StatusCode)
	})
}

func (s *SuiteTest) TestDeleteEvent() {
	s.Run("successful", func() {
		client := &http.Client{}

		newEvent := event.Event{
			Title:       "new event",
			DateStart:   time.Now().Add(1 * time.Hour),
			DateEnd:     time.Now().Add(3 * time.Hour),
			Description: "test new event",
			UserID:      "test user id",
		}
		reqBody, err := json.Marshal(newEvent)
		s.Require().NoError(err)
		createURL := s.ts.URL + "/events"
		req, err := http.NewRequest(http.MethodPost, createURL, bytes.NewBuffer(reqBody)) //nolint:noctx
		s.Require().NoError(err)
		res, err := client.Do(req)
		s.Require().NoError(err)
		response, err := io.ReadAll(res.Body)
		s.Require().NoError(err)
		defer res.Body.Close()
		responseJSON := CreatedEvent{}
		json.Unmarshal(response, &responseJSON)
		s.Require().NotZero(responseJSON.ID)

		deleteEventURL := s.ts.URL + "/events/" + fmt.Sprint(responseJSON.ID)
		req, err = http.NewRequest(http.MethodDelete, deleteEventURL, nil) //nolint:noctx
		s.Require().NoError(err)
		res, err = client.Do(req)
		s.Require().NoError(err)
		defer res.Body.Close()
		s.Require().Equal(http.StatusOK, res.StatusCode)

		getEventByIDUrl := s.ts.URL + "/events/" + fmt.Sprint(responseJSON.ID)
		req, err = http.NewRequest(http.MethodGet, getEventByIDUrl, nil) //nolint:noctx
		s.Require().NoError(err)
		res, err = client.Do(req)
		s.Require().NoError(err)
		response, err = io.ReadAll(res.Body)
		s.Require().NoError(err)
		defer res.Body.Close()
		evtResponseJSON := event.Event{}
		err = json.Unmarshal(response, &evtResponseJSON)
		s.Require().NoError(err)
		s.Require().Equal(true, evtResponseJSON.Deleted)
	})

	s.Run("bad request", func() {
		client := &http.Client{}

		deleteEventURL := s.ts.URL + "/events/0"                            // no event
		req, err := http.NewRequest(http.MethodDelete, deleteEventURL, nil) //nolint:noctx
		s.Require().NoError(err)
		res, err := client.Do(req)
		s.Require().NoError(err)
		defer res.Body.Close()
		s.Require().Equal(http.StatusBadRequest, res.StatusCode)
	})
}

func (s *SuiteTest) TestGetEventList() {
	s.Run("empty list", func() {
		client := &http.Client{}

		getEventListURL := s.ts.URL + "/events/date/" + time.Now().Format("2006-01-02")
		req, err := http.NewRequest(http.MethodGet, getEventListURL, nil) //nolint:noctx
		s.Require().NoError(err)
		res, err := client.Do(req)
		s.Require().NoError(err)
		response, err := io.ReadAll(res.Body)
		s.Require().NoError(err)
		defer res.Body.Close()
		responseJSON := EventList{}
		err = json.Unmarshal(response, &responseJSON)
		s.Require().NoError(err)
		s.Require().Equal(0, len(responseJSON.List))
	})

	s.Run("correct lists (date, week, month)", func() {
		client := &http.Client{}

		// Create events
		createEventURL := s.ts.URL + "/events"
		week := 7 * 24 * time.Hour
		multiplier := 3

		for i := 0; i < 3; i++ {
			start := time.Now().Add(time.Duration(i*multiplier) * week) // *3 week to separate
			newEvent := event.Event{
				Title:       "test_" + fmt.Sprint(i),
				DateStart:   start,
				DateEnd:     start.Add(1 * time.Hour),
				Description: "test event",
				UserID:      "test user id",
			}
			reqBody, err := json.Marshal(newEvent)
			s.Require().NoError(err)
			req, err := http.NewRequest(http.MethodPost, createEventURL, bytes.NewBuffer(reqBody)) //nolint:noctx
			s.Require().NoError(err)
			res, err := client.Do(req)
			s.Require().NoError(err)
			defer res.Body.Close()
			s.Require().Equal(http.StatusOK, res.StatusCode)
		}

		check := func(url string, title string) {
			req, err := http.NewRequest(http.MethodGet, url, nil) //nolint:noctx
			s.Require().NoError(err)
			res, err := client.Do(req)
			s.Require().NoError(err)
			response, err := io.ReadAll(res.Body)
			s.Require().NoError(err)
			defer res.Body.Close()
			responseJSON := EventList{}
			err = json.Unmarshal(response, &responseJSON)
			s.Require().NoError(err)
			s.Require().Equal(1, len(responseJSON.List))
			s.Require().Equal(title, responseJSON.List[0].Title)
		}

		date := time.Now().Format("2006-01-02")
		getEventListByDayURL := s.ts.URL + "/events/date/" + date
		check(getEventListByDayURL, "test_0")

		date = time.Now().Add(1 * time.Duration(multiplier) * week).Format("2006-01-02")
		getEventListByWeekURL := s.ts.URL + "/events/week/" + date
		check(getEventListByWeekURL, "test_1")

		date = time.Now().Add(2 * time.Duration(multiplier) * week).Format("2006-01-02")
		getEventListByMonthURL := s.ts.URL + "/events/month/" + date
		check(getEventListByMonthURL, "test_2")
	})
}
