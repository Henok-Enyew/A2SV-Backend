package middleware

import (
	"net/http"
	"net/http/httptest"
	"task9/delivery/middleware"
	"task9/infrastructure"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthMiddleware_RequireAuth(t *testing.T) {
	tokenGenerator := infrastructure.NewJWTGenerator()
	authMiddleware := middleware.NewAuthMiddleware(tokenGenerator)

	t.Run("valid token", func(t *testing.T) {
		token, _ := tokenGenerator.Generate("123", "testuser", "user")

		router := setupRouter()
		router.Use(authMiddleware.RequireAuth())
		router.GET("/test", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			username, _ := c.Get("username")
			role, _ := c.Get("role")
			c.JSON(http.StatusOK, gin.H{
				"user_id":  userID,
				"username": username,
				"role":     role,
			})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		router := setupRouter()
		router.Use(authMiddleware.RequireAuth())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid token format", func(t *testing.T) {
		router := setupRouter()
		router.Use(authMiddleware.RequireAuth())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "InvalidFormat token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid token", func(t *testing.T) {
		router := setupRouter()
		router.Use(authMiddleware.RequireAuth())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthMiddleware_RequireAdmin(t *testing.T) {
	tokenGenerator := infrastructure.NewJWTGenerator()
	authMiddleware := middleware.NewAuthMiddleware(tokenGenerator)

	t.Run("admin access", func(t *testing.T) {
		token, _ := tokenGenerator.Generate("123", "adminuser", "admin")

		router := setupRouter()
		router.Use(authMiddleware.RequireAuth())
		router.Use(authMiddleware.RequireAdmin())
		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("user access denied", func(t *testing.T) {
		token, _ := tokenGenerator.Generate("123", "regularuser", "user")

		router := setupRouter()
		router.Use(authMiddleware.RequireAuth())
		router.Use(authMiddleware.RequireAdmin())
		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("no role in context", func(t *testing.T) {
		router := setupRouter()
		router.Use(authMiddleware.RequireAdmin())
		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/admin", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

