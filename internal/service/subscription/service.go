package sub

import (
	"avito-backend-bootcamp/internal/infra/repository"
	"avito-backend-bootcamp/pkg/utils/sl"
	"context"
	"errors"
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

var ErrInvalidSubscription = errors.New("this house does not exist or there is no user with this address")
var ErrAlreadyExists = errors.New("you already have subscription for this house")

func (s *Service) CreateSubscription(ctx context.Context, houseID int64, email string) error {
	const op = "subscription.CreateSubscription"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
		slog.Int64("house_id", houseID),
	)

	err := s.repository.SaveSubscritpion(ctx, houseID, email)
	if err != nil {
		if errors.Is(err, repository.ErrConstraintViolation) {
			return ErrInvalidSubscription
		}
		if errors.Is(err, repository.ErrAlreadyExists) {
			return ErrAlreadyExists
		}
		log.Error("failed to save subscription", sl.Err(err))
		return err
	}

	return nil
}
