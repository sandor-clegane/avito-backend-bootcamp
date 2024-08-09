package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"avito-backend-bootcamp/internal/http/handlers"
	mock "avito-backend-bootcamp/internal/http/handlers/create-house/mocks"
	mwr "avito-backend-bootcamp/internal/http/middleware"
	"avito-backend-bootcamp/internal/model"
	housePkg "avito-backend-bootcamp/internal/service/house"
	dbUtil "avito-backend-bootcamp/pkg/utils/db"
	resp "avito-backend-bootcamp/pkg/utils/response"
	"avito-backend-bootcamp/pkg/utils/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouter(houseService HouseService) *chi.Mux {
	// Create router
	r := chi.NewRouter()

	// Create handler
	h := New(sl.SetupLogger(), validator.New(), houseService)

	// Mount handler on router
	r.Post("/house/create", h)

	return r
}

func TestHandleCreateHouse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock house service
		houseService := mock.NewMockHouseService(ctrl)
		now := time.Now().Truncate(time.Second)
		houseService.
			EXPECT().
			CreateHouse(gomock.Any(), "some address", "some developer", int64(2023)).
			Return(&model.House{
				ID:                 123,
				Address:            "some address",
				YearOfConstruction: 2023,
				Developer:          dbUtil.NewNullString("some developer"),
				CreatedAt:          now,
				UpdatedAt:          now,
			}, nil)

		// Create HTTP request
		reqBody := []byte(`{"address": "some address", "year": 2023, "developer": "some developer"}`)
		req := httptest.NewRequest(http.MethodPost, "/house/create", bytes.NewReader(reqBody))

		// Create HTTP response writer
		w := httptest.NewRecorder()

		// Create router
		r := setupRouter(houseService)

		// Execute handler
		r.ServeHTTP(w, req)

		// Assert response status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Assert response body
		var response createHouseResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, createHouseResponse{
			ID:        123,
			Address:   "some address",
			Year:      2023,
			Developer: dbUtil.FromNullString(dbUtil.NewNullString("some developer")),
			CreatedAt: now,
			UpdatedAt: now,
		}, response)
	})

	t.Run("invalid json", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock house service
		houseService := mock.NewMockHouseService(ctrl)

		// Create HTTP request
		req := httptest.NewRequest(http.MethodPost, "/house/create", bytes.NewReader([]byte(`{"address": "some address"}`)))

		// Create HTTP response writer
		w := httptest.NewRecorder()

		// Create router
		r := setupRouter(houseService)

		// Execute handler
		r.ServeHTTP(w, req)

		// Assert response status code
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Assert response body
		var response resp.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response.Error, "Validation error")
	})

	t.Run("invalid input", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock house service
		houseService := mock.NewMockHouseService(ctrl)

		// Create HTTP request
		reqBody := []byte(`{"address": "some address", "year": 0}`)
		req := httptest.NewRequest(http.MethodPost, "/house/create", bytes.NewReader(reqBody))

		// Create HTTP response writer
		w := httptest.NewRecorder()

		// Create router
		r := setupRouter(houseService)

		// Execute handler
		r.ServeHTTP(w, req)

		// Assert response status code
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Assert response body
		var response resp.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response.Error, "Validation error")
	})

	t.Run("address already used", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock house service
		houseService := mock.NewMockHouseService(ctrl)
		houseService.
			EXPECT().
			CreateHouse(gomock.Any(), "some address", "some developer", int64(2023)).
			Return(nil, housePkg.ErrAddressAlreadyUsed)

		// Create HTTP request
		reqBody := []byte(`{"address": "some address", "year": 2023, "developer": "some developer"}`)
		req := httptest.NewRequest(http.MethodPost, "/house/create", bytes.NewReader(reqBody))

		// Create HTTP response writer
		w := httptest.NewRecorder()

		// Create router
		r := setupRouter(houseService)

		// Execute handler
		r.ServeHTTP(w, req)

		// Assert response status code
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Assert response body
		var response resp.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, housePkg.ErrAddressAlreadyUsed.Error(), response.Error)
	})

	t.Run("failed to create house", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock house service
		houseService := mock.NewMockHouseService(ctrl)
		houseService.
			EXPECT().
			CreateHouse(gomock.Any(), "some address", "some developer", int64(2023)).
			Return(nil, errors.New("internal error"))

		// Create HTTP request
		reqBody := []byte(`{"address": "some address", "year": 2023, "developer": "some developer"}`)
		req := httptest.NewRequest(http.MethodPost, "/house/create", bytes.NewReader(reqBody))
		req = req.WithContext(context.WithValue(req.Context(), mwr.RequestIDKey, "test"))

		// Create HTTP response writer
		w := httptest.NewRecorder()

		// Create router
		r := setupRouter(houseService)

		// Execute handler
		r.ServeHTTP(w, req)

		// Assert response status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Assert response body
		var response handlers.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "internal error", response.Message)
	})
}
