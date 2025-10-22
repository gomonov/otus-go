package internalhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/domain"
)

type CreateEventRequest struct {
	Title        string        `json:"title"`
	EventTime    time.Time     `json:"event_time"`
	Duration     time.Duration `json:"duration"`
	Description  string        `json:"description"`
	UserID       int           `json:"user_id"`
	TimeToNotify time.Time     `json:"time_to_notify"`
}

type UpdateEventRequest struct {
	Title        string        `json:"title"`
	EventTime    time.Time     `json:"event_time"`
	Duration     time.Duration `json:"duration"`
	Description  string        `json:"description"`
	UserID       int           `json:"user_id"`
	TimeToNotify time.Time     `json:"time_to_notify"`
}

type EventResponse struct {
	ID           int           `json:"id"`
	Title        string        `json:"title"`
	EventTime    time.Time     `json:"event_time"`
	Duration     time.Duration `json:"duration"`
	Description  string        `json:"description"`
	UserID       int           `json:"user_id"`
	TimeToNotify time.Time     `json:"time_to_notify"`
}

type EventsListResponse struct {
	Events []EventResponse `json:"events"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (s *Server) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.rootHandler)
	mux.HandleFunc("/events", s.eventsHandler)
	mux.HandleFunc("/events/", s.eventByIDHandler)
	mux.HandleFunc("/events/day/", s.eventsByDayHandler)
	mux.HandleFunc("/events/week/", s.eventsByWeekHandler)
	mux.HandleFunc("/events/month/", s.eventsByMonthHandler)

	return mux
}

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	apiInfo := map[string]interface{}{
		"service": "Calendar API",
		"version": "1.0",
		"endpoints": []map[string]string{
			{
				"method":      "POST",
				"path":        "/events",
				"description": "Create new event",
			},
			{
				"method":      "GET",
				"path":        "/events/{id}",
				"description": "Get event by ID",
			},
			{
				"method":      "PUT",
				"path":        "/events/{id}",
				"description": "Update event",
			},
			{
				"method":      "DELETE",
				"path":        "/events/{id}",
				"description": "Delete event",
			},
			{
				"method":      "GET",
				"path":        "/events/day/{date}",
				"description": "Get events for specific day (date format: YYYY-MM-DD)",
			},
			{
				"method":      "GET",
				"path":        "/events/week/{date}",
				"description": "Get events for specific week (date format: YYYY-MM-DD)",
			},
			{
				"method":      "GET",
				"path":        "/events/month/{date}",
				"description": "Get events for specific month (date format: YYYY-MM-DD)",
			},
		},
	}

	json.NewEncoder(w).Encode(apiInfo)
}

func (s *Server) eventsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.createEventHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) eventByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getEventHandler(w, r)
	case http.MethodPut:
		s.updateEventHandler(w, r)
	case http.MethodDelete:
		s.deleteEventHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) createEventHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event := &domain.Event{
		Title:        req.Title,
		EventTime:    req.EventTime,
		Duration:     req.Duration,
		Description:  req.Description,
		UserID:       req.UserID,
		TimeToNotify: req.TimeToNotify,
	}

	if err := s.app.CreateEvent(event); err != nil {
		s.sendError(w, fmt.Sprintf("Failed to create event: %v", err), http.StatusInternalServerError)
		return
	}

	response := EventResponse{
		ID:           event.ID,
		Title:        event.Title,
		EventTime:    event.EventTime,
		Duration:     event.Duration,
		Description:  event.Description,
		UserID:       event.UserID,
		TimeToNotify: event.TimeToNotify,
	}

	s.sendJSON(w, response, http.StatusCreated)
}

func (s *Server) getEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := s.extractIDFromPath(r)
	if err != nil {
		s.sendError(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	event, err := s.app.GetEvent(id)
	if err != nil {
		if err == domain.ErrEventNotFound {
			s.sendError(w, "Event not found", http.StatusNotFound)
		} else {
			s.sendError(w, fmt.Sprintf("Failed to get event: %v", err), http.StatusInternalServerError)
		}
		return
	}

	response := EventResponse{
		ID:           event.ID,
		Title:        event.Title,
		EventTime:    event.EventTime,
		Duration:     event.Duration,
		Description:  event.Description,
		UserID:       event.UserID,
		TimeToNotify: event.TimeToNotify,
	}

	s.sendJSON(w, response, http.StatusOK)
}

func (s *Server) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := s.extractIDFromPath(r)
	if err != nil {
		s.sendError(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	var req UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event := &domain.Event{
		Title:        req.Title,
		EventTime:    req.EventTime,
		Duration:     req.Duration,
		Description:  req.Description,
		UserID:       req.UserID,
		TimeToNotify: req.TimeToNotify,
	}

	if err := s.app.UpdateEvent(id, event); err != nil {
		if err == domain.ErrEventNotFound {
			s.sendError(w, "Event not found", http.StatusNotFound)
		} else {
			s.sendError(w, fmt.Sprintf("Failed to update event: %v", err), http.StatusInternalServerError)
		}
		return
	}

	s.sendJSON(w, map[string]string{"message": "Event updated successfully"}, http.StatusOK)
}

func (s *Server) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := s.extractIDFromPath(r)
	if err != nil {
		s.sendError(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	if err := s.app.DeleteEvent(id); err != nil {
		if err == domain.ErrEventNotFound {
			s.sendError(w, "Event not found", http.StatusNotFound)
		} else {
			s.sendError(w, fmt.Sprintf("Failed to delete event: %v", err), http.StatusInternalServerError)
		}
		return
	}

	s.sendJSON(w, map[string]string{"message": "Event deleted successfully"}, http.StatusOK)
}

func (s *Server) eventsByDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	date, err := s.extractDateFromPath(r)
	if err != nil {
		s.sendError(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	events, err := s.app.ListByDay(date)
	if err != nil {
		s.sendError(w, fmt.Sprintf("Failed to get events: %v", err), http.StatusInternalServerError)
		return
	}

	response := s.convertEventsToResponse(events)
	s.sendJSON(w, response, http.StatusOK)
}

func (s *Server) eventsByWeekHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	date, err := s.extractDateFromPath(r)
	if err != nil {
		s.sendError(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	events, err := s.app.ListByWeek(date)
	if err != nil {
		s.sendError(w, fmt.Sprintf("Failed to get events: %v", err), http.StatusInternalServerError)
		return
	}

	response := s.convertEventsToResponse(events)
	s.sendJSON(w, response, http.StatusOK)
}

func (s *Server) eventsByMonthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	date, err := s.extractDateFromPath(r)
	if err != nil {
		s.sendError(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	events, err := s.app.ListByMonth(date)
	if err != nil {
		s.sendError(w, fmt.Sprintf("Failed to get events: %v", err), http.StatusInternalServerError)
		return
	}

	response := s.convertEventsToResponse(events)
	s.sendJSON(w, response, http.StatusOK)
}

func (s *Server) extractIDFromPath(r *http.Request) (int, error) {
	idStr := r.URL.Path[len("/events/"):]
	return strconv.Atoi(idStr)
}

func (s *Server) extractDateFromPath(r *http.Request) (time.Time, error) {
	var dateStr string

	switch {
	case strings.HasPrefix(r.URL.Path, "/events/day/"):
		dateStr = r.URL.Path[len("/events/day/"):]
	case strings.HasPrefix(r.URL.Path, "/events/week/"):
		dateStr = r.URL.Path[len("/events/week/"):]
	case strings.HasPrefix(r.URL.Path, "/events/month/"):
		dateStr = r.URL.Path[len("/events/month/"):]
	default:
		return time.Time{}, fmt.Errorf("invalid path format")
	}

	dateStr = strings.TrimSuffix(dateStr, "/")

	return time.Parse("2006-01-02", dateStr)
}

func (s *Server) convertEventsToResponse(events []domain.Event) EventsListResponse {
	response := EventsListResponse{
		Events: make([]EventResponse, len(events)),
	}

	for i, event := range events {
		response.Events[i] = EventResponse{
			ID:           event.ID,
			Title:        event.Title,
			EventTime:    event.EventTime,
			Duration:     event.Duration,
			Description:  event.Description,
			UserID:       event.UserID,
			TimeToNotify: event.TimeToNotify,
		}
	}

	return response
}

func (s *Server) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Error(fmt.Sprintf("Failed to encode response: %v", err))
	}
}

func (s *Server) sendError(w http.ResponseWriter, message string, statusCode int) {
	s.sendJSON(w, ErrorResponse{Error: message}, statusCode)
}
