package grpcserver

import (
	"context"
	"errors"
	"net"

	"github.com/abozorov/cinema/user_service/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Server struct {
	logger  *logger.Logger
	App     *grpc.Server
	address string
	notify  chan error
	eg      *errgroup.Group
}

func New(logger *logger.Logger, adrr string) *Server {
	return &Server{
		logger:  logger,
		App:     grpc.NewServer(),
		address: adrr,
		notify:  make(chan error, 1),
		eg:      &errgroup.Group{},
	}
}

func (s *Server) Start() {
	go func() error {
		ln, err := net.Listen("tcp", s.address)
		if err != nil {
			s.notify <- err
			close(s.notify)
			return err
		}

		err = s.App.Serve(ln)
		if err != nil {
			s.notify <- err
			close(s.notify)
			return err
		}

		return nil
	}()

	s.logger.Info("grpc server - Server - Started")
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	var shutdownErrors []error

	s.App.GracefulStop()

	err := s.eg.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error(err.Error(), zap.String("error:", "grpc server - Server - Shutdown - s.eg.Wait"))
		shutdownErrors = append(shutdownErrors, err)
	}

	s.logger.Info("grpc server - Server - Shutdown")

	return errors.Join(shutdownErrors...)
}
