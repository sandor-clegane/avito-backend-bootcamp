package handlers

import (
	"avito-backend-bootcamp/internal/model"
	"context"

	"github.com/google/uuid"
)

type AuthService interface {
	DummyLogin(ctx context.Context, role model.UserType) (string, error)
	Login(ctx context.Context, ID uuid.UUID, password string) (string, error)
	Register(ctx context.Context, email, password string, role model.UserType) (uuid.UUID, error)
}

type FlatService interface {
	CreateFlat(ctx context.Context, houseID, price, fooms int64) (*model.Flat, error)
	UpdateFlat(ctx context.Context, ID int64, status model.FlatStatus) (*model.Flat, error)
	GetFlatListByHouseID(ctx context.Context, houseID int64, userRole model.UserType) ([]*model.Flat, error)
}

type HouseService interface {
	CreateHouse(ctx context.Context, address, developer string, year int64) (*model.House, error)
}

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, houseID int64, email string) error
}
