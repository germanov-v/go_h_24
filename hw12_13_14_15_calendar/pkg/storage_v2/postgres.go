package storage_v2

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgresStorage(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (p *PostgresStorage) CreateEvent(event EventItem, ctx context.Context) (uuid.UUID, error) {
	query := `INSERT INTO events.events (id, title, 
                    start_time, 
                    duration, 
                    description, user_id, 
                    notify_before)
			  VALUES (:id, :title, :start_time,
			          :duration, 
			          :description, :user_id,
			          :notify_before)
			  RETURNING id`

	event.ID = uuid.New()

	// TODO: решить вопрос с идентфикатором инкримент или гуид
	var id uuid.UUID
	err := p.db.GetContext(ctx, &id, query, event)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (p *PostgresStorage) UpdateEvent(event EventItem, ctx context.Context) error {
	query := `UPDATE events.events SET title=:title,
                  start_time=:start_time,
                  duration=:duration,
			  description=:description,
			  user_id=:user_id, 
			  notify_before=:notify_before
			  WHERE id=:id`

	_, err := p.db.NamedExecContext(ctx, query, event)
	return err
}

func (p *PostgresStorage) DeleteEvent(id uuid.UUID, ctx context.Context) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM events WHERE id=$1", id)
	return err
}

func (p *PostgresStorage) GetEventsForDay(time time.Time, ctx context.Context) ([]EventItem, error) {
	query := `SELECT id, title, start_time, 
				  duration, description,
				  user_id, notify_before
			  FROM events.events
              WHERE DATE(start_time)=DATE($1)
                `

	var events []EventItem

	err := p.db.SelectContext(ctx, &events, query, time)

	if err != nil {
		return nil, fmt.Errorf("GetEventsForDay: %w", err)
	}

	return events, nil
}

func (p *PostgresStorage) GetEventsForWeek(start time.Time, ctx context.Context) ([]EventItem, error) {
	since := start.AddDate(0, 0, -int(start.Weekday()))
	//.Truncate(24*time.Hour)

	till := time.Now() //since.AddDate(0, 0, 7)
	return p.GetEventsByPeriod(since, till, ctx)
}

func (p *PostgresStorage) GetEventsForYear(start time.Time, ctx context.Context) ([]EventItem, error) {
	since := start.Truncate(24*time.Hour).AddDate(1, 0, 0)
	till := time.Now()
	return p.GetEventsByPeriod(since, till, ctx)
}

func (p *PostgresStorage) GetEventsByPeriod(since time.Time, till time.Time, ctx context.Context) ([]EventItem, error) {

	query := `SELECT id, title, start_time, 
				  duration, description,
				  user_id, notify_before
			  FROM events.events
              WHERE start_time BETWEEN $1 AND $2
                `

	var events []EventItem

	err := p.db.SelectContext(ctx, &events, query, since, till)

	if err != nil {
		return nil, fmt.Errorf("GetEventsForDay: %w", err)
	}

	return events, nil
}
