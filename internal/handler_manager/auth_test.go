package handler_manager

import (
	customErrors "Avito/internal/errors"
	"Avito/internal/model"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Auth(t *testing.T) {
	t.Parallel()

	var (
		request = model.AuthRequest{
			Login:    "test_user",
			Password: "password",
		}
		requestWithEmptyLogin = model.AuthRequest{
			Login: "",
		}
		requestWithShortPassword = model.AuthRequest{
			Login:    "test_user",
			Password: "pass",
		}
		expectedBody = "{\"token\":\"token\"}\n"
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.mockJWTGen.EXPECT().GenerateJWT(gomock.Any()).Return("token", nil)

		body, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		s.hm.Auth(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedBody, rec.Body.String())
	})

	t.Run("empty login", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		body, _ := json.Marshal(requestWithEmptyLogin)
		req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		s.hm.Auth(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("short password", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		body, _ := json.Marshal(requestWithShortPassword)
		req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		s.hm.Auth(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid http method", func(t *testing.T) {
		s := setUp(t)
		defer s.tearDown()

		body, _ := json.Marshal(requestWithShortPassword)
		req := httptest.NewRequest(http.MethodGet, "/api/auth", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		s.hm.Auth(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to get user"))

		body, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		s.hm.Auth(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("wrong credentials", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(customErrors.ErrAuthFailed)

		body, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		s.hm.Auth(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid joson", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		body, _ := json.Marshal("bad request")
		req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		s.hm.Auth(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
