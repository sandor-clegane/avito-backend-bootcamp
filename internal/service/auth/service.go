package auth

import (
	"avito-backend-bootcamp/internal/model"
	"context"
	"errors"
	"log/slog"
	"math/rand"

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

func (s *Service) DummyLogin(ctx context.Context, role model.UserType) (string, error) {
	// Имитация неуспешной авторизации
	errorProbability := 0.1
	if rand.Float64() < errorProbability {
		return "", errors.New("internal error")
	}

	return "", nil
}

func (s *Service) FindUserByID(ctx context.Context, ID uuid.UUID) (*model.User, error) {
	return nil, nil
}

var ErrUserNotFound = errors.New("user not found")

func (s *Service) Login(ctx context.Context, ID uuid.UUID, password string) (string, error) {
	// TODO
	// lookup in repo

	// compare password

	// create jwt token

	return "", nil
}

func (s *Service) Register(ctx context.Context, email, password string, role model.UserType) (uuid.UUID, error) {
	// TODO
	// save to repo
	return uuid.UUID{}, nil
}
