package postgres

import (
	"avito-backend-bootcamp/internal/model"
	"context"
)

func (r *Repository) GetFlat(
	ctx context.Context,
	ID int64,
) (*model.Flat, error) {
	return nil, nil
}

func (r *Repository) SaveFlat(
	ctx context.Context,
	houseID, price, fooms int64,
) (*model.Flat, error) {
	return nil, nil
}

func (r *Repository) UpdateFlat(
	ctx context.Context,
	flat *model.Flat,
) (*model.Flat, error) {
	return nil, nil
}

func (r *Repository) FlatListByHouseID(
	ctx context.Context,
	houseID int64,
) ([]*model.Flat, error) {
	return nil, nil
}
