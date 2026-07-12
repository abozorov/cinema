package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abozorov/cinema/cmd/movie/internal/config"
	"github.com/abozorov/cinema/cmd/movie/internal/handler"
	"github.com/abozorov/cinema/cmd/movie/internal/repo"
	"github.com/abozorov/cinema/cmd/movie/internal/service"
	moviev1 "github.com/abozorov/cinema/grpc_api/generate/moviepb/movie/v1"
	"github.com/abozorov/cinema/pkg/grpcserver"
	"github.com/abozorov/cinema/pkg/logger"
	"github.com/abozorov/cinema/pkg/postgres"
)

type serv struct {
	grpc *grpcserver.Server
}

func initServ(cfg *config.Config, l *logger.Logger, h *handler.Handler) *serv {

	grpcServer := grpcserver.New(
		l,
		cfg.HTTP.Port,
	)
	moviev1.RegisterMovieServiceServer(grpcServer.App, h)

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
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PG.User,
		cfg.PG.Password,
		cfg.PG.Host,
		cfg.PG.Port,
		cfg.PG.Name,
	)
	pg, err := postgres.NewConn(conn)
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
	s.grpc.Start()
	s.waitForShutdown(logger)
}
