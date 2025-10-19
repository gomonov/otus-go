package memorystorage

import (
	"time"

	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/domain"
)

type EventRepository struct {
	storage *Storage
}

func (r *EventRepository) Create(e *domain.Event) error {
	r.storage.mu.Lock()
	defer r.storage.mu.Unlock()

	e.ID = r.storage.nextID
	r.storage.events[e.ID] = e
	r.storage.nextID++
	return nil
}

func (r *EventRepository) Update(id int, e *domain.Event) error {
	r.storage.mu.Lock()
	defer r.storage.mu.Unlock()

	if _, exists := r.storage.events[id]; !exists {
		return domain.ErrEventNotFound
	}

	e.ID = id
	r.storage.events[id] = e
	return nil
}

func (r *EventRepository) Delete(id int) error {
	r.storage.mu.Lock()
	defer r.storage.mu.Unlock()

	if _, exists := r.storage.events[id]; !exists {
		return domain.ErrEventNotFound
	}

	delete(r.storage.events, id)
	return nil
}

func (r *EventRepository) Get(id int) (domain.Event, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	event, exists := r.storage.events[id]
	if !exists {
		return domain.Event{}, domain.ErrEventNotFound
	}
	return *event, nil
}

func (r *EventRepository) ListByDay(date time.Time) ([]domain.Event, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	var events []domain.Event
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	for _, event := range r.storage.events {
		if (event.EventTime.After(startOfDay) || event.EventTime.Equal(startOfDay)) &&
			event.EventTime.Before(endOfDay) {
			events = append(events, *event)
		}
	}
	return events, nil
}

func (r *EventRepository) ListByWeek(date time.Time) ([]domain.Event, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	var events []domain.Event
	year, week := date.ISOWeek()

	for _, event := range r.storage.events {
		eventYear, eventWeek := event.EventTime.ISOWeek()
		if eventYear == year && eventWeek == week {
			events = append(events, *event)
		}
	}
	return events, nil
}

func (r *EventRepository) ListByMonth(date time.Time) ([]domain.Event, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	var events []domain.Event
	year, month := date.Year(), date.Month()

	for _, event := range r.storage.events {
		eventYear, eventMonth := event.EventTime.Year(), event.EventTime.Month()
		if eventYear == year && eventMonth == month {
			events = append(events, *event)
		}
	}
	return events, nil
}
