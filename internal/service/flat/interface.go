package flat

import (
	"avito-backend-bootcamp/internal/model"
	"context"
)

type FlatRepository interface {
	GetFlat(ctx context.Context, ID int64) (*model.Flat, error)
	SaveFlat(ctx context.Context, houseID, price, fooms int64) (*model.Flat, error)
	UpdateFlat(ctx context.Context, flat *model.Flat) (*model.Flat, error)
	FlatListByHouseID(ctx context.Context, houseID int64) ([]*model.Flat, error)
}

type EventRepository interface {
	PublishEvent(ctx context.Context, eventType model.EventType, payload string) error
}

type Cache interface {
	Set(key int64, value string)
	Get(key int64) (string, bool)
	Remove(key int64)
}

type TrManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) (err error)
}
