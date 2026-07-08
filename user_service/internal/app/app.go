package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abozorov/cinema/user_service/config"
	"github.com/abozorov/cinema/user_service/internal/handler"
	"github.com/abozorov/cinema/user_service/internal/repo"
	"github.com/abozorov/cinema/user_service/internal/service"
	"github.com/abozorov/cinema/user_service/pkg/grpcserver"
	"github.com/abozorov/cinema/user_service/pkg/logger"
	"github.com/abozorov/cinema/user_service/pkg/postgres"
	userv1 "github.com/abozorov/cinema/user_service/userpb/user/v1"
)

type serv struct {
	grpc *grpcserver.Server
}

func initServ(cfg *config.Config, l *logger.Logger, h *handler.Handler) *serv {

	grpcServer := grpcserver.New(
		l,
		cfg.HTTP.Port,
	)
	userv1.RegisterUserServiceServer(grpcServer.App, h)

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

func Run(cfg *config.Config) {
	// load logger
	logger, err := logger.NewLogger(true)
	if err != nil {
		log.Fatal("Eror creating logger %w", err)
	}

	// create db connection
	pg, err := postgres.NewConn(*cfg)
	if err != nil {
		logger.Fatal(fmt.Sprintf("app - Run - postgres.New: %s", err.Error()))
	}
	defer pg.Close()

	// init layers
	rp := repo.New(pg)
	srvc := service.New(rp)
	hndlr := handler.New(logger, srvc)

	// init servers
	s := initServ(cfg, logger, hndlr)

	// start server
	s.waitForShutdown(logger)
}
