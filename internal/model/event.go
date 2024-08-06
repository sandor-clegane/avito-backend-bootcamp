package model

type Event struct {
	ID      int64     `db:"id"`
	Type    EventType `db:"type"`
	Payload string    `db:"payload"`
}
