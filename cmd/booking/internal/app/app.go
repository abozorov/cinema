package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abozorov/cinema/cmd/booking/internal/config"
	"github.com/abozorov/cinema/cmd/booking/internal/handler"
	"github.com/abozorov/cinema/cmd/booking/internal/repo"
	"github.com/abozorov/cinema/cmd/booking/internal/service"
	bookingv1 "github.com/abozorov/cinema/grpc_api/generate/bookingpb/booking/v1"

	moviev1 "github.com/abozorov/cinema/grpc_api/generate/moviepb/movie/v1"
	userv1 "github.com/abozorov/cinema/grpc_api/generate/userpb/user/v1"
	"github.com/abozorov/cinema/pkg/grpcserver"
	"github.com/abozorov/cinema/pkg/logger"
	"github.com/abozorov/cinema/pkg/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type serv struct {
	grpc *grpcserver.Server
}

func initServ(conf *config.Config, l *logger.Logger, h *handler.Handler) *serv {

	grpcServer := grpcserver.New(
		l,
		conf.HTTP.Port,
	)
	bookingv1.RegisterBookingServiceServer(grpcServer.App, h)

	return &serv{
		grpc: grpcServer,
	}
}

func (s *serv) waitForShutdown(l *logger.Logger) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error

	select {
	case sig := <-interrupt:
		l.Info(fmt.Sprintf("app - Run - signal: %s", sig.String()))
	case err = <-s.grpc.Notify():
		l.Error(fmt.Sprintf("app - Run - grpcServer.Notify: %s", err.Error()))
	}
	s.shutdownServers(l)
}

func (s *serv) shutdownServers(l *logger.Logger) {
	if err := s.grpc.Shutdown(); err != nil {
		l.Error(fmt.Sprintf("app - Run - grpcServer.Shutdown: %s", err.Error()))
	}
}

func Run(conf *config.Config) {
	// load logger
	logger, err := logger.NewLogger(true)
	if err != nil {
		log.Fatal("Eror creating logger %w", err)
	}

	// create db connection
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.PG.User,
		conf.PG.Password,
		conf.PG.Host,
		conf.PG.Port,
		conf.PG.Name,
	)
	pg, err := postgres.NewConn(conn)
	if err != nil {
		logger.Fatal(fmt.Sprintf("app - Run - postgres.New: %s", err.Error()))
	}
	defer pg.Close()

	// init clients

	// userClient
	resolver.SetDefaultScheme("dns") // dns:///localhost:5050
	connUserService, err := grpc.Dial(
		fmt.Sprintf("%s:%s", conf.UserService.Host, conf.UserService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatal(fmt.Sprintf("app - Run - grpc.Dial: %s", err.Error()))
	}
	userClient := userv1.NewUserServiceClient(connUserService)

	// movie client
	connMovieService, err := grpc.Dial(
		fmt.Sprintf("%s:%s", conf.MovieService.Host, conf.MovieService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatal(fmt.Sprintf("app - Run - grpc.Dial: %s", err.Error()))
	}
	movieClient := moviev1.NewMovieServiceClient(connMovieService)

	// init layers
	rp := repo.New(pg)
	srvc := service.New(rp, userClient, movieClient)
	hndlr := handler.New(logger, srvc)

	// init servers
	s := initServ(conf, logger, hndlr)

	// start server
	s.grpc.Start()
	s.waitForShutdown(logger)
}
