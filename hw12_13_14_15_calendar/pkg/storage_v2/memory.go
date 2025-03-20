package storage_v2

import (
	"context"
	"github.com/google/uuid"
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
