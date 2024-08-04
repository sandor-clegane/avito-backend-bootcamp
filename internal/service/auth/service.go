package auth

import (
	"avito-backend-bootcamp/internal/model"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type Service struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Service {
	return &Service{
		log: log,
	}
}

func (s *Service) DummyLogin(ctx context.Context, Type model.UserType) string {
	return ""
}

func (s *Service) FindUserByID(ctx context.Context, ID uuid.UUID) (*model.User, error) {
	return nil, nil
}

func (s *Service) Login(ctx context.Context, ID uuid.UUID, Password string) (string, error) {
	// TODO
	// lookup in repo

	// compare password

	// create jwt token

	return "", nil
}

func (s *Service) Register(ctx context.Context, Email, Password string, Type model.UserType) (uuid.UUID, error) {
	// TODO
	// save to repo
	return uuid.UUID{}, nil
}
