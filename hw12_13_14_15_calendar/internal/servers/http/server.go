package server_http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
)

func New(app *app.App, c config.Server) *Server {
	s := &Server{
		addr:   net.JoinHostPort(c.Host, c.HttpPort),
		app:    app,
		logger: app.Logger,
		router: mux.NewRouter(),
	}
	s.configureRouter()
	return s
}

func (s *Server) configureRouter() {
	s.router.Use(s.loggingMiddleware)

	s.router.HandleFunc("/", s.homeHandler).Methods("GET")
	s.router.HandleFunc("/hello-world", s.homeHandler).Methods("GET")

	eventsRouter := s.router.PathPrefix("/events").Subrouter()
	eventsRouter.HandleFunc("/", handleCreateEvent(s.app)).Methods(http.MethodPost)
	eventsRouter.HandleFunc("/", handleUpdateEvent(s.app)).Methods(http.MethodPut)
	eventsRouter.HandleFunc("/", handleDeleteEvent(s.app)).Methods(http.MethodDelete)
	eventsRouter.HandleFunc("/date/{date}", handleGetEventList(s.app, "date")).Methods(http.MethodGet)
	eventsRouter.HandleFunc("/week/{date}", handleGetEventList(s.app, "week")).Methods(http.MethodGet)
	eventsRouter.HandleFunc("/month/{date}", handleGetEventList(s.app, "month")).Methods(http.MethodGet)
}

func (s *Server) Start() error {
	s.srv = &http.Server{
		Addr:         s.addr,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := s.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server closed: %w", err)
	}
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}
	return err
}

func userAgent(r *http.Request) string {
	userAgents := r.Header["User-Agent"]
	if len(userAgents) > 0 {
		return "\"" + userAgents[0] + "\""
	}
	return ""
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (s *Server) homeHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hello, world\n"))
	if err != nil {
		s.logger.Error(fmt.Errorf("http write: %w", err))
	}
}
