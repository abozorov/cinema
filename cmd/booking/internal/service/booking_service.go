package service

import (
	"context"
	"fmt"

	"github.com/abozorov/cinema/cmd/booking/internal/models"
	moviev1 "github.com/abozorov/cinema/grpc_api/generate/moviepb/movie/v1"
	userv1 "github.com/abozorov/cinema/grpc_api/generate/userpb/user/v1"
	"github.com/abozorov/cinema/pkg/errs"
)

type Service struct {
	repo        models.BookingRepository
	userClient  userv1.UserServiceClient
	MovieClient moviev1.MovieServiceClient
}

func New(repo models.BookingRepository,
	userClient userv1.UserServiceClient,
	MovieClient moviev1.MovieServiceClient,
) *Service {
	return &Service{
		repo:        repo,
		userClient:  userClient,
		MovieClient: MovieClient,
	}
}

func (s *Service) Create(ctx context.Context, b *models.Booking) (*models.Booking, error) {
	//validate
	err := b.Validate()
	if err != nil {
		return nil, fmt.Errorf("service.Create: %w", err)
	}

	// get uset by id
	userResp, err := s.userClient.GetByID(ctx, &userv1.GetUserRequest{
		Id: int64(b.UserID),
	})
	if err != nil {
		return nil, fmt.Errorf("service.Create: %w", err)
	}

	// get movie
	movie, err := s.MovieClient.GetByID(ctx, &moviev1.GetMovieRequest{
		Id: int64(b.MovieID),
	})
	if err != nil {
		return nil, fmt.Errorf("service.Create: %w", err)
	}

	// check age
	if userResp.GetAge() < movie.GetAgeLimit() {
		return nil, fmt.Errorf("service.Create: %w", errs.ErrBadRequest)
	}

	// Формируем объект для сохранения с обязательным статусом PENDING.
	booking := &models.Booking{
		UserID:  b.UserID,
		MovieID: b.MovieID,
		Status:  models.StatusPending,
	}

	// Вызываем репозиторий.
	id, err := s.repo.Create(ctx, booking)
	if err != nil {
		return nil, fmt.Errorf("service.Create: %w", err)
	}

	// Присваиваем сгенерированный ID и возвращаем.
	booking.ID = id
	return booking, nil
}

func (s *Service) Get(ctx context.Context, id int) (*models.Booking, error) {
	if id <= 0 {
		return nil, fmt.Errorf("service.Get: %w", models.ErrInvalidBookingUserId) // Или своя ошибка для ID
	}

	booking, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service.Get: %w", err)
	}
	return booking, nil
}

func (s *Service) GetUserBookings(ctx context.Context, userID int) ([]*models.Booking, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("service.GetUserBookings: %w", models.ErrInvalidBookingUserId)
	}

	bookings, err := s.repo.GetUserBookings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service.GetUserBookings: %w", err)
	}
	return bookings, nil
}

func (s *Service) Cancel(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("service.Cancle: %w", models.ErrInvalidBookingUserId) // используем ту же ошибку для неправильного ID
	}

	err := s.repo.Cancel(ctx, id)
	if err != nil {
		return fmt.Errorf("service.Cancle: %w", err)
	}
	return nil
}
