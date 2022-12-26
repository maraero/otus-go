package server_grpc

import (
	"net"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
	"google.golang.org/grpc"

	esgg "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/grpc-gen"
)

func New(app *app.App, c config.Server) *Server {
	s := &Server{
		addr:   net.JoinHostPort(c.Host, c.GrpcPort),
		app:    app,
		logger: app.Logger,
	}
	return s
}

func (s *Server) Start() error {
	lsnr, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.srv = grpc.NewServer(grpc.UnaryInterceptor(loggerInterceptor(*s.logger)))
	esgg.RegisterEventServiceServer(s.srv, NewService(s.app))
	s.logger.Info("starting grpc server on ", s.addr)
	return s.srv.Serve(lsnr)
}

func (s *Server) Stop() error {
	s.srv.GracefulStop()
	return nil
}
