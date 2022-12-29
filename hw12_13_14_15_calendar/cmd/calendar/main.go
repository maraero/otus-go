package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
	es "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/service"
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

	ctx, cancel := context.WithCancel(context.Background())
	go watchSignals(cancel)
	defer cancel()

	strg := storage.New(ctx, logger, config.Storage)
	eventService := es.New(strg)
	calendar := app.New(eventService, logger)

	httpServer := serverhttp.New(calendar, config.Server)
	err = httpServer.Start()
	if err != nil {
		logger.Fatal("can not start http server", err)
	}

	grpcServer := servergrpc.New(calendar, config.Server)
	err = grpcServer.Start()
	if err != nil {
		logger.Fatal("can not start grpc server", err)
	}

	defer shutDown(strg, httpServer, grpcServer, logger)
	logger.Info("calendar is running...")
	<-ctx.Done()
}

func watchSignals(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-signals
	cancel()
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
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Stop(ctx); err != nil {
			logger.Error("failed to stop http server: " + err.Error())
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := grpcServer.Stop(); err != nil {
			logger.Error("failed to stop grpc server: " + err.Error())
		}
	}()

	wg.Wait()
	logger.Info("calendar stopped")
}
