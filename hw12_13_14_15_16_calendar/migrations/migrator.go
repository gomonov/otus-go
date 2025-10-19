package migrations

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/config"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

func AutoMigrate(logger Logger, cfg *config.Config) error {
	if !cfg.Migrations.AutoMigrate {
		logger.Info("Auto migrations disabled")
		return nil
	}

	if cfg.Storage.StorageType != "postgres" {
		logger.Info("Skipping migrations for in-memory storage")
		return nil
	}

	db, err := sql.Open("postgres", cfg.Storage.Dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	log.Println("Applying migrations...")
	if err := goose.Up(db, cfg.Migrations.Dir); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}
