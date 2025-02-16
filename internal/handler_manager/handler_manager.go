package handler_manager

import (
	"Avito/internal/service"
	"Avito/internal/utils"
)

type HandlerManager struct {
	jwtGen utils.JWTGenerator
	svc    service.Service
}

func NewHandlerManager(svc service.Service, jwtGen utils.JWTGenerator) *HandlerManager {
	return &HandlerManager{
		svc:    svc,
		jwtGen: jwtGen,
	}
}
