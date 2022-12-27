package server_http

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
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
