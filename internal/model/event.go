package model

type Event struct {
	ID      int64
	Type    EventType
	Payload string
}
