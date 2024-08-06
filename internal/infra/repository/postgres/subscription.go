package postgres

import (
	"avito-backend-bootcamp/internal/model"
	"context"
)

func (r *Repository) SubsciptionListByHouseID(
	ctx context.Context,
	houseID int64,
) ([]*model.Subscription, error) {
	return nil, nil
}

func (r *Repository) SaveSubscritpion(
	ctx context.Context,
	houseID int64,
	email string,
) error {
	return nil
}
