package server_http

import (
	"bytes"
	"encoding/json"
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
	res, err := http.Get(s.ts.URL + "/")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	response, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	s.Require().NoError(err)
	s.Require().Equal([]byte("Hello, world\n"), response)
}

func (s *SuiteTest) TestHelloWorld() {
	res, err := http.Get(s.ts.URL + "/hello-world")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	response, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	s.Require().NoError(err)
	s.Require().Equal([]byte("Hello, world\n"), response)
}

func (s *SuiteTest) TestCreateEvent() {
	newEvent := event.Event{
		Title:      "test",
		DateStart:  time.Now().Add(1 * time.Hour),
		DateEnd:    time.Now().Add(3 * time.Hour),
		Descripion: "test event",
		UserID:     "test user id",
	}
	reqBody, err := json.Marshal(newEvent)
	s.Require().NoError(err)
	res, err := http.Post(s.ts.URL+"/events/", "application/json", bytes.NewBuffer(reqBody))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	response, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	s.Require().NoError(err)
	responseJson := CreatedEvent{}
	err = json.Unmarshal(response, &responseJson)
	s.Require().NoError(err)
	s.Require().NotZero(responseJson.ID)
}
