package handler_manager

import (
	customErrors "Avito/internal/errors"
	"Avito/internal/middleware"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Buy(t *testing.T) {
	t.Parallel()

	var (
		login = "test_user"
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().MakePurchase(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		req := httptest.NewRequest(http.MethodGet, "/api/get/t-shirt", nil)
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Buy(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("not auth", func(t *testing.T) {
		s := setUp(t)
		defer s.tearDown()

		req := httptest.NewRequest(http.MethodGet, "/api/buy/t-shirt", nil)
		rec := httptest.NewRecorder()

		s.hm.Buy(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("item not found", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().MakePurchase(gomock.Any(), gomock.Any(), gomock.Any()).Return(customErrors.ErrItemNotFound)

		req := httptest.NewRequest(http.MethodGet, "/api/get/shorts", nil)
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Buy(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("user does not exist", func(t *testing.T) {
		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().MakePurchase(gomock.Any(), gomock.Any(), gomock.Any()).Return(customErrors.ErrUserDoesNotExist)

		req := httptest.NewRequest(http.MethodGet, "/api/buy", nil)
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Buy(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("not enough money", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().MakePurchase(gomock.Any(), gomock.Any(), gomock.Any()).Return(customErrors.ErrNotEnoughMoney)

		req := httptest.NewRequest(http.MethodGet, "/api/get/t-shirt", nil)
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Buy(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().MakePurchase(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to buy item"))

		req := httptest.NewRequest(http.MethodGet, "/api/buy/t-shirt", nil)
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Buy(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("invalid http method", func(t *testing.T) {
		s := setUp(t)
		defer s.tearDown()

		req := httptest.NewRequest(http.MethodPost, "/api/buy/t-shirt", nil)
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Buy(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})
}
