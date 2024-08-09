package flat

import (
	"avito-backend-bootcamp/internal/infra/repository"
	repoErr "avito-backend-bootcamp/internal/infra/repository"
	"avito-backend-bootcamp/internal/model"
	mock "avito-backend-bootcamp/internal/service/flat/mocks"
	"avito-backend-bootcamp/pkg/utils/sl"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testHouseID = int64(123)
	testPrice   = int64(100000)
	testRooms   = int64(3)
	testID      = int64(1)
	testStatus  = model.StatusCreated
)

func newTestFlat() *model.Flat {
	return &model.Flat{
		ID:      testID,
		Rooms:   testRooms,
		Price:   testPrice,
		HouseID: testHouseID,
		Status:  testStatus,
	}
}

func TestService_CreateFlat(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock.NewMockFlatRepository(ctrl)
		mockRepo.EXPECT().
			SaveFlat(gomock.Any(), testHouseID, testPrice, testRooms).
			Return(&model.Flat{
				HouseID: testHouseID,
				ID:      testID,
				Price:   testPrice,
				Rooms:   testRooms,
				Status:  testStatus,
			}, nil)

		s := &Service{
			flatRepository: mockRepo,
			log:            sl.SetupLogger(),
		}

		flat, err := s.CreateFlat(context.Background(), testHouseID, testPrice, testRooms)

		require.NoError(t, err)
		assert.Equal(t, flat, newTestFlat())
	})

	t.Run("house not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock.NewMockFlatRepository(ctrl)
		mockRepo.EXPECT().
			SaveFlat(gomock.Any(), testHouseID, testPrice, testRooms).
			Return(nil, repoErr.ErrConstraintViolation)

		s := &Service{
			flatRepository: mockRepo,
			log:            sl.SetupLogger(),
		}

		_, err := s.CreateFlat(context.Background(), testHouseID, testPrice, testRooms)

		require.Error(t, err)
		assert.Equal(t, ErrHouseNotExist, err)
	})

	t.Run("error saving flat", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock.NewMockFlatRepository(ctrl)
		mockRepo.EXPECT().
			SaveFlat(gomock.Any(), testHouseID, testPrice, testRooms).
			Return(nil, errors.New("failed to save flat"))

		s := &Service{
			flatRepository: mockRepo,
			log:            sl.SetupLogger(),
		}

		_, err := s.CreateFlat(context.Background(), testHouseID, testPrice, testRooms)

		require.Error(t, err)
		assert.NotEqual(t, ErrHouseNotExist, err)
	})
}

type mocks struct {
	flatRepository  *mock.MockFlatRepository
	eventRepository *mock.MockEventRepository
	cache           *mock.MockCache
	trManager       *mock.MockTrManager
}

func newMock(ctrl *gomock.Controller) mocks {
	return mocks{
		flatRepository:  mock.NewMockFlatRepository(ctrl),
		eventRepository: mock.NewMockEventRepository(ctrl),
		cache:           mock.NewMockCache(ctrl),
		trManager:       mock.NewMockTrManager(ctrl),
	}
}

func TestUpdateFlat(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := newMock(ctrl)

		flat := model.Flat{
			ID:      1,
			HouseID: 10,
			Status:  model.StatusOnModeration,
		}

		m.flatRepository.
			EXPECT().
			GetFlat(gomock.Any(), int64(1)).
			Return(&flat, nil)
		m.trManager.
			EXPECT().
			Do(gomock.Any(), gomock.Any()).
			Return(nil)

		service := &Service{
			log:             sl.SetupLogger(),
			flatRepository:  m.flatRepository,
			eventRepository: m.eventRepository,
			cache:           m.cache,
			trManager:       m.trManager,
		}

		resultFlat, err := service.UpdateFlat(context.Background(), 1, model.StatusApproved)

		assert.NoError(t, err)
		assert.Equal(t, resultFlat.Status, model.StatusApproved)
	})

	t.Run("flat not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := newMock(ctrl)

		m.flatRepository.
			EXPECT().
			GetFlat(gomock.Any(), int64(1)).
			Return(nil, repository.ErrNotFound)

		service := &Service{
			log:            sl.SetupLogger(),
			flatRepository: m.flatRepository,
		}

		_, err := service.UpdateFlat(context.Background(), 1, model.StatusApproved)

		assert.Equal(t, err, ErrFlatNotExist)
	})

	t.Run("get flat error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := newMock(ctrl)

		m.flatRepository.
			EXPECT().
			GetFlat(gomock.Any(), int64(1)).
			Return(nil, errors.New("failed to get flat"))

		service := &Service{
			log:            sl.SetupLogger(),
			flatRepository: m.flatRepository,
		}

		_, err := service.UpdateFlat(context.Background(), 1, model.StatusApproved)

		assert.Error(t, err)
		assert.NotEqual(t, err, ErrFlatNotExist)
	})

	t.Run("update flat error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := newMock(ctrl)

		m.flatRepository.
			EXPECT().
			GetFlat(gomock.Any(), int64(1)).
			Return(&model.Flat{Status: model.StatusOnModeration}, nil)
		m.trManager.
			EXPECT().
			Do(gomock.Any(), gomock.Any()).
			Return(errors.New("failed to update flat"))

		service := &Service{
			log:            sl.SetupLogger(),
			flatRepository: m.flatRepository,
			trManager:      m.trManager,
		}

		_, err := service.UpdateFlat(context.Background(), 1, model.StatusApproved)

		assert.Error(t, err)
	})

	t.Run("invalid status transition", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := newMock(ctrl)

		m.flatRepository.
			EXPECT().
			GetFlat(gomock.Any(), int64(1)).
			Return(&model.Flat{Status: model.StatusOnModeration}, nil)

		service := &Service{
			log:            sl.SetupLogger(),
			flatRepository: m.flatRepository,
		}

		_, err := service.UpdateFlat(context.Background(), 1, model.StatusOnModeration)

		assert.Equal(t, err, model.ErrImpossibleTransition)
	})
}

