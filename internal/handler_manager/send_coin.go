package handler_manager

import (
	"Avito/internal/errors"
	"Avito/internal/middleware"
	"Avito/internal/model"
	"encoding/json"
	"net/http"
)

func (hm *HandlerManager) SendCoin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errors.ErrInvalidHtppMethod.Error(), http.StatusMethodNotAllowed)
		return
	}

	login, ok := r.Context().Value(middleware.UserLoginKey).(string)
	if !ok {
		http.Error(w, "user not authenticated", http.StatusUnauthorized)
		return
	}

	var req model.SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, errors.ErrInvalidJson.Error(), http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "amount must be greater than 0", http.StatusBadRequest)
		return
	}

	if req.ToUser == "" {
		http.Error(w, "empty receiver", http.StatusBadRequest)
		return
	}

	if req.ToUser == login {
		http.Error(w, "cannot send coins to yourself", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := hm.svc.SendCoin(ctx, login, req.ToUser, req.Amount)

	switch err {
	case errors.ErrReceiverNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case errors.ErrNotEnoughMoney:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case errors.ErrUserDoesNotExist:
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	case nil:
		w.WriteHeader(http.StatusOK)
		return
	default:
		http.Error(w, errors.ErrInternalServerError.Error(), http.StatusInternalServerError)
		return
	}
}
