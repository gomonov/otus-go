package internalhttp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/domain"
	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockApp struct {
	mock.Mock
}

func (m *MockApp) CreateEvent(event *domain.Event) error {
	return m.Called(event).Error(0)
}

func (m *MockApp) GetEvent(id int) (domain.Event, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Event), args.Error(1)
}

func (m *MockApp) UpdateEvent(id int, event *domain.Event) error {
	return m.Called(id, event).Error(0)
}

func (m *MockApp) DeleteEvent(id int) error {
	return m.Called(id).Error(0)
}

func (m *MockApp) ListByDay(date time.Time) ([]domain.Event, error) {
	args := m.Called(date)
	return args.Get(0).([]domain.Event), args.Error(1)
}

func (m *MockApp) ListByWeek(date time.Time) ([]domain.Event, error) {
	args := m.Called(date)
	return args.Get(0).([]domain.Event), args.Error(1)
}

func (m *MockApp) ListByMonth(date time.Time) ([]domain.Event, error) {
	args := m.Called(date)
	return args.Get(0).([]domain.Event), args.Error(1)
}

func createTestLogger() *logger.Logger {
	logg, _ := logger.New(config.LoggerConf{
		Level:    "error",
		FileName: "",
	})
	return logg
}

func TestCreateEvent(t *testing.T) {
	app := &MockApp{}
	testLogger := createTestLogger()
	server := NewServer(testLogger, app, ServerConf{})

	event := map[string]interface{}{
		"title":          "Test",
		"event_time":     time.Now().Format(time.RFC3339),
		"duration":       2100000000000,
		"description":    "Desc",
		"user_id":        1,
		"time_to_notify": time.Now().Format(time.RFC3339),
	}

	app.On("CreateEvent", mock.Anything).Return(nil)

	body, _ := json.Marshal(event)
	req := httptest.NewRequest("POST", "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	server.setupRoutes().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	app.AssertExpectations(t)
}

func TestGetEvent(t *testing.T) {
	app := &MockApp{}
	testLogger := createTestLogger()
	server := NewServer(testLogger, app, ServerConf{})

	testEvent := domain.Event{
		ID:        1,
		Title:     "Test",
		EventTime: time.Now(),
		Duration:  time.Hour,
		UserID:    1,
	}

	app.On("GetEvent", 1).Return(testEvent, nil)

	req := httptest.NewRequest("GET", "/events/1", nil)
	rr := httptest.NewRecorder()

	server.setupRoutes().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response EventResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Test", response.Title)

	app.AssertExpectations(t)
}

func TestUpdateEvent(t *testing.T) {
	app := &MockApp{}
	testLogger := createTestLogger()
	server := NewServer(testLogger, app, ServerConf{})

	event := map[string]interface{}{
		"title": "Updated",
	}

	app.On("UpdateEvent", 1, mock.Anything).Return(nil)

	body, _ := json.Marshal(event)
	req := httptest.NewRequest("PUT", "/events/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	server.setupRoutes().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	app.AssertExpectations(t)
}

func TestDeleteEvent(t *testing.T) {
	app := &MockApp{}
	testLogger := createTestLogger()
	server := NewServer(testLogger, app, ServerConf{})

	app.On("DeleteEvent", 1).Return(nil)

	req := httptest.NewRequest("DELETE", "/events/1", nil)
	rr := httptest.NewRecorder()

	server.setupRoutes().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	app.AssertExpectations(t)
}

func TestListEvents(t *testing.T) {
	app := &MockApp{}
	testLogger := createTestLogger()
	server := NewServer(testLogger, app, ServerConf{})

	testEvents := []domain.Event{
		{
			ID:        1,
			Title:     "Event 1",
			EventTime: time.Now(),
			Duration:  time.Hour,
			UserID:    1,
		},
	}

	app.On("ListByDay", mock.Anything).Return(testEvents, nil)

	req := httptest.NewRequest("GET", "/events/day/2023-12-25", nil)
	rr := httptest.NewRecorder()

	server.setupRoutes().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response EventsListResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response.Events, 1)
	assert.Equal(t, "Event 1", response.Events[0].Title)

	app.AssertExpectations(t)
}

func TestGetEventNotFound(t *testing.T) {
	app := &MockApp{}
	testLogger := createTestLogger()
	server := NewServer(testLogger, app, ServerConf{})

	app.On("GetEvent", 999).Return(domain.Event{}, domain.ErrEventNotFound)

	req := httptest.NewRequest("GET", "/events/999", nil)
	rr := httptest.NewRecorder()

	server.setupRoutes().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	app.AssertExpectations(t)
}

func TestInvalidEventID(t *testing.T) {
	app := &MockApp{}
	testLogger := createTestLogger()
	server := NewServer(testLogger, app, ServerConf{})

	req := httptest.NewRequest("GET", "/events/invalid", nil)
	rr := httptest.NewRecorder()

	server.setupRoutes().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
