package tests

import (
	"Avito/internal/handler_manager"
	"Avito/internal/middleware"
	"Avito/internal/model"
	"Avito/internal/repository"
	"Avito/internal/service"
	"Avito/internal/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SendCoin(t *testing.T) {
	sendCoinRequest := model.SendCoinRequest{
		ToUser: "user2",
		Amount: 100,
	}
	regUserReq1 := model.AuthRequest{
		Login:    "user1",
		Password: "password",
	}

	regUserReq2 := model.AuthRequest{
		Login:    "user2",
		Password: "password",
	}

	t.Run("send coin", func(t *testing.T) {
		database.SetUp(t, "users", "send_coin_events")
		repo := repository.NewRepository(database.DB)
		svc := service.NewService(repo)
		jwtGen := &utils.JWTGen{}
		hm := handler_manager.NewHandlerManager(svc, jwtGen)

		body, err := json.Marshal(regUserReq1)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		hm.Auth(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		res := map[string]string{}
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		require.NoError(t, err)

		token, ok := res["token"]
		assert.True(t, ok)

		body, err = json.Marshal(regUserReq2)
		require.NoError(t, err)

		req = httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
		rec = httptest.NewRecorder()

		hm.Auth(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		body, err = json.Marshal(sendCoinRequest)
		require.NoError(t, err)

		req = httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+token)
		rec = httptest.NewRecorder()

		handler := middleware.AuthMiddleware(http.HandlerFunc(hm.SendCoin))
		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func Test_Buy(t *testing.T) {
	regUserReq := model.AuthRequest{
		Login:    "user1",
		Password: "password",
	}

	t.Run("buy", func(t *testing.T) {
		database.SetUp(t, "users", "purchases")
		repo := repository.NewRepository(database.DB)
		svc := service.NewService(repo)
		jwtGen := &utils.JWTGen{}
		hm := handler_manager.NewHandlerManager(svc, jwtGen)

		body, err := json.Marshal(regUserReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		hm.Auth(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		res := map[string]string{}
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		require.NoError(t, err)

		token, ok := res["token"]
		assert.True(t, ok)

		req = httptest.NewRequest(http.MethodGet, "/api/buy/t-shirt", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec = httptest.NewRecorder()

		handler := middleware.AuthMiddleware(http.HandlerFunc(hm.Buy))
		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
