package handlers

import (
	pkgCtx "avito-backend-bootcamp/pkg/utils/ctx"
	resp "avito-backend-bootcamp/pkg/utils/response"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"avito-backend-bootcamp/internal/http/handlers"
	mock "avito-backend-bootcamp/internal/http/handlers/get-house/mocks"
	mwr "avito-backend-bootcamp/internal/http/middleware"
	"avito-backend-bootcamp/internal/model"

	"avito-backend-bootcamp/pkg/utils/sl"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouter(flatService FlatService) *chi.Mux {
	// Create router
	r := chi.NewRouter()

	// Create handler
	h := New(sl.SetupLogger(), flatService)

	// Mount handler on router
	r.Get("/house/{id}", h)

	return r
}

func TestHandleGetHouse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock flat service
		flatService := mock.NewMockFlatService(ctrl)
		flatService.
			EXPECT().
			GetFlatListByHouseID(gomock.Any(), int64(123), model.Moderator).
			Return([]*model.Flat{{ID: 1}}, nil)

		// Create HTTP request
		req := httptest.NewRequest(http.MethodGet, "/house/123", nil)
		req = req.WithContext(context.WithValue(req.Context(), pkgCtx.KeyUserType, model.Moderator))

		// Create HTTP response writer
		w := httptest.NewRecorder()

		// Create router
		r := setupRouter(flatService)

		// Execute handler
		r.ServeHTTP(w, req)

		// Assert response status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Assert response body
		var response getHouseResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, []*model.Flat{{ID: 1}}, response.Flats)
	})

	t.Run("invalid house ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock flat service
		flatService := mock.NewMockFlatService(ctrl)

		// Create HTTP request
		req := httptest.NewRequest(http.MethodGet, "/house/abc", nil)
		req = req.WithContext(context.WithValue(req.Context(), pkgCtx.KeyUserType, model.Moderator))

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
		assert.Equal(t, `strconv.ParseInt: parsing "abc": invalid syntax`, response.Error)
	})

	t.Run("failed to get flats", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup mock flat service
		flatService := mock.NewMockFlatService(ctrl)
		flatService.
			EXPECT().
			GetFlatListByHouseID(gomock.Any(), int64(123), model.Moderator).
			Return(nil, errors.New("internal"))

		// Create HTTP request
		req := httptest.NewRequest(http.MethodGet, "/house/123", nil)
		req = req.WithContext(context.WithValue(req.Context(), pkgCtx.KeyUserType, model.Moderator))
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
			Message:   "internal",
			RequestID: "test",
			Code:      http.StatusInternalServerError,
		}, response)
	})
}
