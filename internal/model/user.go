package model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Email    string
	Password string
	Type     UserType
}
