package middleware

import (
	"net/http"
	"strings"
	"task8/infrastructure"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	tokenGenerator *infrastructure.JWTGenerator
}

func NewAuthMiddleware(tokenGenerator *infrastructure.JWTGenerator) *AuthMiddleware {
	return &AuthMiddleware{tokenGenerator: tokenGenerator}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "authorization header required",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := m.tokenGenerator.Validate(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("username", claims["username"])
		c.Set("role", claims["role"])

		c.Next()
	}
}

func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "unauthorized",
			})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

