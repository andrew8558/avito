package handler_manager

import (
	"Avito/internal/errors"
	"Avito/internal/model"
	"encoding/json"
	"net/http"
)

func (hm *HandlerManager) Auth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errors.ErrInvalidHtppMethod.Error(), http.StatusMethodNotAllowed)
		return
	}

	var req model.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, errors.ErrInvalidJson.Error(), http.StatusBadRequest)
		return
	}

	if req.Login == "" {
		http.Error(w, "empty login", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 8 {
		http.Error(w, "password must contain at least 8 characters", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := hm.svc.GetUser(ctx, req.Login, req.Password)

	switch err {
	case errors.ErrAuthFailed:
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	case nil:
		token, err := hm.jwtGen.GenerateJWT(req.Login)
		if err != nil {
			http.Error(w, errors.ErrInternalServerError.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"token": token})
		return
	default:
		http.Error(w, errors.ErrInternalServerError.Error(), http.StatusInternalServerError)
		return
	}
}
