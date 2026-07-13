package booking

import (
	"context"
	"fmt"
	"strconv"

	"github.com/abozorov/cinema/cmd/api_gateway/internal/models"
	"github.com/abozorov/cinema/cmd/api_gateway/internal/services"
	bookingpb "github.com/abozorov/cinema/grpc_api/generate/bookingpb/booking/v1"
)

type BookingService struct {
	serviceManager services.IServiceManager
}

func NewBookingService(
	serviceManager services.IServiceManager) *BookingService {

	return &BookingService{
		serviceManager: serviceManager,
	}
}

// Create - создание брони
func (b *BookingService) Create(ctx context.Context, booking models.Booking) (*models.Booking, error) {
	resp, err := b.serviceManager.BookingService().CreateBooking(ctx, &bookingpb.CreateBookingRequest{
		UserId:  strconv.Itoa(booking.UserID),
		MovieId: strconv.Itoa(booking.MovieID),
	})

	if err != nil {
		return nil, fmt.Errorf("booking_service.Create: %w", err)
	}

	bId, err := strconv.Atoi(resp.Booking.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid booking id %q: %w", resp.Booking.GetId(), err)
	}

	userId, err := strconv.Atoi(resp.Booking.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("invalid booking id %q: %w", resp.Booking.GetUserId(), err)
	}

	movieId, err := strconv.Atoi(resp.Booking.GetMovieId())
	if err != nil {
		return nil, fmt.Errorf("invalid booking id %q: %w", resp.Booking.GetMovieId(), err)
	}

	return &models.Booking{
		ID:      bId,
		UserID:  userId,
		MovieID: movieId,
		Status:  models.BookingStatus[int(resp.Booking.GetStatus())],
	}, nil
}

// GetByID - получение брони по id
func (b *BookingService) GetByID(ctx context.Context, id int) (*models.Booking, error) {
	resp, err := b.serviceManager.BookingService().GetBooking(ctx, &bookingpb.GetBookingRequest{
		Id: strconv.Itoa(id),
	})

	if err != nil {
		return nil, fmt.Errorf("booking_service.GetByID: %w", err)
	}

	bId, err := strconv.Atoi(resp.Booking.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid booking id %q: %w", resp.Booking.GetId(), err)
	}

	userId, err := strconv.Atoi(resp.Booking.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("invalid booking id %q: %w", resp.Booking.GetUserId(), err)
	}

	movieId, err := strconv.Atoi(resp.Booking.GetMovieId())
	if err != nil {
		return nil, fmt.Errorf("invalid booking id %q: %w", resp.Booking.GetMovieId(), err)
	}

	return &models.Booking{
		ID:      bId,
		UserID:  userId,
		MovieID: movieId,
		Status:  models.BookingStatus[int(resp.Booking.GetStatus())],
	}, nil
}

// GetByUserID - получение всех броней пользователя
func (b *BookingService) GetByUserID(ctx context.Context, userID int) ([]models.Booking, error) {
	resp, err := b.serviceManager.BookingService().GetUserBookings(ctx, &bookingpb.GetUserBookingsRequest{
		UserId: strconv.Itoa(userID),
	})
	if err != nil {
		return nil, fmt.Errorf("booking_service.GetByUserID: %w", err)
	}

	bookings := make([]models.Booking, 0, len(resp.GetBookings()))
	for _, item := range resp.GetBookings() {

		bId, err := strconv.Atoi(item.GetId())
		if err != nil {
			return nil, fmt.Errorf("invalid booking id %q: %w", item.GetId(), err)
		}

		userId, err := strconv.Atoi(item.GetUserId())
		if err != nil {
			return nil, fmt.Errorf("invalid booking id %q: %w", item.GetUserId(), err)
		}

		movieId, err := strconv.Atoi(item.GetMovieId())
		if err != nil {
			return nil, fmt.Errorf("invalid booking id %q: %w", item.GetMovieId(), err)
		}

		booking := &models.Booking{
			ID:      bId,
			UserID:  userId,
			MovieID: movieId,
			Status:  models.BookingStatus[int(item.GetStatus())],
		}

		bookings = append(bookings, *booking)
	}

	return bookings, nil
}

// Cancel - отмена брони с проверкой владельца
func (b *BookingService) Cancel(ctx context.Context, id int, userID int) error {
	_, err := b.serviceManager.BookingService().CancelBooking(ctx, &bookingpb.CancelBookingRequest{
		Id: strconv.Itoa(id),
	})
	if err != nil {
		return fmt.Errorf("booking_service.Cancel: %w", err)
	}

	return nil
}
