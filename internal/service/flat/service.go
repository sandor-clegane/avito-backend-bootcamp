package flat

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

func (s *Service) CreateFlat(ctx context.Context, HouseID, Price, Rooms int64) (*model.Flat, error) {
	// TODO
	// insert to repo
	return nil, nil
}

func (s *Service) UpdateFlat(ctx context.Context, ID int64, status model.FlatStatus) (*model.Flat, error) {
	// TODO
	// update from repo
	// craete flat approved event
	return nil, nil
}

func (s *Service) GetFlatApprovedEvents(ctx context.Context) error {
	return nil
}
