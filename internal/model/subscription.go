package model

import "time"

// Подписка пользователя на получение уведомлений о доме
type Subscription struct {
	HouseID   int64     `db:"house_id"`
	Email     string    `db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
