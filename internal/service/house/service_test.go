package house

import (
	"avito-backend-bootcamp/internal/model"
	"context"
	"log/slog"
	"reflect"
	"testing"
)

func TestService_CreateHouse(t *testing.T) {
	t.Parallel()
	type fields struct {
		log            *slog.Logger
		houseRpository HouseRepository
	}
	type args struct {
		ctx       context.Context
		address   string
		developer string
		year      int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.House
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				log:            tt.fields.log,
				houseRpository: tt.fields.houseRpository,
			}
			got, err := s.CreateHouse(tt.args.ctx, tt.args.address, tt.args.developer, tt.args.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreateHouse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.CreateHouse() = %v, want %v", got, tt.want)
			}
		})
	}
}
