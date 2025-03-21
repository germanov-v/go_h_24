package storage

import (
	"context"
)

type Store interface {
	CreateEvent(event EventItem, ctx context.Context) error
}
