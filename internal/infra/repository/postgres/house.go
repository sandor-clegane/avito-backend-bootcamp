package postgres

import (
	repo "avito-backend-bootcamp/internal/infra/repository"
	"avito-backend-bootcamp/internal/model"
	dbUtil "avito-backend-bootcamp/pkg/utils/db"
	"errors"

	"context"
	"database/sql"
)

// SaveHouse saves a new house to the database.
func (r *Repository) SaveHouse(ctx context.Context, address, developer string, year int64) (*model.House, error) {
	// Prepare the query to insert the house
	query :=
		"INSERT INTO houses (address, developer, year_of_construction) " +
			"VALUES ($1, $2, $3) RETURNING id"

	// Insert the house using the prepared query
	var houseID int64
	err := r.getter.DefaultTrOrDB(ctx, r.db).
		GetContext(ctx, &houseID, query, address, dbUtil.NewNullString(developer), year)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repo.ErrConstraintViolation
		}
		return nil, err
	}

	// Create a new house object with the inserted ID
	house := &model.House{
		ID:                 houseID,
		Address:            address,
		Developer:          sql.NullString{},
		YearOfConstruction: year,
	}

	return house, nil
}

// GetHouse retrieves a house by its ID from the database.
func (r *Repository) GetHouse(ctx context.Context, id int64) (*model.House, error) {
	// Prepare the query to fetch the house by ID
	query :=
		"SELECT * " +
			"FROM houses " +
			"WHERE id = $1"

	// Fetch the house using the prepared query
	var house model.House
	err := r.getter.DefaultTrOrDB(ctx, r.db).
		GetContext(ctx, &house, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repo.ErrNotFound
		}
		return nil, err
	}

	return &house, nil
}
