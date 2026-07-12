package handler

import (
	"context"
	"fmt"

	"github.com/abozorov/cinema/cmd/user/internal/models"
	userv1 "github.com/abozorov/cinema/grpc_api/generate/userpb/user/v1"
	"github.com/abozorov/cinema/pkg/logger"
	"google.golang.org/grpc/codes"
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
		r.GetName(),
		r.GetEmail(),
		r.GetPhone(),
		r.GetPassword(),
		int(r.GetAge()),
	)
	if err != nil {
		h.logger.Error(fmt.Sprintf("user_handler.Add: %s", err))
		return &userv1.CreateUserResponse{}, responseErr(err)
	}

	// add user(call service)
	id, err := h.service.Add(ctx, user)
	if err != nil {
		h.logger.Error(fmt.Sprintf("user_handler.Add: %s", err))
		return nil, responseErr(err)
	}

	// return
	return &userv1.CreateUserResponse{
		Id: int64(id),
	}, nil
}

func (h *Handler) GetByID(ctx context.Context, r *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	// get id
	id := int(r.GetId())

	// get by id
	user, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.logger.Error(fmt.Sprintf("user_handler.GetByID: %d %s", id,  err))
		return nil, responseErr(err)
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
		ID:    int(r.GetId()),
		Name:  r.GetName(),
		Phone: r.GetPhone(),
	}

	// update
	err := h.service.Update(ctx, user)
	if err != nil {
		return nil, responseErr(err)
	}

	// answer
	return &userv1.UpdateUserResponse{
		Code:    int64(codes.OK),
		Message: "User update successfuly!",
	}, nil
}
