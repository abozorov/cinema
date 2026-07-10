package handler

import (
	"context"

	"github.com/abozorov/cinema/cmd/user/internal/service"
	userv1 "github.com/abozorov/cinema/grpc_api/generate/userpb/user/v1"
	"github.com/abozorov/cinema/pkg/logger"
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

func (c *Handler) Add(ctx context.Context, r *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	
	return &userv1.CreateUserResponse{}, nil
}

func (c *Handler) GetByID(ctx context.Context, r *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {

	return &userv1.GetUserResponse{}, nil
}

func (c *Handler) Update(ctx context.Context, r *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {

	return &userv1.UpdateUserResponse{}, nil
}
