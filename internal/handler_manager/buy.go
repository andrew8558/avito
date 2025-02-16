package handler_manager

import (
	"Avito/internal/errors"
	"Avito/internal/middleware"
	"net/http"
	"strings"
)

func (hm *HandlerManager) Buy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, errors.ErrInvalidHtppMethod.Error(), http.StatusMethodNotAllowed)
		return
	}

	login, ok := r.Context().Value(middleware.UserLoginKey).(string)
	if !ok {
		http.Error(w, "user not authenticated", http.StatusUnauthorized)
		return
	}

	item := strings.TrimPrefix(r.URL.Path, "/api/buy/")

	ctx := r.Context()
	err := hm.svc.MakePurchase(ctx, login, item)

	switch err {
	case errors.ErrNotEnoughMoney:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case errors.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case errors.ErrUserDoesNotExist:
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	case nil:
		w.WriteHeader(http.StatusOK)
		return
	default:
		http.Error(w, errors.ErrInternalServerError.Error(), http.StatusInternalServerError)
	}
}
