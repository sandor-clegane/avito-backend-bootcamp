package house

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	repoErr "avito-backend-bootcamp/internal/infra/repository"
	"avito-backend-bootcamp/internal/model"
	repository "avito-backend-bootcamp/internal/service/house/mocks"
	"avito-backend-bootcamp/pkg/utils/sl"
)

const (
	testAddress   = "123 Main St"
	testDeveloper = "Acme Developers"
	testYear      = int64(2024)
	testID        = int64(3)
)

type MockHouseRepository struct {
	ctrl *gomock.Controller
	mock *repository.MockHouseRepository
}

func NewMockHouseRepository(ctrl *gomock.Controller) *MockHouseRepository {
	mock := repository.NewMockHouseRepository(ctrl)
	return &MockHouseRepository{
		ctrl: ctrl,
		mock: mock,
	}
}

func TestService_CreateHouse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockHouseRepository(ctrl)
		mockRepo.mock.EXPECT().
			SaveHouse(gomock.Any(), testAddress, testDeveloper, testYear).
			Return(&model.House{ID: testID}, nil)

		s := &Service{
			houseRpository: mockRepo.mock,
			log:            sl.SetupLogger(),
		}

		house, err := s.CreateHouse(context.Background(), testAddress, testDeveloper, testYear)

		require.NoError(t, err)
		assert.Equal(t, testID, house.ID)
	})

	t.Run("house already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockHouseRepository(ctrl)
		mockRepo.mock.EXPECT().
			SaveHouse(gomock.Any(), testAddress, testDeveloper, testYear).
			Return(nil, repoErr.ErrAlreadyExists)

		s := &Service{
			houseRpository: mockRepo.mock,
			log:            sl.SetupLogger(),
		}

		_, err := s.CreateHouse(context.Background(), testAddress, testDeveloper, testYear)

		require.Error(t, err)
		assert.Equal(t, ErrAddressAlreadyUsed, err)
	})

	t.Run("error saving house", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockHouseRepository(ctrl)
		mockRepo.mock.EXPECT().
			SaveHouse(gomock.Any(), testAddress, testDeveloper, testYear).
			Return(nil, errors.New("failed to save house"))

		s := &Service{
			houseRpository: mockRepo.mock,
			log:            sl.SetupLogger(),
		}

		_, err := s.CreateHouse(context.Background(), testAddress, testDeveloper, testYear)

		require.Error(t, err)
		assert.NotEqual(t, ErrAddressAlreadyUsed, err)
	})
}
