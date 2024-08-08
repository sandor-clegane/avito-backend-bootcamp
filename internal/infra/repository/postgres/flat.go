package postgres

import (
	repo "avito-backend-bootcamp/internal/infra/repository"
	"avito-backend-bootcamp/internal/model"
	"context"
	"database/sql"
)

// GetFlat retrieves a flat by its ID from the database.
func (r *Repository) GetFlat(ctx context.Context, ID int64) (*model.Flat, error) {
	query :=
		"SELECT * " +
			"FROM flats " +
			"WHERE id = $1"

	var flat model.Flat
	err := r.getter.DefaultTrOrDB(ctx, r.db).
		GetContext(ctx, &flat, query, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repo.ErrNotFound
		}
		return nil, err
	}

	return &flat, nil
}

// SaveFlat saves a new flat to the database.
func (r *Repository) SaveFlat(ctx context.Context, houseID, price, rooms int64) (*model.Flat, error) {
	query :=
		"INSERT INTO flats (house_id, price, rooms) " +
			"VALUES ($1, $2, $3) " +
			"RETURNING *"

	var flat model.Flat
	err := r.getter.DefaultTrOrDB(ctx, r.db).
		GetContext(ctx, &flat, query, houseID, price, rooms)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repo.ErrConstraintViolation
		}
		return nil, err
	}

	return &flat, nil
}

// UpdateFlat updates an existing flat in the database.
func (r *Repository) UpdateFlat(ctx context.Context, flat *model.Flat) (*model.Flat, error) {
	query :=
		"UPDATE flats " +
			"SET house_id = $1, price = $2, rooms = $3, status = $4 " +
			"WHERE id = $5 " +
			"RETURNING *"

	err := r.getter.DefaultTrOrDB(ctx, r.db).
		GetContext(ctx, flat, query, flat.HouseID, flat.Price, flat.Rooms, flat.Status, flat.ID)
	if err != nil {
		return nil, err
	}

	return flat, nil
}

// FlatListByHouseID retrieves a list of flats associated with a given house ID.
func (r *Repository) FlatListByHouseID(ctx context.Context, houseID int64) ([]*model.Flat, error) {
	query :=
		"SELECT * " +
			"FROM flats " +
			"WHERE house_id = ?"

	var flats []*model.Flat
	err := r.getter.DefaultTrOrDB(ctx, r.db).
		SelectContext(ctx, &flats, query, houseID)
	if err != nil {
		return nil, err
	}

	return flats, nil
}
