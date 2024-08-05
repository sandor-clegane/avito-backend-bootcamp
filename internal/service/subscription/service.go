package sub

import (
	"avito-backend-bootcamp/internal/model"
	"context"
	"log/slog"
)

type Service struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Service {
	return &Service{
		log: log,
	}
}

func (s *Service) CreateSubscription(ctx context.Context, houseID int64, email string) error {
	// insert subscription to repo
	// think how to check email??
	return nil
}

func (s *Service) GetSubsciberListByHouseID(ctx context.Context, houseID int64) ([]*model.Subscription, error) {
	// get list from repo
	return nil, nil
}
