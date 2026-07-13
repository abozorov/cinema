package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) Logging() gin.HandlerFunc {
	return func(c *gin.Context) {

		startTime := time.Now()
		log.Printf("[INFO]	START http.method: {%s} url: {%s}\n",
			c.Request.Method,
			c.Request.URL.Path,
		)

		c.Next()

		endTime := time.Now()
		log.Printf("[INFO]	END request duration: {%d} ms\n\n",
			int(endTime.UnixMilli())-int(startTime.UnixMilli()))
	}
}
