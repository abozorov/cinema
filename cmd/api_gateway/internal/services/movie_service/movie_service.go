package movie

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/abozorov/cinema/cmd/api_gateway/internal/models"
	"github.com/abozorov/cinema/cmd/api_gateway/internal/services"
	moviev1 "github.com/abozorov/cinema/grpc_api/generate/moviepb/movie/v1"
)

type MovieService struct {
	serviceManager services.IServiceManager
}

func NewMovieService(
	serviceManager services.IServiceManager) *MovieService {

	return &MovieService{
		serviceManager: serviceManager,
	}
}

// Create - создание фильма
func (m *MovieService) Create(ctx context.Context, movie models.Movie) (int, error) {
	// validation
	movie.Title = strings.TrimSpace(movie.Title)
	movie.Description = strings.TrimSpace(movie.Description)

	// creating
	id, err := m.serviceManager.MovieService().Create(ctx, &moviev1.CreateMovieRequest{
		Title:       movie.Title,
		Description: movie.Description,
		Duration:    int32(movie.Duration),
		AgeLimit:    int32(movie.AgeLimit),
	})
	if err != nil {
		return 0, fmt.Errorf("movie_service.Create: %w", err)
	}

	return int(id.GetId()), nil
}

// GetByID - получение фильма по id
func (m *MovieService) GetByID(ctx context.Context, id int) (*models.Movie, error) {
	movie, err := m.serviceManager.MovieService().GetByID(ctx, &moviev1.GetMovieRequest{
		Id: int64(id),
	})
	if err != nil {
		return nil, fmt.Errorf("movie_service.GetByID: %w", err)
	}

	tm, err := time.Parse(models.TimeFormat, movie.GetCreatedAt())
	if err != nil {
		return nil, fmt.Errorf("movie_service.GetByID: %w", err)
	}

	return &models.Movie{
		ID:          int(movie.GetId()),
		Title:       movie.GetTitle(),
		Description: movie.GetDescription(),
		Duration:    int(movie.GetDuration()),
		AgeLimit:    int(movie.GetAgeLimit()),
		CreatedAt:   tm,
	}, nil
}

// Update - обновление фильма
func (m *MovieService) Update(ctx context.Context, movie models.Movie) error {
	// validation
	movie.Title = strings.TrimSpace(movie.Title)
	movie.Description = strings.TrimSpace(movie.Description)

	// updating
	_, err := m.serviceManager.MovieService().Update(ctx, &moviev1.UpdateMovieRequest{
		Id:          int64(movie.ID),
		Title:       movie.Title,
		Description: movie.Description,
		Duration:    int32(movie.Duration),
		AgeLimit:    int32(movie.AgeLimit),
	})
	if err != nil {
		return fmt.Errorf("movie_service.Update: %w", err)
	}

	return nil
}

// List - получение списка фильмов
func (m *MovieService) List(ctx context.Context) ([]models.Movie, error) {
	movies, err := m.serviceManager.MovieService().List(ctx, &moviev1.ListMovieRequest{})
	if err != nil {
		return nil, fmt.Errorf("movie_service.List: %w", err)
	}
	moviesResp := make([]models.Movie, 0, len(movies.Movies))

	for _, movie := range movies.Movies {
		tm, err := time.Parse(models.TimeFormat, movie.GetCreatedAt())
		if err != nil {
			return nil, fmt.Errorf("movie_service.List: %w", err)
		}
		moviesResp = append(moviesResp, models.Movie{
			ID:          int(movie.GetId()),
			Title:       movie.GetTitle(),
			Description: movie.GetDescription(),
			Duration:    int(movie.GetDuration()),
			AgeLimit:    int(movie.GetAgeLimit()),
			CreatedAt:   tm,
		})
	}
	return moviesResp, nil
}
