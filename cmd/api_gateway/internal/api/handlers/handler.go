package handlers

import (
	"github.com/abozorov/cinema/cmd/api_gateway/internal/services"
	booking_service "github.com/abozorov/cinema/cmd/api_gateway/internal/services/booking_service"
	movie_service "github.com/abozorov/cinema/cmd/api_gateway/internal/services/movie_service"
	user_service "github.com/abozorov/cinema/cmd/api_gateway/internal/services/user_service"
	"github.com/abozorov/cinema/pkg/logger"
)

type Handler struct {
	serviceManager services.IServiceManager
	userService    *user_service.UserService
	movieService   *movie_service.MovieService
	bookingService *booking_service.BookingService
	logger         *logger.Logger
}

func NewHandler(serviceManager services.IServiceManager,
	userService *user_service.UserService,
	movieService *movie_service.MovieService,
	bookingService *booking_service.BookingService,
	logger *logger.Logger) *Handler {
	return &Handler{
		serviceManager: serviceManager,
		userService:    userService,
		movieService:   movieService,
		bookingService: bookingService,
		logger:         logger,
	}
}
