package handler_manager

import (
	customErrors "Avito/internal/errors"
	"Avito/internal/middleware"
	"Avito/internal/model"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_SendCoin(t *testing.T) {
	t.Parallel()

	var (
		login   = "test_user"
		request = model.SendCoinRequest{
			ToUser: "receiver",
			Amount: 100,
		}
		requestWithNegativeBalance = model.SendCoinRequest{
			ToUser: "receiver",
			Amount: -10,
		}
		requestWithEmptyReceiver = model.SendCoinRequest{
			ToUser: "",
			Amount: 10,
		}
		requestSendToYourself = model.SendCoinRequest{
			ToUser: "test_user",
			Amount: 10,
		}
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().SendCoin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		body, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.SendCoin(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("amount coins <= 0", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		body, _ := json.Marshal(requestWithNegativeBalance)
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.SendCoin(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("empty receiver", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		body, _ := json.Marshal(requestWithEmptyReceiver)
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.SendCoin(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("try send to yorself", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		body, _ := json.Marshal(requestSendToYourself)
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.SendCoin(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("not auth", func(t *testing.T) {
		s := setUp(t)
		defer s.tearDown()

		req := httptest.NewRequest(http.MethodGet, "/api/sendCoin", nil)
		rec := httptest.NewRecorder()

		s.hm.SendCoin(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("receiver not found", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().SendCoin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(customErrors.ErrReceiverNotFound)

		body, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.SendCoin(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("not enough money", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().SendCoin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(customErrors.ErrNotEnoughMoney)

		body, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.SendCoin(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("user does not exist", func(t *testing.T) {
		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().SendCoin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(customErrors.ErrUserDoesNotExist)

		body, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.SendCoin(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().SendCoin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to send coins"))

		body, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.SendCoin(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("invalid http method", func(t *testing.T) {
		s := setUp(t)
		defer s.tearDown()

		body, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodGet, "/api/sendCoin", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.SendCoin(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		s := setUp(t)
		defer s.tearDown()

		body, _ := json.Marshal("bad request")
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.UserLoginKey, login)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.SendCoin(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
