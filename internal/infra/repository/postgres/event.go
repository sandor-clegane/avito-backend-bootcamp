package postgres

import (
	"avito-backend-bootcamp/internal/model"
	"context"
)

func (r *Repository) PublishEvent(
	ctx context.Context,
	eventType model.EventType,
	payload string,
) error {
	return nil
}

func (r *Repository) GetNewEvent(ctx context.Context) (*model.Event, error) {
	return nil, nil
}

func (r *Repository) SetDone(ctx context.Context, eventID int64) error {
	return nil
}
