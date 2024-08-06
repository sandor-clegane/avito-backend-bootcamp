package model

import (
	"database/sql"
	"time"
)

// Дом
type House struct {
	ID                 int64          `json:"id" db:"id"`
	Address            string         `json:"address" db:"address"`
	YearOfConstruction int64          `json:"year_of_construction" db:"year_of_construction"`
	Developer          sql.NullString `json:"developer,omitempty" db:"develope"`
	CreatedAt          time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at" db:"updated_at"`
}
