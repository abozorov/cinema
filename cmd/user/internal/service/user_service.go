package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/abozorov/cinema/cmd/user/internal/models"
)

type Service struct {
	repo models.UserRepository
}

func New(repo models.UserRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Add(ctx context.Context, user *models.User) (int, error) {
	// check validations
	err := user.IsValid()
	if err != nil {
		return 0, fmt.Errorf("user_service.Add: %w", err)
	}

	// hash password
	// user.PasswordHash, err = password.Hash(user.PasswordHash)
	// if err != nil {
	// 	return 0, fmt.Errorf("user_service.Add: %w", err)
	// }

	// saving in db
	id, err := s.repo.Add(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("user_service.Add: %w", err)
	}

	// return
	return id, nil
}

func (s *Service) GetByID(ctx context.Context, id int) (*models.User, error) {
	// validate id
	if id < 1 {
		return &models.User{}, fmt.Errorf("user_service.GetUserByID: %w", models.ErrInvalidUSerId)
	}

	// get by id
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return &models.User{}, fmt.Errorf("user_service.GetUserByID: %w", err)
	}
	return user, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	// validate id
	email = strings.TrimSpace(email)
	if email == "" {
		return &models.User{}, fmt.Errorf("user_service.GetByEmail: %w", models.ErrInvalidEmail) 
	}

	// get by email
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return &models.User{}, fmt.Errorf("user_service.GetByEmail: %w", err)
	}
	return user, nil
}

func (s *Service) Update(ctx context.Context, user *models.User) error {
	// validate id
	if user.ID < 1 {
		return fmt.Errorf("user_service.Update: %w", models.ErrInvalidUSerId)
	}

	// get old user
	old, err := s.repo.GetByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("user_service.Update: %w", err)
	}

	// validate
	err = old.Update(user.Name, old.Email, user.Phone, old.Age)
	if err != nil {
		return fmt.Errorf("user_service.Update: %w", err)
	}

	// update
	err = s.repo.Update(ctx, old)
	if err != nil {
		return fmt.Errorf("user_service.Update: %w", err)
	}

	// answer
	return nil
}
