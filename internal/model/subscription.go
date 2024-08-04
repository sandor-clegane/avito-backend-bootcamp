package model

import "time"

// Подписка пользователя на получение уведомлений о доме
type Subscription struct {
	HouseID   int64
	Email     string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
