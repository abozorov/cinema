package api

import (
	"github.com/abozorov/cinema/cmd/api_gateway/internal/api/handlers"
	"github.com/abozorov/cinema/cmd/api_gateway/internal/api/middleware"
	"github.com/abozorov/cinema/cmd/api_gateway/internal/config"
	"github.com/gin-gonic/gin"
)

type Option struct {
	Conf        *config.Config
	Midddleware *middleware.Middleware
	Handler     *handlers.Handler
}

func NewRouter(opt *Option) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), opt.Midddleware.Logging())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// auth
	authApi := router.Group("/api/auth")
	authApi.POST("register", opt.Handler.Register) // registration
	authApi.POST("verify", opt.Handler.Verify)     // verify
	// authApi.POST("login")    // login
	// // authApi.POST("refresh")    // login

	// user
	userApi := router.Group("/api/user")
	userApi.GET("/:id", opt.Handler.GetGyId)      //get by id
	userApi.PATCH("/:id", opt.Handler.UpdateById) //update by id

	// movie
	movieApi := router.Group("/api/movie")
	movieApi.POST("", opt.Handler.CreateMovie)
	movieApi.GET("/:id", opt.Handler.GetMovieById)
	movieApi.PUT("/:id", opt.Handler.UpdateMovie)
	movieApi.GET("", opt.Handler.ListMovies)

	// // booking
	bookingApi := router.Group("/api/booking")
	bookingApi.POST("", opt.Handler.CreateBooking)
	bookingApi.GET("/:id", opt.Handler.GetBooking)
	bookingApi.GET("/user/:id", opt.Handler.GetUserBookings)
	bookingApi.DELETE("/:id", opt.Handler.CancelBooking)

	// // admin
	// adminApi := router.Group("/api/admin")

	return router
}
