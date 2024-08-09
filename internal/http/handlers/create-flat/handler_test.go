package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"avito-backend-bootcamp/internal/http/handlers"
	mock "avito-backend-bootcamp/internal/http/handlers/create-flat/mocks"
	mwr "avito-backend-bootcamp/internal/http/middleware"
	"avito-backend-bootcamp/internal/model"
	flatPkg "avito-backend-bootcamp/internal/service/flat"

	resp "avito-backend-bootcamp/pkg/utils/response"
	"avito-backend-bootcamp/pkg/utils/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouter(flatService FlatService) *chi.Mux {
	// Create router
	r := chi.NewRouter()

	// Create handler
	h := New(sl.SetupLogger(), validator.New(), flatService)

	// Mount handler on router
	r.Post("/flat/create", h)

	return r
}

func TestHandleCreateFlat(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock flat service
		flatService := mock.NewMockFlatService(ctrl)
		flatService.
			EXPECT().
			CreateFlat(gomock.Any(), int64(123), int64(1000000), int64(2)).
			Return(&model.Flat{ID: 456, HouseID: 123, Price: 1000000, Rooms: 2, Status: model.StatusCreated}, nil)

		// Create HTTP request
		reqBody := []byte(`{"house_id": 123, "price": 1000000, "rooms": 2}`)
		req := httptest.NewRequest(http.MethodPost, "/flat/create", bytes.NewReader(reqBody))

		// Create HTTP response writer
		w := httptest.NewRecorder()

		// Create router
		r := setupRouter(flatService)

		// Execute handler
		r.ServeHTTP(w, req)

		// Assert response status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Assert response body
		var response createFlatResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, createFlatResponse{
			ID:      456,
			HouseID: 123,
			Price:   1000000,
			Rooms:   2,
			Status:  "created",
		}, response)
	})

	t.Run("invalid json", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock flat service
		flatService := mock.NewMockFlatService(ctrl)

		// Create HTTP request
		req := httptest.NewRequest(http.MethodPost, "/flat/create", bytes.NewReader([]byte(`{"house_id": 123}`)))

		// Create HTTP response writer
		w := httptest.NewRecorder()

		// Create router
		r := setupRouter(flatService)

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

		// Setup mock flat service
		flatService := mock.NewMockFlatService(ctrl)

		// Create HTTP request
		reqBody := []byte(`{"house_id": 0, "price": 1000000, "rooms": 2}`)
		req := httptest.NewRequest(http.MethodPost, "/flat/create", bytes.NewReader(reqBody))

		// Create HTTP response writer
		w := httptest.NewRecorder()

		// Create router
		r := setupRouter(flatService)

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

	t.Run("house not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock flat service
		flatService := mock.NewMockFlatService(ctrl)
		flatService.
			EXPECT().
			CreateFlat(gomock.Any(), int64(123), int64(1000000), int64(2)).
			Return(nil, flatPkg.ErrHouseNotExist)

		// Create HTTP request
		reqBody := []byte(`{"house_id": 123, "price": 1000000, "rooms": 2}`)
		req := httptest.NewRequest(http.MethodPost, "/flat/create", bytes.NewReader(reqBody))

		// Create HTTP response writer
		w := httptest.NewRecorder()

		// Create router
		r := setupRouter(flatService)

		// Execute handler
		r.ServeHTTP(w, req)

		// Assert response status code
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Assert response body
		var response resp.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, flatPkg.ErrHouseNotExist.Error(), response.Error)
	})

	t.Run("failed to create flat", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock flat service
		flatService := mock.NewMockFlatService(ctrl)
		flatService.
			EXPECT().
			CreateFlat(gomock.Any(), int64(123), int64(1000000), int64(2)).
			Return(nil, errors.New("internal error"))

		// Create HTTP request
		reqBody := []byte(`{"house_id": 123, "price": 1000000, "rooms": 2}`)
		req := httptest.NewRequest(http.MethodPost, "/flat/create", bytes.NewReader(reqBody))
		req = req.WithContext(context.WithValue(req.Context(), mwr.RequestIDKey, "test"))

		// Create HTTP response writer
		w := httptest.NewRecorder()

		// Create router
		r := setupRouter(flatService)

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
