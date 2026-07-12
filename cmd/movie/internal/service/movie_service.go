package service

import (
	"context"
	"fmt"

	"github.com/abozorov/cinema/cmd/movie/internal/models"
)

type Service struct {
	repo models.MovieRepository
}

func New(repo models.MovieRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, m *models.Movie) (int, error) {
	// check validations
	err := m.IsValid()
	if err != nil {
		return 0, fmt.Errorf("movie_service.Create: %w", err)
	}

	// saving in db
	id, err := s.repo.Create(ctx, m)
	if err != nil {
		return 0, fmt.Errorf("movie_service.Create: %w", err)
	}

	// return
	return id, nil
}

func (s *Service) GetByID(ctx context.Context, id int) (*models.Movie, error) {
	// validate id
	if id < 1 {
		return &models.Movie{}, fmt.Errorf("movie_service.GetByID: %w", models.ErrInvalidMovieId)
	}

	// get by id
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return &models.Movie{}, fmt.Errorf("movie_service.GetByID: %w", err)
	}
	return user, nil
}

func (s *Service) List(ctx context.Context) ([]models.Movie, error) {
	// get all
	movies, err := s.repo.List(ctx)
	if err != nil {
		return []models.Movie{}, err

	}

	// return
	return movies, nil
}

func (s *Service) Update(ctx context.Context, m *models.Movie) error {
	// validate id
	if m.ID < 1 {
		return fmt.Errorf("movie_service.Update: %w", models.ErrInvalidMovieId)
	}

	// get old user
	old, err := s.repo.GetByID(ctx, m.ID)
	if err != nil {
		return fmt.Errorf("movie_service.Update: %w", err)
	}

	// validate
	err = old.Update(m.Title, m.Description, m.Duration, m.AgeLimit)
	if err != nil {
		return fmt.Errorf("movie_service.Update: %w", err)
	}

	// update
	err = s.repo.Update(ctx, old)
	if err != nil {
		return fmt.Errorf("movie_service.Update: %w", err)
	}

	// answer
	return nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return nil
}
