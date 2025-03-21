package storage_v2

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

// TODO: мапа внутри и мьютекс
type MemoryStorage struct{}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (p *MemoryStorage) CreateEvent(event EventItem, ctx context.Context) (uuid.UUID, error) {

	event.ID = uuid.New()

	return event.ID, nil
}

func (p *MemoryStorage) UpdateEvent(event EventItem, ctx context.Context) error {

	return nil
}

func (p *MemoryStorage) DeleteEvent(id uuid.UUID, ctx context.Context) error {
	return nil
}

func (p *MemoryStorage) GetEventsForDay(time time.Time, ctx context.Context) ([]EventItem, error) {
	return nil, errors.New("not implemented")
}

func (p *MemoryStorage) GetEventsForWeek(start time.Time, ctx context.Context) ([]EventItem, error) {

	return nil, errors.New("GetEventsForWeek is not yet implemented")
}

func (p *MemoryStorage) GetEventsForYear(start time.Time, ctx context.Context) ([]EventItem, error) {
	return nil, errors.New("GetEventsForYear is not yet implemented")
}
