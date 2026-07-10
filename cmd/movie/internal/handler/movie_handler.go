package handler

import (
	"context"

	"github.com/abozorov/cinema/cmd/movie/internal/models"
	moviev1 "github.com/abozorov/cinema/grpc_api/generate/moviepb/movie/v1"
	"github.com/abozorov/cinema/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	moviev1.UnimplementedMovieServiceServer
	logger  *logger.Logger
	service models.MovieService
}

func New(logger *logger.Logger, service models.MovieService) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) Create(ctx context.Context, r *moviev1.CreateMovieRequest) (*moviev1.CreateMovieResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}

func (h *Handler) GetByID(ctx context.Context, r *moviev1.GetMovieRequest) (*moviev1.GetMovieResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByID not implemented")
}

func (h *Handler) List(ctx context.Context, r *moviev1.ListMovieRequest) (*moviev1.ListMovieResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}

func (h *Handler) Update(ctx context.Context, r *moviev1.UpdateMovieRequest) (*moviev1.UpdateMovieResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func (h *Handler) Delete(ctx context.Context, r *moviev1.DeleteMovieRequest) (*moviev1.DeleteMovieResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
