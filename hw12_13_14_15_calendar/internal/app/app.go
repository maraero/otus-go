package app

import (
	eventservice "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/service"
)

type App struct { // TODO
}

type Logger interface { // TODO
}

type Storage interface { // TODO
}

func New(logger Logger, eventservice *eventservice.EventService) *App {
	return &App{}
}
