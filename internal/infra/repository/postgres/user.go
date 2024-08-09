package postgres

import (
	"avito-backend-bootcamp/internal/model"
	"context"

	"github.com/google/uuid"
)

// SaveUser saves a new user to the database.
func (r *Repository) SaveUser(ctx context.Context, email, password string, role model.UserType) (uuid.UUID, error) {
	// Generate a unique UUID for the user
	userID := uuid.New()

	// Prepare the query to insert the user
	query :=
		"INSERT INTO users (id, email, password, type) " +
			"VALUES ($1, $2, $3, $4)"

	// Insert the user using the prepared query
	_, err := r.getter.DefaultTrOrDB(ctx, r.db).
		ExecContext(ctx, query, userID, email, password, role)
	if err != nil {
		return uuid.UUID{}, PostgresErrorTransform(err)
	}

	return userID, nil
}

// GetUser retrieves a user by its ID from the database.
func (r *Repository) GetUser(ctx context.Context, ID uuid.UUID) (*model.User, error) {
	// Prepare the query to fetch the user by ID
	query :=
		"SELECT * " +
			"FROM users " +
			"WHERE id = $1"

	// Fetch the user using the prepared query
	var user model.User
	err := r.getter.DefaultTrOrDB(ctx, r.db).
		GetContext(ctx, &user, query, ID)
	if err != nil {
		return nil, PostgresErrorTransform(err)
	}

	return &user, nil
}
