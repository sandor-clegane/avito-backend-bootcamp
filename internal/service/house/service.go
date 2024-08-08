package house

import (
	"avito-backend-bootcamp/internal/infra/repository"
	"avito-backend-bootcamp/internal/model"
	"avito-backend-bootcamp/pkg/utils/sl"
	"context"
	"errors"
	"log/slog"
)

type Service struct {
	log            *slog.Logger
	houseRpository HouseRepository
}

func New(log *slog.Logger, houseRpository HouseRepository) *Service {
	return &Service{
		log:            log,
		houseRpository: houseRpository,
	}
}

var ErrAddressAlreadyUsed = errors.New("house with given address already exist")

func (s *Service) CreateHouse(ctx context.Context, address, developer string, year int64) (*model.House, error) {
	const op = "house.CreateHouse"

	log := s.log.With(
		slog.String("op", op),
		slog.String("address", address),
		slog.String("developer", developer),
		slog.Int64("year", year),
	)

	house, err := s.houseRpository.SaveHouse(ctx, address, developer, year)
	if err != nil {
		if errors.Is(err, repository.ErrConstraintViolation) {
			log.Error("attempt to create invalid house", sl.Err(err))
			return nil, ErrAddressAlreadyUsed
		}
		log.Error("failed to save house", sl.Err(err))
		return nil, err
	}

	return house, nil
}
