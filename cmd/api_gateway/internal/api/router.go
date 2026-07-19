package api

import (
	"github.com/abozorov/cinema/cmd/api_gateway/internal/api/handlers"
	"github.com/abozorov/cinema/cmd/api_gateway/internal/api/middleware"
	"github.com/abozorov/cinema/cmd/api_gateway/internal/config"
	"github.com/abozorov/cinema/cmd/api_gateway/internal/models/permission"
	"github.com/gin-gonic/gin"
)

type Option struct {
	Conf       *config.Config
	Middleware *middleware.Middleware
	Handler    *handlers.Handler
}

func NewRouter(opt *Option) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), opt.Middleware.Logging())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// auth
	authApi := router.Group("/api/auth")
	authApi.POST(
		"register",
		opt.Handler.Register,
	)

	authApi.POST(
		"verify",
		opt.Handler.Verify,
	)

	authApi.POST(
		"login",
		opt.Handler.Login,
	)
	// // authApi.POST("refresh")    // refresh access token

	// user
	userApi := router.Group("/api/user")
	userApi.Use(opt.Middleware.Auth())
	userApi.GET(
		"/:id",
		opt.Middleware.RBAC(permission.UserViewMe),
		opt.Handler.GetGyId,
	)

	userApi.PATCH(
		"/:id",
		opt.Middleware.RBAC(permission.UserUpdate),
		opt.Handler.UpdateById,
	)

	// movie
	movieApi := router.Group("/api/movie")
	movieApi.Use(opt.Middleware.Auth())
	movieApi.GET(
		"/:id",
		opt.Middleware.RBAC(permission.MovieView),
		opt.Handler.GetMovieById,
	)

	movieApi.GET(
		"",
		opt.Middleware.RBAC(permission.MovieList),
		opt.Handler.ListMovies,
	)

	// // booking
	bookingApi := router.Group("/api/booking")
	bookingApi.Use(opt.Middleware.Auth())
	bookingApi.GET(
		"/:id",
		opt.Middleware.RBAC(permission.BookingView),
		opt.Handler.GetBooking,
	)

	bookingApi.POST(
		"",
		opt.Middleware.RBAC(permission.BookingCreate),
		opt.Handler.CreateBooking,
	)

	bookingApi.GET(
		"/user/:id",
		opt.Middleware.RBAC(permission.BookingViewMe),
		opt.Handler.GetUserBookings,
	)
	bookingApi.DELETE(
		"/:id",
		opt.Middleware.RBAC(permission.BookingCancel),
		opt.Handler.CancelBooking,
	)

	// // admin
	adminApi := router.Group("/api/admin")
	adminApi.Use(opt.Middleware.Auth())
	adminApi.POST(
		"movie",
		opt.Middleware.RBAC(permission.MovieCreate),
		opt.Handler.CreateMovie,
	)

	adminApi.PATCH(
		"movie/:id",
		opt.Middleware.RBAC(permission.MovieUpdate),
		opt.Handler.UpdateMovie,
	)

	return router
}
