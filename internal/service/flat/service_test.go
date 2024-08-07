package flat

import (
	"avito-backend-bootcamp/internal/model"
	"context"
	"log/slog"
	"reflect"
	"testing"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

func TestService_CreateFlat(t *testing.T) {
	t.Parallel()
	type fields struct {
		log             *slog.Logger
		flatRepository  FlatRepository
		eventRepository EventRepository
		cache           Cache
		trManager       *manager.Manager
	}
	type args struct {
		ctx     context.Context
		houseID int64
		price   int64
		rooms   int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Flat
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
			got, err := s.CreateFlat(tt.args.ctx, tt.args.houseID, tt.args.price, tt.args.rooms)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreateFlat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.CreateFlat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_UpdateFlat(t *testing.T) {
	t.Parallel()
	type fields struct {
		log             *slog.Logger
		flatRepository  FlatRepository
		eventRepository EventRepository
		cache           Cache
		trManager       *manager.Manager
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

func TestService_GetFlatListByHouseID(t *testing.T) {
	t.Parallel()
	type fields struct {
		log             *slog.Logger
		flatRepository  FlatRepository
		eventRepository EventRepository
		cache           Cache
		trManager       *manager.Manager
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
