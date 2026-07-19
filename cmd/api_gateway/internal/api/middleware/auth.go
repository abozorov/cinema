package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/abozorov/cinema/cmd/api_gateway/internal/models/permission"
	mycontext "github.com/abozorov/cinema/cmd/api_gateway/internal/my_context"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) RBAC(requiredPermission string) gin.HandlerFunc {

	return func(c *gin.Context) {
		role := c.Request.Context().Value(mycontext.RoleKey).(string)
		if role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}

		permissions, ok := permission.RolePermission[role]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}

		if _, ok = permissions[requiredPermission]; !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Permission denied",
			})
			return
		}
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
		ctx = context.WithValue(ctx, mycontext.RoleKey, claims.Role)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
