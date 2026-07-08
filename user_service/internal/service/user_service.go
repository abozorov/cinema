package service

import "github.com/abozorov/cinema/user_service/internal/repo"

type Service struct {
	repo *repo.Repo
}

func New(repo *repo.Repo) *Service {
	return &Service{
		repo: repo,
	}
}
