package memorystorage

import (
	"testing"
	"time"

	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_CRUD(t *testing.T) {
	storage := NewStorage()
	eventRepo := storage.Event()

	event := &domain.Event{
		Title:     "Test Event",
		EventTime: time.Now().Add(24 * time.Hour),
		Duration:  2 * time.Hour,
		UserID:    1,
	}

	err := eventRepo.Create(event)
	require.NoError(t, err)
	assert.Equal(t, 1, event.ID)

	retrieved, err := eventRepo.Get(event.ID)
	require.NoError(t, err)
	assert.Equal(t, event.Title, retrieved.Title)

	updatedEvent := &domain.Event{
		Title:     "Updated Event",
		EventTime: time.Now().Add(48 * time.Hour),
		Duration:  3 * time.Hour,
		UserID:    2,
	}

	err = eventRepo.Update(event.ID, updatedEvent)
	require.NoError(t, err)

	updated, err := eventRepo.Get(event.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Event", updated.Title)

	err = eventRepo.Delete(event.ID)
	require.NoError(t, err)

	_, err = eventRepo.Get(event.ID)
	assert.ErrorIs(t, err, domain.ErrEventNotFound)
}

func TestStorage_GetNotFound(t *testing.T) {
	storage := NewStorage()
	eventRepo := storage.Event()

	_, err := eventRepo.Get(999)
	assert.ErrorIs(t, err, domain.ErrEventNotFound)
}

func TestStorage_UpdateNotFound(t *testing.T) {
	storage := NewStorage()
	eventRepo := storage.Event()

	event := &domain.Event{
		Title:     "Test Event",
		EventTime: time.Now(),
		Duration:  1 * time.Hour,
		UserID:    1,
	}

	err := eventRepo.Update(999, event)
	assert.ErrorIs(t, err, domain.ErrEventNotFound)
}

func TestStorage_DeleteNotFound(t *testing.T) {
	storage := NewStorage()
	eventRepo := storage.Event()

	err := eventRepo.Delete(999)
	assert.ErrorIs(t, err, domain.ErrEventNotFound)
}

func TestStorage_ListByPeriods(t *testing.T) {
	storage := NewStorage()
	eventRepo := storage.Event()

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	event1 := &domain.Event{
		Title:     "Today Morning",
		EventTime: today.Add(10 * time.Hour),
		Duration:  1 * time.Hour,
		UserID:    1,
	}

	event2 := &domain.Event{
		Title:     "Today Evening",
		EventTime: today.Add(18 * time.Hour),
		Duration:  2 * time.Hour,
		UserID:    2,
	}

	event3 := &domain.Event{
		Title:     "Tomorrow",
		EventTime: today.Add(26 * time.Hour),
		Duration:  1 * time.Hour,
		UserID:    3,
	}

	require.NoError(t, eventRepo.Create(event1))
	require.NoError(t, eventRepo.Create(event2))
	require.NoError(t, eventRepo.Create(event3))

	dayEvents, err := eventRepo.ListByDay(today)
	require.NoError(t, err)
	assert.Len(t, dayEvents, 2)

	weekEvents, err := eventRepo.ListByWeek(today)
	require.NoError(t, err)
	assert.Len(t, weekEvents, 3)

	monthEvents, err := eventRepo.ListByMonth(today)
	require.NoError(t, err)
	assert.Len(t, monthEvents, 3)
}

func TestStorage_ConcurrentAccess(t *testing.T) {
	storage := NewStorage()
	eventRepo := storage.Event()

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			event := &domain.Event{
				Title:     "Event",
				EventTime: time.Now().Add(time.Duration(id) * time.Hour),
				Duration:  1 * time.Hour,
				UserID:    id,
			}
			_ = eventRepo.Create(event)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	monthEvents, err := eventRepo.ListByMonth(time.Now())
	require.NoError(t, err)
	assert.Len(t, monthEvents, 10)
}
