package app

import (
	es "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/service"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
)

func New(es *es.EventService, logger *logger.Log) *App {
	return &App{Event_service: es, Logger: logger}
}
