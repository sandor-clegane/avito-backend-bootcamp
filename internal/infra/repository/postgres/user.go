package postgres

import (
	repo "avito-backend-bootcamp/internal/infra/repository"
	"avito-backend-bootcamp/internal/model"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

// SaveUser saves a new user to the database.
func (r *Repository) SaveUser(ctx context.Context, email, password string, role model.UserType) (uuid.UUID, error) {
	// Generate a unique UUID for the user
	userID := uuid.New()

	// Prepare the query to insert the user
	query :=
		`INSERT INTO user (id, email, password, role)
	 VALUES (?, ?, ?, ?)`

	// Insert the user using the prepared query
	_, err := r.db.ExecContext(ctx, query, userID, email, password, role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, repo.ErrConstraintViolation
		}
		return uuid.UUID{}, err
	}

	return userID, nil
}

// GetUser retrieves a user by its ID from the database.
func (r *Repository) GetUser(ctx context.Context, ID uuid.UUID) (*model.User, error) {
	// Prepare the query to fetch the user by ID
	query :=
		`SELECT *
	 FROM user
	 WHERE id = ?`

	// Fetch the user using the prepared query
	var user model.User
	err := r.db.GetContext(ctx, &user, query, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repo.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
