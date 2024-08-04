package model

import "time"

// Дом
type House struct {
	ID                 int       `json:"id"`
	Address            string    `json:"address"`
	YearOfConstruction int       `json:"year_of_construction"`
	Developer          *string   `json:"developer,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
