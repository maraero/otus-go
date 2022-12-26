package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
	es "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/service"
	l "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
	server_grpc "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/servers/grpc"
	server_http "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/servers/http"
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

	dbConnection := newDBConnection(ctx, logger, config.Storage)
	defer func() {
		if dbConnection != nil {
			err := dbConnection.Close()
			if err != nil {
				log.Fatal("can not close database connection: %w", err)
			}
		}
	}()

	if dbConnection != nil {
		migrate(dbConnection.DB, logger)
	}

	eventService := es.New(dbConnection)
	calendar := app.New(eventService, logger)

	httpServer := server_http.New(calendar, config.Server)
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		if err := httpServer.Stop(ctx); err != nil {
			logger.Error("failed to stop http server: " + err.Error())
		}
	}()

	grpcServer := server_grpc.New(calendar, config.Server)
	go func() {
		<-ctx.Done()
		if err := grpcServer.Stop(); err != nil {
			logger.Error("failed to stop grpc server: " + err.Error())
		}
	}()

	logger.Info("calendar is running...")

	if err := httpServer.Start(); err != nil {
		logger.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
