package app

import (
	es "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/events-service"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
)

func New(es *es.EventsService, logger *logger.Log) *App {
	return &App{EventsService: es, Logger: logger}
}
