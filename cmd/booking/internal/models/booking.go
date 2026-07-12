package models

import (
	"context"
	"errors"
	"strings"
)

type Booking struct {
	ID      int
	UserID  int
	MovieID int
	Status  string
}

type BookingService interface {
	Create(ctx context.Context, b *Booking) (*Booking, error)
	Get(ctx context.Context, id int) (*Booking, error)
	GetUserBookings(ctx context.Context, userID int) ([]*Booking, error)
	Cancel(ctx context.Context, id int) error
}

type BookingRepository interface {
	Create(ctx context.Context, b *Booking) (int, error)
	Get(ctx context.Context, id int) (*Booking, error)
	GetUserBookings(ctx context.Context, userID int) ([]*Booking, error)
	Cancel(ctx context.Context, id int) error
}

// NewBooking создает новое бронирование с указанными userID и movieID.
// Статус автоматически устанавливается в "pending".
// Возвращает указатель на Booking и nil в случае успеха.
// Возвращает ошибку, если userID или movieID не прошли валидацию.
func NewBooking(userID, movieID int) (*Booking, error) {
	if err := validateBookingUserID(userID); err != nil {
		return nil, err
	}
	if err := validateBookingMovieID(movieID); err != nil {
		return nil, err
	}

	return &Booking{
		UserID:  userID,
		MovieID: movieID,
		Status:  "BOOKING_STATUS_UNSPECIFIED",
	}, nil
}

// Ошибки для валидации бронирования
var (
	ErrInvalidBookingUserId  = errors.New("invalid booking user id")
	ErrInvalidBookingMovieId = errors.New("invalid booking movie id")
	ErrInvalidBookingStatus  = errors.New("invalid booking status")
	ErrEmptyBookingStatus    = errors.New("empty booking status")
)

// Допустимые статусы бронирования
var (
	StatusUnspecified = "BOOKING_STATUS_UNSPECIFIED"
	StatusPending     = "BOOKING_STATUS_PENDING"
	StatusConfirmed   = "BOOKING_STATUS_CONFIRMED"
	StatusCanceled    = "BOOKING_STATUS_CANCELED"

	validBookingStatuses = map[string]bool{
		"BOOKING_STATUS_UNSPECIFIED": true,
		"BOOKING_STATUS_PENDING":     true,
		"BOOKING_STATUS_CONFIRMED":   true,
		"BOOKING_STATUS_CANCELED":    true,
	}
)

func (b *Booking) Validate() error {
	if err := validateBookingUserID(b.UserID); err != nil {
		return err
	}
	if err := validateBookingMovieID(b.MovieID); err != nil {
		return err
	}
	return nil
}

// validateBookingUserID проверяет, что идентификатор пользователя положительный.
func validateBookingUserID(userID int) error {
	if userID <= 0 {
		return ErrInvalidBookingUserId
	}
	return nil
}

// validateBookingMovieID проверяет, что идентификатор фильма положительный.
func validateBookingMovieID(movieID int) error {
	if movieID <= 0 {
		return ErrInvalidBookingMovieId
	}
	return nil
}

// validateBookingStatus проверяет статус бронирования.
// Статус обязателен (не nil и не пустая строка) и должен входить в список разрешённых.
func validateBookingStatus(status *string) error {
	if status == nil || strings.TrimSpace(*status) == "" {
		return ErrEmptyBookingStatus
	}
	if !validBookingStatuses[*status] {
		return ErrInvalidBookingStatus
	}
	return nil
}
