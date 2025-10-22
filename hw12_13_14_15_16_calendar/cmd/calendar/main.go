package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/app"
	config2 "github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/server/http"
	storage2 "github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/migrations"
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

	config, err := config2.LoadConfig(configFile)
	if err != nil {
		panic(err)
	}

	logg, err := logger.New(config.Logger)
	if err != nil {
		panic(err)
	}

	if err = migrations.AutoMigrate(logg, config); err != nil {
		panic(err)
	}

	var storage storage2.Storage

	storageType := config.Storage.StorageType

	switch storageType {
	case "postgres":
		storage, err = sqlstorage.NewStorage(config.Storage.Dsn)
		if err != nil {
			log.Fatal("Failed to connect to PostgreSQL:", err)
		}
		defer storage.Close()

		logg.Info("Using PostgresSQL storage")
	default:
		storage = memorystorage.NewStorage()
		logg.Info("Using in-memory storage")
	}

	calendar := app.New(logg, storage)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	context.AfterFunc(ctx, func() {
		logg.Info("calendar is stopping...")
	})

	httpServer := internalhttp.NewServer(logg, calendar, internalhttp.ServerConf{
		Host: config.Server.Host,
		Port: config.Server.Port,
	})
	go func() {
		logg.Info("HTTP server starting...")
		if err := httpServer.Start(ctx); err != nil {
			logg.Error("failed to start HTTP server: " + err.Error())
		}
	}()
	defer func() {
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer stopCancel()
		if err := httpServer.Stop(stopCtx); err != nil {
			logg.Error("failed to stop HTTP server: " + err.Error())
		}
	}()

	grpcServer := internalgrpc.NewServer(logg, calendar, internalgrpc.GRPCConf{
		Host: config.GRPC.Host,
		Port: config.GRPC.Port,
	})
	go func() {
		logg.Info("gRPC server starting...")
		if err := grpcServer.Start(ctx); err != nil {
			logg.Error("failed to start gRPC server: " + err.Error())
		}
	}()
	defer func() {
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer stopCancel()
		if err := grpcServer.Stop(stopCtx); err != nil {
			logg.Error("failed to stop gRPC server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")
	<-ctx.Done()
}
