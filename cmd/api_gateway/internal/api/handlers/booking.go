package handlers

import (
	"net/http"
	"strconv"

	"github.com/abozorov/cinema/cmd/api_gateway/internal/models"
	mycontext "github.com/abozorov/cinema/cmd/api_gateway/internal/my_context"
	"github.com/abozorov/cinema/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type createBookingRequest struct {
	MovieID int `json:"movie_id"`
}

type bookingResponse struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	MovieID int    `json:"movie_id"`
	Status  string `json:"status"`
}

func newBookingResponse(b models.Booking) *bookingResponse {
	return &bookingResponse{
		ID:      b.ID,
		UserID:  b.UserID,
		MovieID: b.MovieID,
		Status:  b.Status,
	}
}

// CreateBooking - создание брони текущим пользователем
func (h *Handler) CreateBooking(c *gin.Context) {
	// get id текущего пользователя из контекста
	userID, ok := c.Request.Context().Value(mycontext.UserIDKey).(int)
	if !ok {
		h.logger.Error("booking.Create: ", zap.String("error", errs.ErrIncorrectLoginOrPassword.Error()))
		errsToHttp(c.Writer, errs.ErrIncorrectLoginOrPassword)
		return
	}

	// get booking from body
	req := createBookingRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("booking.Create: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// creating & transform request -> models.Booking
	booking := models.Booking{
		UserID:  userID,
		MovieID: req.MovieID,
	}

	created, err := h.bookingService.Create(c.Request.Context(), booking)
	if err != nil {
		h.logger.Error("booking.Create: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	resp := *newBookingResponse(*created)

	c.JSON(http.StatusCreated, resp)
}

// GetBooking - получение брони по id
func (h *Handler) GetBooking(c *gin.Context) {
	// get id из URL параметра
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("booking.GetBooking: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// get by id
	booking, err := h.bookingService.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("booking.GetBooking: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	resp := *newBookingResponse(*booking)

	c.JSON(http.StatusOK, resp)
}

// GetUserBookings - получение всех броней текущего пользователя
func (h *Handler) GetUserBookings(c *gin.Context) {
	// get id из URL параметра
	idStr := c.Param("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("booking.GetBooking: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	bookings, err := h.bookingService.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("booking.GetUserBookings: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// transform []models.Booking -> []response
	resp := make([]bookingResponse, 0, len(bookings))
	for _, b := range bookings {
		resp = append(resp, *newBookingResponse(b))
	}

	c.JSON(http.StatusOK, resp)
}

// CancelBooking - отмена брони по id
func (h *Handler) CancelBooking(c *gin.Context) {
	// get id из URL параметра
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("booking.CancelBooking: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// get id текущего пользователя из контекста (нужен, чтобы проверить владельца брони)
	userID, ok := c.Request.Context().Value(mycontext.UserIDKey).(int)
	if !ok {
		h.logger.Error("booking.CancelBooking: ", zap.String("error", errs.ErrIncorrectLoginOrPassword.Error()))
		errsToHttp(c.Writer, errs.ErrIncorrectLoginOrPassword)
		return
	}

	if err := h.bookingService.Cancel(c.Request.Context(), id, userID); err != nil {
		h.logger.Error("booking.CancelBooking: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	c.String(http.StatusOK, "Booking cancelled")
}
