package service

import "github.com/abozorov/cinema/cmd/movie/internal/models"

type Service struct {
	repo models.MovieRepository
}

func New(repo models.MovieRepository) *Service {
	return &Service{
		repo: repo,
	}
}
