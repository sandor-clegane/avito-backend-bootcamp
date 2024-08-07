package handlers

import (
	"log/slog"
	"net/http"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type args struct {
		log         *slog.Logger
		flatService FlatService
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
			if got := New(tt.args.log, tt.args.flatService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
