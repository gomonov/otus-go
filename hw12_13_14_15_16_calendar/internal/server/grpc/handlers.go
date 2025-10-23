package internalgrpc

import (
	"context"
	"errors"

	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreateEvent(_ context.Context, req *CreateEventRequest) (*CreateEventResponse, error) {
	eventProto := req.GetEvent()
	if eventProto == nil {
		return nil, status.Error(codes.InvalidArgument, "event is required")
	}

	event := &domain.Event{
		Title:        eventProto.GetTitle(),
		EventTime:    eventProto.GetEventTime().AsTime(),
		Duration:     eventProto.GetDuration().AsDuration(),
		Description:  eventProto.GetDescription(),
		UserID:       int(eventProto.GetUserId()),
		TimeToNotify: eventProto.GetTimeToNotify().AsTime(),
	}

	if err := s.app.CreateEvent(event); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &CreateEventResponse{
		Event: s.domainEventToProto(event),
	}, nil
}

func (s *Server) UpdateEvent(_ context.Context, req *UpdateEventRequest) (*UpdateEventResponse, error) {
	eventProto := req.GetEvent()
	if eventProto == nil {
		return nil, status.Error(codes.InvalidArgument, "event is required")
	}

	event := &domain.Event{
		Title:        eventProto.GetTitle(),
		EventTime:    eventProto.GetEventTime().AsTime(),
		Duration:     eventProto.GetDuration().AsDuration(),
		Description:  eventProto.GetDescription(),
		UserID:       int(eventProto.GetUserId()),
		TimeToNotify: eventProto.GetTimeToNotify().AsTime(),
	}

	if err := s.app.UpdateEvent(int(req.GetId()), event); err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, "event not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &UpdateEventResponse{
		Event: s.domainEventToProto(event),
	}, nil
}

func (s *Server) DeleteEvent(_ context.Context, req *DeleteEventRequest) (*DeleteEventResponse, error) {
	if err := s.app.DeleteEvent(int(req.GetId())); err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, "event not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &DeleteEventResponse{
		Success: true,
	}, nil
}

func (s *Server) GetEvent(_ context.Context, req *GetEventRequest) (*GetEventResponse, error) {
	event, err := s.app.GetEvent(int(req.GetId()))
	if err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, "event not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &GetEventResponse{
		Event: s.domainEventToProto(&event),
	}, nil
}

func (s *Server) ListByDay(_ context.Context, req *ListByDayRequest) (*ListEventsResponse, error) {
	date := req.GetDate().AsTime()
	events, err := s.app.ListByDay(date)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ListEventsResponse{
		Events: s.domainEventsToProto(events),
	}, nil
}

func (s *Server) ListByWeek(_ context.Context, req *ListByWeekRequest) (*ListEventsResponse, error) {
	date := req.GetDate().AsTime()
	events, err := s.app.ListByWeek(date)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ListEventsResponse{
		Events: s.domainEventsToProto(events),
	}, nil
}

func (s *Server) ListByMonth(_ context.Context, req *ListByMonthRequest) (*ListEventsResponse, error) {
	date := req.GetDate().AsTime()
	events, err := s.app.ListByMonth(date)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ListEventsResponse{
		Events: s.domainEventsToProto(events),
	}, nil
}

func (s *Server) domainEventToProto(event *domain.Event) *Event {
	return &Event{
		Id:           int64(event.ID),
		Title:        event.Title,
		EventTime:    timestamppb.New(event.EventTime),
		Duration:     durationpb.New(event.Duration),
		Description:  event.Description,
		UserId:       int64(event.UserID),
		TimeToNotify: timestamppb.New(event.TimeToNotify),
	}
}

func (s *Server) domainEventsToProto(events []domain.Event) []*Event {
	protoEvents := make([]*Event, len(events))
	for i, event := range events {
		protoEvents[i] = s.domainEventToProto(&event)
	}
	return protoEvents
}
