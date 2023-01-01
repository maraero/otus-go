package serverhttp

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
)

type Router interface {
	Routes() []Route
}

type Route struct {
	Name   string
	Method string
	Path   string
	Func   http.HandlerFunc
}

type responseWriter struct {
	http.ResponseWriter
	code int
}

type Server struct {
	addr   string
	app    *app.App
	logger *logger.Log
	srv    *http.Server
	router *mux.Router
}
