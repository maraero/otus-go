package app

import (
	eventservice "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/service"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
)

func New(logger *logger.Log, eventservice *eventservice.EventService) *App {
	return &App{logger: logger, eventservice: eventservice}
}
