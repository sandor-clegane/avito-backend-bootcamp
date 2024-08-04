package house

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

func (s *Service) CreateHouse(ctx context.Context, Address, Developer string, Year int64) (*model.House, error) {
	// TODO
	// insert to repo
	return nil, nil
}
