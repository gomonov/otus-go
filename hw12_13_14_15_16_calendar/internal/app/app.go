package app

import (
	"time"

	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/domain"
	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
}

type Storage interface {
	Event() storage.EventRepository
}

type EventRepository interface {
	Create(e *domain.Event) error
	Update(id int, e *domain.Event) error
	Delete(id int) error
	Get(id int) (domain.Event, error)
	ListByDay(date time.Time) ([]domain.Event, error)
	ListByWeek(date time.Time) ([]domain.Event, error)
	ListByMonth(date time.Time) ([]domain.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) GetEvent(id int) (domain.Event, error) {
	return a.storage.Event().Get(id)
}

func (a *App) UpdateEvent(id int, event *domain.Event) error {
	return a.storage.Event().Update(id, event)
}

func (a *App) CreateEvent(event *domain.Event) error {
	return a.storage.Event().Create(event)
}

func (a *App) ListByDay(date time.Time) ([]domain.Event, error) {
	return a.storage.Event().ListByDay(date)
}

func (a *App) ListByWeek(date time.Time) ([]domain.Event, error) {
	return a.storage.Event().ListByWeek(date)
}

func (a *App) ListByMonth(date time.Time) ([]domain.Event, error) {
	return a.storage.Event().ListByMonth(date)
}

func (a *App) DeleteEvent(id int) error {
	return a.storage.Event().Delete(id)
}
