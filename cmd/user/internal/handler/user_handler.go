package handler

import (
	"context"
	"fmt"

	"github.com/abozorov/cinema/cmd/user/internal/models"
	userv1 "github.com/abozorov/cinema/grpc_api/generate/userpb/user/v1"
	"github.com/abozorov/cinema/pkg/logger"
)

type Handler struct {
	userv1.UnimplementedUserServiceServer
	logger  *logger.Logger
	service models.UserService
}

func New(logger *logger.Logger, service models.UserService) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) Add(ctx context.Context, r *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {

	// userv1.CreateUserRequest -> models.user
	user, err := models.NewUser(
		r.Name,
		r.Email,
		r.Phone,
		r.Password,
		int(r.Age),
	)
	if err != nil {
		h.logger.Error(fmt.Sprintf("user_handler.Add: %s", err))
		return &userv1.CreateUserResponse{}, responseErr(err)
	}

	// add user(call service)
	id, err := h.service.Add(ctx, user)
	if err != nil {
		h.logger.Error(fmt.Sprintf("user_handler.Add: %s", err))
		return &userv1.CreateUserResponse{}, responseErr(err)
	}

	// return
	return &userv1.CreateUserResponse{
		Id: int64(id),
	}, nil
}

func (h *Handler) GetByID(ctx context.Context, r *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	// get id
	id := int(r.Id)

	// get by id
	user, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.logger.Error(fmt.Sprintf("user_handler.GetByID: %s", err))
		return &userv1.GetUserResponse{}, responseErr(err)
	}

	// transform model.user -> userResponse
	responseUser := &userv1.GetUserResponse{
		Id:    int64(user.ID),
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Age:   int32(user.Age),
	}

	// return
	return responseUser, nil
}

func (h *Handler) Update(ctx context.Context, r *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	// userv1.UpdateUserRequest -> models.user
	user := &models.User{
		ID: int(r.Id),
		Name: r.Name,
		Phone: r.Phone,
	}
	
	// update
	err := h.service.Update(ctx, user)
	if err != nil {
		return &userv1.UpdateUserResponse{}, responseErr(err) 
	}

	// answer
	return &userv1.UpdateUserResponse{
		Code: 200,
		Message: "User update successfuly!",
	}, nil
}
