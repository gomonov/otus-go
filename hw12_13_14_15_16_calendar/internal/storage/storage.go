package storage

import (
	"time"

	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/domain"
)

type Storage interface {
	Event() EventRepository
	Close() error
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
