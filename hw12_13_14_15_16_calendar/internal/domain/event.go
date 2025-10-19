package domain

import (
	"errors"
	"time"
)

type Event struct {
	ID           int           `json:"id"`
	Title        string        `json:"title"`
	EventTime    time.Time     `json:"-"`
	Duration     time.Duration `json:"-"`
	Description  string        `json:"description"`
	UserID       int           `json:"-"`
	TimeToNotify time.Time     `json:"-"`
}

func (e *Event) GetEndTime() time.Time {
	return e.EventTime.Add(e.Duration)
}

func (e *Event) Validate() error {
	if e.Title == "" {
		return ErrEmptyTitle
	}
	if e.EventTime.IsZero() {
		return ErrInvalidEventTime
	}
	if e.Duration <= 0 {
		return ErrInvalidDuration
	}
	return nil
}

var (
	ErrEmptyTitle       = errors.New("event title cannot be empty")
	ErrInvalidEventTime = errors.New("event time must be set")
	ErrInvalidDuration  = errors.New("event duration must be positive")
	ErrEventNotFound    = errors.New("event not found")
)
