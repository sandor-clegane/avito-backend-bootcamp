package postgres

import (
	"avito-backend-bootcamp/internal/model"
	"context"
)

func (r *Repository) SaveHouse(
	ctx context.Context,
	address, developer string,
	year int64,
) (*model.House, error) {
	return nil, nil
}

func (r *Repository) GetHouse(
	ctx context.Context,
	id int64,
) (*model.House, error) {
	return nil, nil
}
