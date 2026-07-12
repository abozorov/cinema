package handler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/abozorov/cinema/cmd/booking/internal/models"
	bookingv1 "github.com/abozorov/cinema/grpc_api/generate/bookingpb/booking/v1"
	"github.com/abozorov/cinema/pkg/logger"
)

type Handler struct {
	bookingv1.UnimplementedBookingServiceServer
	logger  *logger.Logger
	service models.BookingService
}

func New(logger *logger.Logger, service models.BookingService) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func bookingToResponseBooking(booking *models.Booking) *bookingv1.Booking {
	return &bookingv1.Booking{
		Id:      strconv.Itoa(booking.ID),
		UserId:  strconv.Itoa(booking.UserID),
		MovieId: strconv.Itoa(booking.MovieID),
		Status:  bookingv1.BookingStatus(bookingv1.BookingStatus_value[booking.Status]),
	}
}

func (h *Handler) CreateBooking(ctx context.Context, r *bookingv1.CreateBookingRequest) (*bookingv1.CreateBookingResponse, error) {
	// get model
	userId, err := strconv.Atoi(r.GetUserId())
	if err != nil {
		h.logger.Error(fmt.Sprintf("booking_handler.Create: %s", err))
		return nil, responseErr(err)
	}

	movieId, err := strconv.Atoi(r.GetMovieId())
	if err != nil {
		h.logger.Error(fmt.Sprintf("booking_handler.Create: %s", err))
		return nil, responseErr(err)
	}

	booking, err := models.NewBooking(
		userId,
		movieId,
	)
	// create
	booking, err = h.service.Create(ctx, booking)
	if err != nil {
		h.logger.Error(fmt.Sprintf("booking_handler.Create: %s", err))
		return nil, responseErr(err)
	}

	// return
	return &bookingv1.CreateBookingResponse{
		Booking: bookingToResponseBooking(booking),
	}, nil

}

func (h *Handler) GetBooking(ctx context.Context, r *bookingv1.GetBookingRequest) (*bookingv1.GetBookingResponse, error) {
	id, err := strconv.Atoi(r.GetId())
	if err != nil {
		h.logger.Error(fmt.Sprintf("booking_handler.GetBooking: %s", err))
		return nil, responseErr(err)
	}

	booking, err := h.service.Get(ctx, id)
	if err != nil {
		h.logger.Error(fmt.Sprintf("booking_handler.GetBooking: %s", err))
		return nil, responseErr(err)
	}

	return &bookingv1.GetBookingResponse{
		Booking: bookingToResponseBooking(booking),
	}, nil

}

func (h *Handler) GetUserBookings(ctx context.Context, r *bookingv1.GetUserBookingsRequest) (*bookingv1.GetUserBookingsResponse, error) {
	// get id
	userId, err := strconv.Atoi(r.GetUserId())
	if err != nil {
		h.logger.Error(fmt.Sprintf("booking_handler.GetUserBookings: %s", err))
		return nil, responseErr(err)
	}

	// get bookings
	bookings, err := h.service.GetUserBookings(ctx, userId)
	if err != nil {
		h.logger.Error(fmt.Sprintf("booking_handler.GetUserBookings: %s", err))
		return nil, responseErr(err)
	}

	// transform
	bookingsResp := make([]*bookingv1.Booking, 0, len(bookings))
	for _, v := range bookings {
		bookingsResp = append(bookingsResp, bookingToResponseBooking(v))
	}

	// return
	return &bookingv1.GetUserBookingsResponse{
		Bookings: bookingsResp,
	}, nil
}

func (h *Handler) CancelBooking(ctx context.Context, r *bookingv1.CancelBookingRequest) (*bookingv1.CancelBookingResponse, error) {
	id, err := strconv.Atoi(r.GetId())
	if err != nil {
		h.logger.Error(fmt.Sprintf("booking_handler.CancelBooking: %s", err))
		return nil, responseErr(err)
	}

	err = h.service.Cancel(ctx, id)
	if err != nil {
		h.logger.Error(fmt.Sprintf("booking_handler.CancelBooking: %s", err))
		return nil, responseErr(err)
	}

	return nil, nil
}
