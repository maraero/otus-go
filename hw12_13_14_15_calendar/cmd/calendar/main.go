package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
	er "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/events-repository"
	es "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/events-service"
	l "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
	servergrpc "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/servers/grpc"
	serverhttp "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/servers/http"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := config.New(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := l.New(config.Logger)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	strg := storage.New(ctx, logger, config.Storage)
	eventsRepository := er.New(strg)
	eventsService := es.New(eventsRepository)
	calendar := app.New(eventsService, logger)

	httpServer := serverhttp.New(calendar, config.Server)
	go func() {
		err = httpServer.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server closed:", err)
			cancel()
		}
	}()

	grpcServer := servergrpc.New(calendar, config.Server)
	go func() {
		err = grpcServer.Start()
		if err != nil {
			logger.Error("grpc server closed:", err)
			cancel()
		}
	}()

	logger.Info("calendar is running...")
	<-ctx.Done()
	shutDown(strg, httpServer, grpcServer, logger)
}

func shutDown(strg *storage.Storage, httpServer *serverhttp.Server, grpcServer *servergrpc.Server, logger *l.Log) {
	logger.Info("calendar is turning off...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if strg.Connection != nil {
			err := strg.Connection.Close()
			if err != nil {
				logger.Error("can not close database connection: %w", err)
			} else {
				logger.Info("database connection closed")
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Stop(ctx); err != nil {
			logger.Error("failed to stop http server: " + err.Error())
		} else {
			logger.Info("http server closed")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := grpcServer.Stop(); err != nil {
			logger.Error("failed to stop grpc server: " + err.Error())
		} else {
			logger.Info("grpc server closed")
		}
	}()

	wg.Wait()
	logger.Info("calendar stopped")
}
