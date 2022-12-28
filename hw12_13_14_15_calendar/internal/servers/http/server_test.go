package server_http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func (s *SuiteTest) SetupTest() {
	config, err := config.New("../../../configs/config_test.json")
	if err != nil {
		log.Fatal("can not parse test config file", err)
	}

	logger, err := l.New(config.Logger)
	if err != nil {
		log.Fatal("can not init logger", err)
	}

	var dbConnection *sqlx.DB
	eventService := es.New(dbConnection)
	calendar := app.New(eventService, logger)
	s.ts = httptest.NewServer(New(calendar, config.Server).router)
}

func TestHttpServer(t *testing.T) {
	suite.Run(t, new(SuiteTest))
}

func (s *SuiteTest) TestRoot() {
	client := &http.Client{}
	rootUrl := s.ts.URL + "/"
	req, err := http.NewRequest(http.MethodGet, rootUrl, nil)
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
	rootUrl := s.ts.URL + "/hello-world"
	req, err := http.NewRequest(http.MethodGet, rootUrl, nil)
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
		createEventUrl := s.ts.URL + "/events"
		reqBody, err := json.Marshal(newEvent)
		s.Require().NoError(err)

		req, err := http.NewRequest(http.MethodPost, createEventUrl, bytes.NewBuffer(reqBody))
		res, err := client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, res.StatusCode)

		response, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		s.Require().NoError(err)

		responseJson := CreatedEvent{}
		err = json.Unmarshal(response, &responseJson)
		s.Require().NoError(err)
		s.Require().NotZero(responseJson.ID)

		getEventByIdUrl := s.ts.URL + "/events/" + fmt.Sprint(responseJson.ID)
		req, err = http.NewRequest(http.MethodGet, getEventByIdUrl, nil)
		res, err = client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, res.StatusCode)

		response, err = io.ReadAll(res.Body)
		defer res.Body.Close()
		s.Require().NoError(err)

		evtResponseJson := event.Event{}
		err = json.Unmarshal(response, &evtResponseJson)
		s.Require().NoError(err)
		s.Require().Equal(newEvent.Title, evtResponseJson.Title)
		s.Require().Equal(newEvent.Description, evtResponseJson.Description)
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
		createEventUrl := s.ts.URL + "/events"
		reqBody, err := json.Marshal(newEvent)
		s.Require().NoError(err)

		req, err := http.NewRequest(http.MethodPost, createEventUrl, bytes.NewBuffer(reqBody))
		res, err := client.Do(req)
		s.Require().NoError(err)
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
		createUrl := s.ts.URL + "/events"
		req, err := http.NewRequest(http.MethodPost, createUrl, bytes.NewBuffer(reqBody))
		res, err := client.Do(req)
		response, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		responseJson := CreatedEvent{}
		json.Unmarshal(response, &responseJson)
		s.Require().NotZero(responseJson.ID)

		updatedEvent := newEvent
		updatedEvent.Title = "updated event"
		updateUrl := s.ts.URL + "/events/" + fmt.Sprint(responseJson.ID)
		reqBody, err = json.Marshal(updatedEvent)
		s.Require().NoError(err)
		req, err = http.NewRequest(http.MethodPut, updateUrl, bytes.NewBuffer(reqBody))
		res, err = client.Do(req)
		response, err = io.ReadAll(res.Body)
		defer res.Body.Close()
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, res.StatusCode)

		getEventByIdUrl := s.ts.URL + "/events/" + fmt.Sprint(responseJson.ID)
		req, err = http.NewRequest(http.MethodGet, getEventByIdUrl, nil)
		res, err = client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, res.StatusCode)

		response, err = io.ReadAll(res.Body)
		defer res.Body.Close()
		s.Require().NoError(err)

		evtResponseJson := event.Event{}
		err = json.Unmarshal(response, &evtResponseJson)
		s.Require().NoError(err)
		s.Require().Equal(updatedEvent.Title, evtResponseJson.Title)
		s.Require().Equal(updatedEvent.Description, evtResponseJson.Description)
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
		createUrl := s.ts.URL + "/events"
		req, err := http.NewRequest(http.MethodPost, createUrl, bytes.NewBuffer(reqBody))
		res, err := client.Do(req)
		response, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		responseJson := CreatedEvent{}
		json.Unmarshal(response, &responseJson)
		s.Require().NotZero(responseJson.ID)

		updatedEvent := newEvent
		updatedEvent.Title = "updated event"
		updateUrl := s.ts.URL + "/events/0" // wrong id
		reqBody, err = json.Marshal(updatedEvent)
		s.Require().NoError(err)
		req, err = http.NewRequest(http.MethodPut, updateUrl, bytes.NewBuffer(reqBody))
		res, err = client.Do(req)
		response, err = io.ReadAll(res.Body)
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
		createUrl := s.ts.URL + "/events"
		req, err := http.NewRequest(http.MethodPost, createUrl, bytes.NewBuffer(reqBody))
		res, err := client.Do(req)
		response, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		responseJson := CreatedEvent{}
		json.Unmarshal(response, &responseJson)
		s.Require().NotZero(responseJson.ID)

		deleteEventUrl := s.ts.URL + "/events/" + fmt.Sprint(responseJson.ID)
		req, err = http.NewRequest(http.MethodDelete, deleteEventUrl, nil)
		res, err = client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, res.StatusCode)

		getEventByIdUrl := s.ts.URL + "/events/" + fmt.Sprint(responseJson.ID)
		req, err = http.NewRequest(http.MethodGet, getEventByIdUrl, nil)
		res, err = client.Do(req)
		response, err = io.ReadAll(res.Body)
		defer res.Body.Close()
		evtResponseJson := event.Event{}
		err = json.Unmarshal(response, &evtResponseJson)
		s.Require().NoError(err)
		s.Require().Equal(true, evtResponseJson.Deleted)
	})

	s.Run("bad request", func() {
		client := &http.Client{}

		deleteEventUrl := s.ts.URL + "/events/0" // no event
		req, err := http.NewRequest(http.MethodDelete, deleteEventUrl, nil)
		res, err := client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, res.StatusCode)
	})
}

