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
	mock "avito-backend-bootcamp/internal/http/handlers/update-flat/mocks"
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
	r.Post("/flat/update", h)

	return r
}

func TestHandleUpdateFlat(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock flat service
		flatService := mock.NewMockFlatService(ctrl)
		flatService.
			EXPECT().
			UpdateFlat(gomock.Any(), int64(123), model.StatusApproved).
			Return(&model.Flat{ID: 123, HouseID: 456, Price: 1000000, Rooms: 2, Status: model.StatusApproved}, nil)

		// Create HTTP request
		reqBody := []byte(`{"id": 123, "status": "approved"}`)
		req := httptest.NewRequest(http.MethodPost, "/flat/update", bytes.NewReader(reqBody))

		// Create HTTP response writer
		w := httptest.NewRecorder()

		// Create router
		r := setupRouter(flatService)

		// Execute handler
		r.ServeHTTP(w, req)

		// Assert response status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Assert response body
		var response updateFlatResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, updateFlatResponse{
			ID:      123,
			HouseID: 456,
			Price:   1000000,
			Rooms:   2,
			Status:  "approved",
		}, response)
	})

	t.Run("invalid json", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock flat service
		flatService := mock.NewMockFlatService(ctrl)

		// Create HTTP request
		req := httptest.NewRequest(http.MethodPost, "/flat/update", bytes.NewReader([]byte(`{"id": 123}`)))

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

	t.Run("invalid id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock flat service
		flatService := mock.NewMockFlatService(ctrl)

		// Create HTTP request
		reqBody := []byte(`{"id": 0, "status": "approved"}`)
		req := httptest.NewRequest(http.MethodPost, "/flat/update", bytes.NewReader(reqBody))

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

	t.Run("invalid status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock flat service
		flatService := mock.NewMockFlatService(ctrl)

		// Create HTTP request
		reqBody := []byte(`{"id": 123, "status": "invalid"}`)
		req := httptest.NewRequest(http.MethodPost, "/flat/update", bytes.NewReader(reqBody))

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
		assert.Equal(t, "unknown enum value invalid", response.Error)
	})

	t.Run("flat not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock flat service
		flatService := mock.NewMockFlatService(ctrl)
		flatService.
			EXPECT().
			UpdateFlat(gomock.Any(), int64(123), model.StatusApproved).
			Return(nil, flatPkg.ErrFlatNotExist)

		// Create HTTP request
		reqBody := []byte(`{"id": 123, "status": "approved"}`)
		req := httptest.NewRequest(http.MethodPost, "/flat/update", bytes.NewReader(reqBody))

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
		assert.Equal(t, flatPkg.ErrFlatNotExist.Error(), response.Error)
	})

	t.Run("impossible transition", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock flat service
		flatService := mock.NewMockFlatService(ctrl)
		flatService.
			EXPECT().
			UpdateFlat(gomock.Any(), int64(123), model.StatusApproved).
			Return(nil, model.ErrImpossibleTransition)

		// Create HTTP request
		reqBody := []byte(`{"id": 123, "status": "approved"}`)
		req := httptest.NewRequest(http.MethodPost, "/flat/update", bytes.NewReader(reqBody))

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
		assert.Equal(t, model.ErrImpossibleTransition.Error(), response.Error)
	})

	t.Run("failed to update flat", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock flat service
		flatService := mock.NewMockFlatService(ctrl)
		flatService.
			EXPECT().
			UpdateFlat(gomock.Any(), int64(123), model.StatusApproved).
			Return(nil, errors.New("internal error"))

		// Create HTTP request
		reqBody := []byte(`{"id": 123, "status": "approved"}`)
		req := httptest.NewRequest(http.MethodPost, "/flat/update", bytes.NewReader(reqBody))
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
		assert.Equal(t, handlers.ErrorResponse{
			Message:   "internal error",
			RequestID: "test",
			Code:      http.StatusInternalServerError,
		}, response)
	})
}
