package storage_v2

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrDateBusy      = errors.New("event date is already exist")
)

type Storage interface {
	CreateEvent(event EventItem, ctx context.Context) (uuid.UUID, error)
	UpdateEvent(event EventItem, ctx context.Context) error
	DeleteEvent(id uuid.UUID, ctx context.Context) error
}
