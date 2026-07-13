package services

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"

	"github.com/abozorov/cinema/cmd/api_gateway/internal/config"
	bookingv1 "github.com/abozorov/cinema/grpc_api/generate/bookingpb/booking/v1"
	moviev1 "github.com/abozorov/cinema/grpc_api/generate/moviepb/movie/v1"
	userv1 "github.com/abozorov/cinema/grpc_api/generate/userpb/user/v1"
)

type IServiceManager interface {
	UserService() userv1.UserServiceClient
	MovieService() moviev1.MovieServiceClient
	BookingService() bookingv1.BookingServiceClient
}

type serviceManager struct {
	userService    userv1.UserServiceClient
	movieService   moviev1.MovieServiceClient
	bookingService bookingv1.BookingServiceClient
}

func (s *serviceManager) UserService() userv1.UserServiceClient {
	return s.userService
}

func (s *serviceManager) MovieService() moviev1.MovieServiceClient {
	return s.movieService
}

func (s *serviceManager) BookingService() bookingv1.BookingServiceClient {
	return s.bookingService
}

func NewServiceManager(config config.Config) (IServiceManager, error) {
	resolver.SetDefaultScheme("dns") // dns:///localhost:50051

	connUserService, err := grpc.Dial(
		fmt.Sprintf("%s:%s", config.UserService.Host, config.UserService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	connMovieService, err := grpc.Dial(
		fmt.Sprintf("%s:%s", config.MovieService.Host, config.MovieService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to movie service: %w", err)
	}

	connBookingService, err := grpc.Dial(
		fmt.Sprintf("%s:%s", config.BookingService.Host, config.BookingService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to booking service: %w", err)
	}

	return &serviceManager{
		userService:    userv1.NewUserServiceClient(connUserService),
		movieService:   moviev1.NewMovieServiceClient(connMovieService),
		bookingService: bookingv1.NewBookingServiceClient(connBookingService),
	}, nil
}
