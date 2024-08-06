package postgres

import (
	"avito-backend-bootcamp/internal/model"
	"context"

	"github.com/google/uuid"
)

func (r *Repository) SaveUser(
	ctx context.Context,
	email, password string,
	role model.UserType,
) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}

func (r *Repository) GetUser(
	ctx context.Context,
	ID uuid.UUID,
) (*model.User, error) {
	return nil, nil
}
