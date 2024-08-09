package flat

import (
	repoErr "avito-backend-bootcamp/internal/infra/repository"
	"avito-backend-bootcamp/internal/model"
	repository "avito-backend-bootcamp/internal/service/flat/mocks"
	"avito-backend-bootcamp/pkg/utils/sl"
	"context"
	"errors"
	"log/slog"
	"reflect"
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

type MockFlatRepository struct {
	ctrl *gomock.Controller
	mock *repository.MockFlatRepository
}

func newTestFlat() *model.Flat {
	return &model.Flat{
		ID:      testID,
		Rooms:   testRooms,
		Price:   testPrice,
		HouseID: testHouseID,
		Status:  testStatus,
	}
}

func NewMockFlatRepository(ctrl *gomock.Controller) *MockFlatRepository {
	mock := repository.NewMockFlatRepository(ctrl)
	return &MockFlatRepository{
		ctrl: ctrl,
		mock: mock,
	}
}

func TestService_CreateFlat(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockFlatRepository(ctrl)
		mockRepo.mock.EXPECT().
			SaveFlat(gomock.Any(), testHouseID, testPrice, testRooms).
			Return(&model.Flat{
				HouseID: testHouseID,
				ID:      testID,
				Price:   testPrice,
				Rooms:   testRooms,
				Status:  testStatus,
			}, nil)

		s := &Service{
			flatRepository: mockRepo.mock,
			log:            sl.SetupLogger(),
		}

		flat, err := s.CreateFlat(context.Background(), testHouseID, testPrice, testRooms)

		require.NoError(t, err)
		assert.Equal(t, flat, newTestFlat())
	})

	t.Run("house not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockFlatRepository(ctrl)
		mockRepo.mock.EXPECT().
			SaveFlat(gomock.Any(), testHouseID, testPrice, testRooms).
			Return(nil, repoErr.ErrConstraintViolation)

		s := &Service{
			flatRepository: mockRepo.mock,
			log:            sl.SetupLogger(),
		}

		_, err := s.CreateFlat(context.Background(), testHouseID, testPrice, testRooms)

		require.Error(t, err)
		assert.Equal(t, ErrHouseNotExist, err)
	})

	t.Run("error saving flat", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockFlatRepository(ctrl)
		mockRepo.mock.EXPECT().
			SaveFlat(gomock.Any(), testHouseID, testPrice, testRooms).
			Return(nil, errors.New("failed to save flat"))

		s := &Service{
			flatRepository: mockRepo.mock,
			log:            sl.SetupLogger(),
		}

		_, err := s.CreateFlat(context.Background(), testHouseID, testPrice, testRooms)

		require.Error(t, err)
		assert.NotEqual(t, ErrHouseNotExist, err)
	})
}

func TestService_UpdateFlat(t *testing.T) {
	type fields struct {
		log             *slog.Logger
		flatRepository  FlatRepository
		eventRepository EventRepository
		cache           Cache
		trManager       TrManager
	}
	type args struct {
		ctx    context.Context
		ID     int64
		status model.FlatStatus
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantFlat *model.Flat
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				log:             tt.fields.log,
				flatRepository:  tt.fields.flatRepository,
				eventRepository: tt.fields.eventRepository,
				cache:           tt.fields.cache,
				trManager:       tt.fields.trManager,
			}
			gotFlat, err := s.UpdateFlat(tt.args.ctx, tt.args.ID, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdateFlat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFlat, tt.wantFlat) {
				t.Errorf("Service.UpdateFlat() = %v, want %v", gotFlat, tt.wantFlat)
			}
		})
	}
}

func TestService_flatListForClient(t *testing.T) {
	type fields struct {
		log             *slog.Logger
		flatRepository  FlatRepository
		eventRepository EventRepository
		cache           Cache
		trManager       TrManager
	}
	type args struct {
		ctx     context.Context
		houseID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Flat
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				log:             tt.fields.log,
				flatRepository:  tt.fields.flatRepository,
				eventRepository: tt.fields.eventRepository,
				cache:           tt.fields.cache,
				trManager:       tt.fields.trManager,
			}
			got, err := s.flatListForClient(tt.args.ctx, tt.args.houseID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.flatListForClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.flatListForClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetFlatListByHouseID(t *testing.T) {
	type fields struct {
		log             *slog.Logger
		flatRepository  FlatRepository
		eventRepository EventRepository
		cache           Cache
		trManager       TrManager
	}
	type args struct {
		ctx      context.Context
		houseID  int64
		userRole model.UserType
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantFlatList []*model.Flat
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				log:             tt.fields.log,
				flatRepository:  tt.fields.flatRepository,
				eventRepository: tt.fields.eventRepository,
				cache:           tt.fields.cache,
				trManager:       tt.fields.trManager,
			}
			gotFlatList, err := s.GetFlatListByHouseID(tt.args.ctx, tt.args.houseID, tt.args.userRole)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetFlatListByHouseID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFlatList, tt.wantFlatList) {
				t.Errorf("Service.GetFlatListByHouseID() = %v, want %v", gotFlatList, tt.wantFlatList)
			}
		})
	}
}
