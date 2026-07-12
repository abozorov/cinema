package handler

import (
	"context"
	"fmt"

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

	// get model
	movie, err := models.NewMovie(
		r.GetTitle(),
		r.GetDescription(),
		int(r.GetDuration()),
		int(r.GetAgeLimit()),
	)
	if err != nil {
		h.logger.Error(fmt.Sprintf("movie_handler.Create: %s", err))
		return nil, responseErr(err)
	}

	// create
	id, err := h.service.Create(ctx, movie)
	if err != nil {
		h.logger.Error(fmt.Sprintf("movie_handler.Create: %s", err))
		return nil, responseErr(err)
	}

	// return
	return &moviev1.CreateMovieResponse{
		Id: int64(id),
	}, nil
}

func (h *Handler) GetByID(ctx context.Context, r *moviev1.GetMovieRequest) (*moviev1.GetMovieResponse, error) {
	// get id
	id := int(r.GetId())

	// get by id
	movie, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.logger.Error(fmt.Sprintf("movie_handler.GetByID: %s", err))
		return nil, responseErr(err)
	}

	// transform model.user -> userResponse
	responseUser := &moviev1.GetMovieResponse{
		Id:          int64(movie.ID),
		Title:       movie.Title,
		Description: movie.Description,
		Duration:    int32(movie.Duration),
		AgeLimit:    int32(movie.AgeLimit),
		CreatedAt:   movie.CreatedAt.Format(models.TimeFormat),
	}

	// return
	return responseUser, nil
}

func (h *Handler) List(ctx context.Context, r *moviev1.ListMovieRequest) (*moviev1.ListMovieResponse, error) {
	// get all
	m, err := h.service.List(ctx)
	if err != nil {
		return nil, responseErr(err)
	}

	// models.movie -> moviev1.ListMovieRequest
	movies := make([]*moviev1.GetMovieResponse, 0, len(m))
	for _, v := range m {
		movies = append(movies, &moviev1.GetMovieResponse{
			Id:          int64(v.ID),
			Title:       v.Title,
			Description: v.Description,
			Duration:    int32(v.Duration),
			AgeLimit:    int32(v.AgeLimit),
			CreatedAt:   v.CreatedAt.Format(models.TimeFormat),
		})
	}

	// return
	return &moviev1.ListMovieResponse{
		Movies: movies,
	}, nil
}

func (h *Handler) Update(ctx context.Context, r *moviev1.UpdateMovieRequest) (*moviev1.UpdateMovieResponse, error) {
	// moviev1.UpdateMovieRequest -> models.movie
	user := &models.Movie{
		ID:          int(r.GetId()),
		Title:       r.GetTitle(),
		Description: r.GetDescription(),
		Duration:    int(r.GetDuration()),
		AgeLimit:    int(r.GetAgeLimit()),
	}

	// update
	err := h.service.Update(ctx, user)
	if err != nil {
		return nil, responseErr(err)
	}

	// answer
	return &moviev1.UpdateMovieResponse{
		Code:    int64(codes.OK),
		Message: "User update successfuly!",
	}, nil
}

func (h *Handler) Delete(ctx context.Context, r *moviev1.DeleteMovieRequest) (*moviev1.DeleteMovieResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