func (s *SuiteTest) TestGetEventList() {
	s.Run("empty list", func() {
		client := &http.Client{}

		getEventListUrl := s.ts.URL + "/events/date/" + time.Now().Format("2006-01-02")
		req, err := http.NewRequest(http.MethodGet, getEventListUrl, nil)
		s.Require().NoError(err)
		res, err := client.Do(req)
		s.Require().NoError(err)
		response, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		responseJson := EventList{}
		err = json.Unmarshal(response, &responseJson)
		s.Require().NoError(err)
		s.Require().Equal(0, len(responseJson.List))
	})

	s.Run("correct lists (date, week, month)", func() {
		client := &http.Client{}

		// Create events
		createEventUrl := s.ts.URL + "/events"
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
			req, err := http.NewRequest(http.MethodPost, createEventUrl, bytes.NewBuffer(reqBody))
			res, err := client.Do(req)
			s.Require().Equal(http.StatusOK, res.StatusCode)
		}

		check := func(url string, title string) {
			req, err := http.NewRequest(http.MethodGet, url, nil)
			res, err := client.Do(req)
			response, err := io.ReadAll(res.Body)
			defer res.Body.Close()
			responseJson := EventList{}
			err = json.Unmarshal(response, &responseJson)
			s.Require().NoError(err)
			s.Require().Equal(1, len(responseJson.List))
			s.Require().Equal(title, responseJson.List[0].Title)
		}

		getEventListByDayUrl := s.ts.URL + "/events/date/" + time.Now().Format("2006-01-02")
		check(getEventListByDayUrl, "test_0")

		getEventListByWeekUrl := s.ts.URL + "/events/week/" + time.Now().Add(1*time.Duration(multiplier)*week).Format("2006-01-02")
		check(getEventListByWeekUrl, "test_1")

		getEventListByMonthUrl := s.ts.URL + "/events/month/" + time.Now().Add(2*time.Duration(multiplier)*week).Format("2006-01-02")
		check(getEventListByMonthUrl, "test_2")
	})
}
