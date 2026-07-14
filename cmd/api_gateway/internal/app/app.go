package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abozorov/cinema/cmd/api_gateway/internal/api"
	"github.com/abozorov/cinema/cmd/api_gateway/internal/api/handlers"
	"github.com/abozorov/cinema/cmd/api_gateway/internal/api/middleware"
	"github.com/abozorov/cinema/cmd/api_gateway/internal/config"
	"github.com/abozorov/cinema/cmd/api_gateway/internal/services"
	booking_service "github.com/abozorov/cinema/cmd/api_gateway/internal/services/booking_service"
	movie_service "github.com/abozorov/cinema/cmd/api_gateway/internal/services/movie_service"
	user_service "github.com/abozorov/cinema/cmd/api_gateway/internal/services/user_service"
	"github.com/abozorov/cinema/pkg/jwt"
	"github.com/abozorov/cinema/pkg/logger"
	mailsender "github.com/abozorov/cinema/pkg/mail_sender"

	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

func Run(conf *config.Config) {
	// load logger
	logger, err := logger.NewLogger(true)
	if err != nil {
		log.Fatal("Eror creating logger %w", err)
	}

	// create SecretJWT
	sJWT := jwt.NewSecretJWT(
		conf.JWT.SecretToken,
		time.Duration(conf.JWT.JWTLiveTime*int(time.Second)),
	)

	// create memCache
	memCache := cache.New(time.Minute*5, time.Second*10)

	// make email sender
	mailSender := mailsender.NewMailSender(
		conf.Email.Email,
		conf.Email.Password,
		conf.Email.Host,
		conf.Email.Port,
	)

	// init layers
	srvc, err := services.NewServiceManager(*conf)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	userService := user_service.NewUserService(
		srvc,
		sJWT,
		memCache,
		mailSender,
	)
	movieService := movie_service.NewMovieService(
		srvc,
	)
	bookingService := booking_service.NewBookingService(
		srvc,
	)

	handler := handlers.NewHandler(
		srvc,
		userService,
		movieService,
		bookingService,
		logger,
	)

	router := api.NewRouter(&api.Option{
		Conf:       conf,
		Middleware: middleware.NewMiddlware(sJWT),
		Handler:    handler,
	})

	// init servers
	server := &http.Server{
		Addr:    conf.HTTP.Port,
		Handler: router,
	}

	// start server
	go func() {
		logger.Info(fmt.Sprintf("Server started localhost:%s started", server.Addr))
		err := server.ListenAndServe()
		if err != nil {
			logger.Error("app: ", zap.Error(err))
			return
		}
	}()

	// gracefull shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	logger.Info("Shutdown server started")
	stopCtx, stopCancle := context.WithTimeout(context.Background(), time.Second*5)
	defer stopCancle()

	server.Shutdown(stopCtx)

	logger.Info("Server shutdown completed")
}
