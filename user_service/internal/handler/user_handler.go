package handler

import (
	"github.com/abozorov/cinema/user_service/internal/service"
	"github.com/abozorov/cinema/user_service/pkg/logger"
	userv1 "github.com/abozorov/cinema/user_service/userpb/user/v1"
)

type Handler struct {
	userv1.UnimplementedUserServiceServer
	logger  *logger.Logger
	service *service.Service
}

func New(logger *logger.Logger, service *service.Service) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}
