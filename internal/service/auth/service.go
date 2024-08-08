package auth

import (
	"avito-backend-bootcamp/internal/infra/repository"
	"avito-backend-bootcamp/internal/model"
	"avito-backend-bootcamp/pkg/utils/sl"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type JWT interface {
	CreateToken(role string) (string, error)
}

type UserRepository interface {
	SaveUser(ctx context.Context, email, password string, role model.UserType) (uuid.UUID, error)
	GetUser(ctx context.Context, ID uuid.UUID) (*model.User, error)
}

type Service struct {
	repository UserRepository
	log        *slog.Logger
	jwt        JWT
}

func New(log *slog.Logger, jwt JWT, repository UserRepository) *Service {
	return &Service{
		log:        log,
		repository: repository,
		jwt:        jwt,
	}
}

func (s *Service) DummyLogin(ctx context.Context, role model.UserType) (string, error) {
	const op = "Auth.DummyLogin"

	// Имитация неуспешной авторизации
	errorProbability := 0.1
	if rand.Float64() < errorProbability {
		return "", errors.New("internal error")
	}

	token, err := s.jwt.CreateToken(string(role))
	if err != nil {
		return "", err
	}

	return token, nil
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (s *Service) Login(ctx context.Context, ID uuid.UUID, password string) (string, error) {
	const op = "Auth.Login"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to login user")

	user, err := s.repository.GetUser(ctx, ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.Warn("user not found", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Info("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	token, err := s.jwt.CreateToken(string(user.Type))
	if err != nil {
		log.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

var (
	ErrEmailAlreadyUser = errors.New("this email already used")
)

func (s *Service) Register(ctx context.Context, email, password string, role model.UserType) (uuid.UUID, error) {
	const op = "Auth.Register"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	id, err := s.repository.SaveUser(ctx, email, string(passHash), role)
	if err != nil {
		if errors.Is(err, repository.ErrConstraintViolation) {
			log.Error("email already user", sl.Err(err))
			return uuid.UUID{}, ErrEmailAlreadyUser
		}
		log.Error("failed to save user", sl.Err(err))
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
