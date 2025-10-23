package sqlstorage

import (
	"time"

	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/domain"
	"github.com/jmoiron/sqlx"
)

type EventRepository struct {
	db *sqlx.DB
}

type eventDB struct {
	ID           int           `db:"id"`
	Title        string        `db:"title"`
	EventTime    time.Time     `db:"event_time"`
	Duration     time.Duration `db:"duration"`
	Description  string        `db:"description"`
	UserID       int           `db:"user_id"`
	TimeToNotify time.Time     `db:"time_to_notify"`
}

func (e eventDB) toDomain() domain.Event {
	return domain.Event{
		ID:           e.ID,
		Title:        e.Title,
		EventTime:    e.EventTime,
		Duration:     e.Duration,
		Description:  e.Description,
		UserID:       e.UserID,
		TimeToNotify: e.TimeToNotify,
	}
}

func toEventDB(e domain.Event) eventDB {
	return eventDB{
		ID:           e.ID,
		Title:        e.Title,
		EventTime:    e.EventTime,
		Duration:     e.Duration,
		Description:  e.Description,
		UserID:       e.UserID,
		TimeToNotify: e.TimeToNotify,
	}
}

func (r *EventRepository) Create(e *domain.Event) error {
	query := `
        INSERT INTO events (title, event_time, duration, description, user_id, time_to_notify)
        VALUES (:title, :event_time, :duration, :description, :user_id, :time_to_notify)
        RETURNING id
    `

	eventDB := toEventDB(*e)
	rows, err := r.db.NamedQuery(query, &eventDB)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&e.ID); err != nil {
			return err
		}
	}

	return nil
}

func (r *EventRepository) Update(id int, e *domain.Event) error {
	query := `
        UPDATE events 
        SET title = :title, event_time = :event_time, duration = :duration,
            description = :description, user_id = :user_id, time_to_notify = :time_to_notify
        WHERE id = :id
    `

	eventDB := toEventDB(*e)
	eventDB.ID = id

	result, err := r.db.NamedExec(query, &eventDB)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (r *EventRepository) Delete(id int) error {
	query := `DELETE FROM events WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (r *EventRepository) Get(id int) (domain.Event, error) {
	query := `SELECT * FROM events WHERE id = $1`

	var event eventDB
	err := r.db.Get(&event, query, id)
	if err != nil {
		return domain.Event{}, domain.ErrEventNotFound
	}

	return event.toDomain(), nil
}

func (r *EventRepository) ListByDay(date time.Time) ([]domain.Event, error) {
	query := `
        SELECT * FROM events 
        WHERE event_time::date = $1::date
        ORDER BY event_time
    `

	var eventsDB []eventDB
	err := r.db.Select(&eventsDB, query, date)
	if err != nil {
		return nil, err
	}

	events := make([]domain.Event, len(eventsDB))
	for i, event := range eventsDB {
		events[i] = event.toDomain()
	}

	return events, nil
}

func (r *EventRepository) ListByWeek(date time.Time) ([]domain.Event, error) {
	startOfWeek := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfWeek := startOfWeek.Add(7 * 24 * time.Hour)

	query := `
        SELECT * FROM events 
        WHERE event_time >= $1 AND event_time < $2
        ORDER BY event_time
    `

	var eventsDB []eventDB
	err := r.db.Select(&eventsDB, query, startOfWeek, endOfWeek)
	if err != nil {
		return nil, err
	}

	events := make([]domain.Event, len(eventsDB))
	for i, event := range eventsDB {
		events[i] = event.toDomain()
	}

	return events, nil
}

func (r *EventRepository) ListByMonth(date time.Time) ([]domain.Event, error) {
	startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	startOfNextMonth := startOfMonth.AddDate(0, 1, 0)

	query := `
        SELECT * FROM events 
        WHERE event_time >= $1 AND event_time < $2
        ORDER BY event_time
    `

	var eventsDB []eventDB
	err := r.db.Select(&eventsDB, query, startOfMonth, startOfNextMonth)
	if err != nil {
		return nil, err
	}

	events := make([]domain.Event, len(eventsDB))
	for i, event := range eventsDB {
		events[i] = event.toDomain()
	}

	return events, nil
}
