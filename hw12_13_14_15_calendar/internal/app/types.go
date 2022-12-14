package app

import (
	es "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/events-service"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
)

type App struct {
	EventsService *es.EventsService
	Logger        *logger.Log
}
