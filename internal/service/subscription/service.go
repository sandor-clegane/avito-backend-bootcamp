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
	// insert to repo
	return nil
}

func (s *Service) GetSubsciberListByHouseID(ctx context.Context, houseID int64) ([]*model.Subscription, error) {
	// get list from repo
	return nil, nil
}

func (s *Service) NotifySubscriber(ctx context.Context, subEmail, message string) error {
	// use mock client
	return nil
}
