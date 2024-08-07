package handlers

import (
	"log/slog"
	"net/http"
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type args struct {
		log          *slog.Logger
		validate     *validator.Validate
		houseService HouseService
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.log, tt.args.validate, tt.args.houseService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
