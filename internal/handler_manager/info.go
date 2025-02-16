package handler_manager

import (
	"Avito/internal/errors"
	"Avito/internal/middleware"
	"encoding/json"
	"net/http"
)

func (hm *HandlerManager) Info(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, errors.ErrInvalidHtppMethod.Error(), http.StatusMethodNotAllowed)
		return
	}

	login, ok := r.Context().Value(middleware.UserLoginKey).(string)
	if !ok {
		http.Error(w, "user not authenticated", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	res, err := hm.svc.GetUserInfo(ctx, login)

	switch err {
	case errors.ErrUserDoesNotExist:
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	case nil:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return
	default:
		http.Error(w, errors.ErrInternalServerError.Error(), http.StatusInternalServerError)
		return
	}
}
