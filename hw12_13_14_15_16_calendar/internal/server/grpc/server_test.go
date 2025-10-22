package internalgrpc

import (
	"context"
	"testing"
	"time"

	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/domain"
	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func createTestGRPCLogger() *logger.Logger {
	logg, _ := logger.New(config.LoggerConf{
		Level:    "error",
		FileName: "",
	})
	return logg
}

func TestGRPCCreateEvent(t *testing.T) {
	app := &MockApp{}
	logger := createTestGRPCLogger()
	server := NewServer(logger, app, GRPCConf{})

	event := &Event{
		Title:        "Test",
		EventTime:    timestamppb.Now(),
		Duration:     durationpb.New(time.Hour),
		Description:  "Desc",
		UserId:       1,
		TimeToNotify: timestamppb.Now(),
	}

	app.On("CreateEvent", mock.Anything).Return(nil)

	resp, err := server.CreateEvent(context.Background(), &CreateEventRequest{Event: event})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Test", resp.Event.Title)
	app.AssertExpectations(t)
}

func TestGRPCGetEvent(t *testing.T) {
	app := &MockApp{}
	logger := createTestGRPCLogger()
	server := NewServer(logger, app, GRPCConf{})

	testEvent := domain.Event{
		ID:           1,
		Title:        "Test",
		EventTime:    time.Now(),
		Duration:     time.Hour,
		Description:  "Desc",
		UserID:       1,
		TimeToNotify: time.Now(),
	}

	app.On("GetEvent", 1).Return(testEvent, nil)

	resp, err := server.GetEvent(context.Background(), &GetEventRequest{Id: 1})

	assert.NoError(t, err)
	assert.Equal(t, "Test", resp.Event.Title)
	app.AssertExpectations(t)
}

func TestGRPCUpdateEvent(t *testing.T) {
	app := &MockApp{}
	logger := createTestGRPCLogger()
	server := NewServer(logger, app, GRPCConf{})

	event := &Event{
		Title: "Updated",
	}

	app.On("UpdateEvent", 1, mock.Anything).Return(nil)

	resp, err := server.UpdateEvent(context.Background(), &UpdateEventRequest{
		Id:    1,
		Event: event,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	app.AssertExpectations(t)
}

func TestGRPCDeleteEvent(t *testing.T) {
	app := &MockApp{}
	logger := createTestGRPCLogger()
	server := NewServer(logger, app, GRPCConf{})

	app.On("DeleteEvent", 1).Return(nil)

	resp, err := server.DeleteEvent(context.Background(), &DeleteEventRequest{Id: 1})

	assert.NoError(t, err)
	assert.True(t, resp.Success)
	app.AssertExpectations(t)
}

func TestGRPCListByDay(t *testing.T) {
	app := &MockApp{}
	logger := createTestGRPCLogger()
	server := NewServer(logger, app, GRPCConf{})

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

	resp, err := server.ListByDay(context.Background(), &ListByDayRequest{
		Date: timestamppb.Now(),
	})

	assert.NoError(t, err)
	assert.Len(t, resp.Events, 1)
	assert.Equal(t, "Event 1", resp.Events[0].Title)
	app.AssertExpectations(t)
}

func TestGRPCErrorHandling(t *testing.T) {
	app := &MockApp{}
	logger := createTestGRPCLogger()
	server := NewServer(logger, app, GRPCConf{})

	app.On("GetEvent", 999).Return(domain.Event{}, domain.ErrEventNotFound)

	_, err := server.GetEvent(context.Background(), &GetEventRequest{Id: 999})

	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
	app.AssertExpectations(t)
}
