package server_http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
)

func New(logger *logger.Log, app *app.App, c config.Server) *Server {
	s := &Server{
		addr:   net.JoinHostPort(c.Host, c.Port),
		app:    app,
		logger: logger,
		router: mux.NewRouter(),
	}
	s.configureRouter()
	return s
}

func (s *Server) configureRouter() {
	s.router.Use(s.loggingMiddleware)
	s.router.HandleFunc("/", s.homeHandler).Methods("GET")
	s.router.HandleFunc("/hello-world", s.homeHandler).Methods("GET")
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

func requestAddr(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
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
