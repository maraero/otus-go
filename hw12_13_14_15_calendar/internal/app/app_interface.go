package app

import (
	eventservice "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/service"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
)

type App struct {
	logger       *logger.Log
	eventservice *eventservice.EventService
}
