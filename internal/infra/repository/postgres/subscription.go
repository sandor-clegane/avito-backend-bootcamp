package postgres

import (
	"avito-backend-bootcamp/internal/model"
	"context"
)

// SubsciptionListByHouseID retrieves a list of subscriptions associated with a given house ID.
func (r *Repository) SubsciptionListByHouseID(ctx context.Context, houseID int64) ([]*model.Subscription, error) {
	// Prepare the query to fetch subscriptions for a specific house
	query :=
		"SELECT * " +
			"FROM subscriptions " +
			"WHERE house_id = $1"

	// Fetch the subscriptions using the prepared query
	var subscriptions []*model.Subscription
	err := r.getter.DefaultTrOrDB(ctx, r.db).
		SelectContext(ctx, &subscriptions, query, houseID)
	if err != nil {
		return nil, PostgresErrorTransform(err)
	}

	return subscriptions, nil
}

// SaveSubscritpion saves a new subscription for a given house ID and email.
func (r *Repository) SaveSubscritpion(ctx context.Context, houseID int64, email string) error {
	// Prepare the query to insert the subscription
	query :=
		"INSERT INTO subscriptions (house_id, email) " +
			"VALUES ($1, $2)"

	// Insert the subscription using the prepared query
	_, err := r.getter.DefaultTrOrDB(ctx, r.db).
		ExecContext(ctx, query, houseID, email)
	if err != nil {
		return PostgresErrorTransform(err)
	}

	return nil
}
