package house

import (
	"avito-backend-bootcamp/internal/model"
	"context"
)

type HouseRepository interface {
	SaveHouse(ctx context.Context, address, developer string, year int64) (*model.House, error)
}
