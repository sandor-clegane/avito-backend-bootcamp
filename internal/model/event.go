package model

import "time"

type Event struct {
	ID          int64     `db:"id"`
	Type        EventType `db:"type"`
	Payload     string    `db:"payload"`
	CreatedAt   time.Time `db:"created_at"`
	ProcessedAt time.Time `db:"processed_at"`
}
