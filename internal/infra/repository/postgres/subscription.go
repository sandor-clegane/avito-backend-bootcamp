package postgres

import (
	repo "avito-backend-bootcamp/internal/infra/repository"
	"avito-backend-bootcamp/internal/model"
	"context"
	"database/sql"
	"errors"
)

// SubsciptionListByHouseID retrieves a list of subscriptions associated with a given house ID.
func (r *Repository) SubsciptionListByHouseID(ctx context.Context, houseID int64) ([]*model.Subscription, error) {
	// Prepare the query to fetch subscriptions for a specific house
	query :=
		`SELECT *
	 FROM subscription
	 WHERE house_id = ?`

	// Fetch the subscriptions using the prepared query
	var subscriptions []*model.Subscription
	err := r.db.SelectContext(ctx, &subscriptions, query, houseID)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// SaveSubscritpion saves a new subscription for a given house ID and email.
func (r *Repository) SaveSubscritpion(ctx context.Context, houseID int64, email string) error {
	// Prepare the query to insert the subscription
	query :=
		`INSERT INTO subscription (house_id, email)
	 VALUES (?, ?)`

	// Insert the subscription using the prepared query
	_, err := r.db.ExecContext(ctx, query, houseID, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repo.ErrConstraintViolation
		}
		return err
	}

	return nil
}
