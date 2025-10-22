package memorystorage

import (
	"sync"

	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/domain"
	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[int]*domain.Event
	mu     sync.RWMutex
	nextID int
}

func NewStorage() *Storage {
	return &Storage{
		events: make(map[int]*domain.Event),
		nextID: 1,
	}
}

func (s *Storage) Event() storage.EventRepository {
	return &EventRepository{
		storage: s,
	}
}

func (s *Storage) Close() error {
	return nil
}
