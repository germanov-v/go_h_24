package storage_v2

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgresStorage(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (p *PostgresStorage) CreateEvent(event EventItem, ctx context.Context) (uuid.UUID, error) {
	query := `INSERT INTO events (id, title, 
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
	query := `UPDATE events SET title=:title,
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
