package postgres

import (
	repo "avito-backend-bootcamp/internal/infra/repository"
	"avito-backend-bootcamp/internal/model"
	"context"
	"database/sql"
	"errors"
)

// PublishEvent publishes a new event.
func (r *Repository) PublishEvent(ctx context.Context, eventType model.EventType, payload string) error {
	// Insert the event into the database
	_, err := r.getter.DefaultTrOrDB(ctx, r.db).ExecContext(ctx,
		"INSERT INTO events (type, payload) "+
			"VALUES ($1, $2)",
		eventType, payload)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repo.ErrConstraintViolation
		}
		return err
	}

	return nil
}

// GetNewEvent retrieves the next unprocessed event.
func (r *Repository) GetNewEvent(ctx context.Context) (*model.Event, error) {
	// Select the oldest unprocessed event
	var event model.Event
	err := r.getter.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &event,
		"SELECT * "+
			"FROM events "+
			"WHERE processed_at IS NULL "+
			"ORDER BY created_at ASC "+
			"LIMIT 1",
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repo.ErrNotFound
		}
		return nil, err
	}

	return &event, nil
}

// SetDone marks an event as processed.
func (r *Repository) SetDone(ctx context.Context, eventID int64) error {
	// Update the event to mark it as processed
	_, err := r.getter.DefaultTrOrDB(ctx, r.db).ExecContext(ctx,
		`UPDATE events
		SET processed_at = NOW()
		WHERE id = $1`,
		eventID)
	if err != nil {
		return err
	}

	return nil
}
