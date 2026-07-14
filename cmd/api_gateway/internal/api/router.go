package api

import (
	"github.com/abozorov/cinema/cmd/api_gateway/internal/api/handlers"
	"github.com/abozorov/cinema/cmd/api_gateway/internal/api/middleware"
	"github.com/abozorov/cinema/cmd/api_gateway/internal/config"
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
	authApi.POST("register", opt.Handler.Register) // registration
	authApi.POST("verify", opt.Handler.Verify)     // verify
	authApi.POST("login", opt.Handler.Login)       // login
	// // authApi.POST("refresh")    // refresh access token

	// user
	userApi := router.Group("/api/user")
	userApi.Use(opt.Middleware.Auth())
	userApi.GET("/:id", opt.Handler.GetGyId)      //get by id
	userApi.PATCH("/:id", opt.Handler.UpdateById) //update by id

	// movie
	movieApi := router.Group("/api/movie")
	movieApi.Use(opt.Middleware.Auth())
	movieApi.GET("/:id", opt.Handler.GetMovieById)
	movieApi.GET("", opt.Handler.ListMovies)

	// // booking
	bookingApi := router.Group("/api/booking")
	bookingApi.Use(opt.Middleware.Auth())
	bookingApi.GET("/:id", opt.Handler.GetBooking)
	bookingApi.POST("", opt.Handler.CreateBooking)
	bookingApi.GET("/user/:id", opt.Handler.GetUserBookings)
	bookingApi.DELETE("/:id", opt.Handler.CancelBooking)

	// // admin
	adminApi := router.Group("/api/admin")
	adminApi.Use(opt.Middleware.AuthAdmin())
	adminApi.POST("movie", opt.Handler.CreateMovie)
	adminApi.PATCH("movie/:id", opt.Handler.UpdateMovie)

	return router
}
