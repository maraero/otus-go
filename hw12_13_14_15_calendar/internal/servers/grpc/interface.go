package server_grpc

import (
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
)

type Server struct {
	addr   string
	app    *app.App
	logger *logger.Log
	srv    *grpc.Server
}
