package utils

import (
	"Avito/internal/middleware"
	"context"
)

func GetUserLogin(ctx context.Context) (string, bool) {
	login, ok := ctx.Value(middleware.UserLoginKey).(string)
	return login, ok
}
