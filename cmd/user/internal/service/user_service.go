package service

import "github.com/abozorov/cinema/cmd/user/internal/repo"

type Service struct {
	repo *repo.Repo
}

func New(repo *repo.Repo) *Service {
	return &Service{
		repo: repo,
	}
}
