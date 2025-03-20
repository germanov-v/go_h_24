package storage

import (
	"context"
	"github.com/google/uuid"
)

type Event struct {
	ID uuid.UUID
}

type Storage interface {
	CreateEvent(event Event, ctx context.Context) error
}