func TestFlatListForClient(t *testing.T) {
	t.Run("cache hit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := newMock(ctrl)

		houseID := int64(10)
		flatList := []*model.Flat{
			{ID: 1, HouseID: houseID, Status: model.StatusApproved},
			{ID: 2, HouseID: houseID, Status: model.StatusApproved},
		}
		flatsJSON, _ := json.Marshal(flatList)

		m.cache.
			EXPECT().
			Get(houseID).
			Return(string(flatsJSON), true)

		service := &Service{
			log:            sl.SetupLogger(),
			flatRepository: m.flatRepository,
			cache:          m.cache,
		}

		resultFlatList, err := service.flatListForClient(context.Background(), houseID)

		assert.NoError(t, err)
		assert.Equal(t, resultFlatList, flatList)
	})

	t.Run("cache hit with invalid json", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := newMock(ctrl)

		houseID := int64(10)
		invalidJSON := []byte("invalid json")

		m.cache.
			EXPECT().
			Get(houseID).
			Return(string(invalidJSON), true)
		m.cache.
			EXPECT().
			Remove(houseID)
		m.flatRepository.
			EXPECT().
			FlatListByHouseID(gomock.Any(), gomock.Any()).
			Return(nil, nil)
		m.cache.
			EXPECT().
			Set(houseID, gomock.Any())

		service := &Service{
			log:            sl.SetupLogger(),
			flatRepository: m.flatRepository,
			cache:          m.cache,
		}

		resultFlatList, err := service.flatListForClient(context.Background(), houseID)

		assert.NoError(t, err)
		assert.Equal(t, resultFlatList, []*model.Flat{})
	})

	t.Run("cache miss", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := newMock(ctrl)

		houseID := int64(10)
		flatList := []*model.Flat{
			{ID: 1, HouseID: houseID, Status: model.StatusApproved},
			{ID: 2, HouseID: houseID, Status: model.StatusApproved},
		}

		m.cache.
			EXPECT().
			Get(houseID).
			Return("", false)
		m.flatRepository.
			EXPECT().
			FlatListByHouseID(gomock.Any(), houseID).
			Return(flatList, nil)
		m.cache.
			EXPECT().
			Set(houseID, gomock.Any())
		service := &Service{
			log:            sl.SetupLogger(),
			flatRepository: m.flatRepository,
			cache:          m.cache,
		}

		resultFlatList, err := service.flatListForClient(context.Background(), houseID)

		assert.NoError(t, err)
		assert.Equal(t, resultFlatList, flatList)
	})

	t.Run("database error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := newMock(ctrl)

		houseID := int64(10)
		databaseError := errors.New("database error")

		m.cache.
			EXPECT().
			Get(houseID).
			Return("", false)
		m.flatRepository.
			EXPECT().
			FlatListByHouseID(gomock.Any(), houseID).
			Return(nil, databaseError)

		service := &Service{
			log:            sl.SetupLogger(),
			flatRepository: m.flatRepository,
			cache:          m.cache,
		}

		resultFlatList, err := service.flatListForClient(context.Background(), houseID)

		assert.Error(t, err)
		assert.Equal(t, err, databaseError)
		assert.Nil(t, resultFlatList)
	})

	t.Run("filter active flats", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := newMock(ctrl)

		houseID := int64(10)
		flatList := []*model.Flat{
			{ID: 1, HouseID: houseID, Status: model.StatusApproved},
			{ID: 2, HouseID: houseID, Status: model.StatusOnModeration},
			{ID: 3, HouseID: houseID, Status: model.StatusApproved},
		}

		m.cache.
			EXPECT().
			Get(houseID).
			Return("", false)
		m.flatRepository.
			EXPECT().
			FlatListByHouseID(gomock.Any(), houseID).
			Return(flatList, nil)
		m.cache.
			EXPECT().
			Set(houseID, gomock.Any())

		service := &Service{
			log:            sl.SetupLogger(),
			flatRepository: m.flatRepository,
			cache:          m.cache,
		}

		resultFlatList, err := service.flatListForClient(context.Background(), houseID)

		assert.NoError(t, err)
		assert.Equal(t, len(resultFlatList), 2)
		assert.Equal(t, resultFlatList[0].ID, int64(1))
		assert.Equal(t, resultFlatList[1].ID, int64(3))
	})
}
