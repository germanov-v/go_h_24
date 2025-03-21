package storage_v2

import (
	"github.com/google/uuid"
	"time"
)

type EventItem struct {
	ID           uuid.UUID      `db:"id"` // TODO: может вообще сделать поле guid_id и id сделать pk int64 nexval
	Title        string         `db:"title"`
	StartTime    time.Time      `db:"start_time"`
	Duration     time.Duration  `db:"duration"`
	Description  *string        `db:"description"`
	UserID       uuid.UUID      `db:"user_id"`
	NotifyBefore *time.Duration `db:"notify_before"`
}
