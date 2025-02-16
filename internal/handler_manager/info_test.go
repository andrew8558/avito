package handler_manager

import (
	customErrors "Avito/internal/errors"
	"Avito/internal/middleware"
	"Avito/internal/model"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Info(t *testing.T) {
	t.Parallel()

	var (
		login = "test_user"
	)

	t.Run("smoke", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(&model.InfoResponse{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Info(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get info"))

		req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Info(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("user does not exist", func(t *testing.T) {
		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(nil, customErrors.ErrUserDoesNotExist)

		req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Info(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("not auth", func(t *testing.T) {
		s := setUp(t)
		defer s.tearDown()

		req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
		rec := httptest.NewRecorder()

		s.hm.Info(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid http method", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		req := httptest.NewRequest(http.MethodPost, "/api/info", nil)
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Info(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})
}
