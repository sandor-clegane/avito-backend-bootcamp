package postgres

import (
	"avito-backend-bootcamp/internal/model"
	dbUtil "avito-backend-bootcamp/pkg/utils/db"

	"context"
	"database/sql"
)

// SaveHouse saves a new house to the database.
func (r *Repository) SaveHouse(ctx context.Context, address, developer string, year int64) (*model.House, error) {
	// Prepare the query to insert the house
	query :=
		`INSERT INTO houses (address, developer, year)
	 VALUES (?, ?, ?)`

	// Insert the house using the prepared query
	result, err := r.getter.DefaultTrOrDB(ctx, r.db).
		ExecContext(ctx, query, address, dbUtil.NewNullString(developer), year)
	if err != nil {
		return nil, err
	}

	// Get the last inserted ID
	houseID, err := result.LastInsertId()
	if err != nil {
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
		`SELECT *
	 FROM houses
	 WHERE id = ?`

	// Fetch the house using the prepared query
	var house model.House
	err := r.getter.DefaultTrOrDB(ctx, r.db).
		GetContext(ctx, &house, query, id)
	if err != nil {
		return nil, err
	}

	return &house, nil
}
