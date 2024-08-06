package sub

import (
	"avito-backend-bootcamp/pkg/utils/sl"
	"context"
	"log/slog"
)

type SubscriberRepository interface {
	SaveSubscritpion(ctx context.Context, houseID int64, email string) error
}

type Service struct {
	log        *slog.Logger
	repository SubscriberRepository
}

func New(log *slog.Logger, repository SubscriberRepository) *Service {
	return &Service{
		log:        log,
		repository: repository,
	}
}

func (s *Service) CreateSubscription(ctx context.Context, houseID int64, email string) error {
	const op = "subscription.CreateSubscription"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
		slog.Int64("house_id", houseID),
	)

	err := s.repository.SaveSubscritpion(ctx, houseID, email)
	if err != nil {
		log.Error("failed to save subscription", sl.Err(err))
		return err
	}

	return nil
}
