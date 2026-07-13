package middleware

import (
	"context"
	"strings"

	mycontext "github.com/abozorov/cinema/cmd/api_gateway/internal/my_context"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		if !strings.HasPrefix(token, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := m.jwt.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		if claims.Role != "admin" {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		ctx := context.WithValue(c.Request.Context(), mycontext.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, mycontext.EmailKey, claims.Email)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func (m *Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		if !strings.HasPrefix(token, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := m.jwt.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		ctx := context.WithValue(c.Request.Context(), mycontext.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, mycontext.EmailKey, claims.Email)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
