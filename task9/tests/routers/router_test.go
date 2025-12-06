package routers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task9/delivery/http"
	"task9/delivery/middleware"
	"task9/infrastructure"
	"task9/tests/mocks"
	"task9/usecase"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouterWithMocks() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	mockTaskRepo := new(mocks.MockTaskRepository)
	taskUseCase := usecase.NewTaskUseCase(mockTaskRepo)
	taskHandler := http.NewTaskHandler(taskUseCase)

	mockUserRepo := new(mocks.MockUserRepository)
	passwordHasher := infrastructure.NewBcryptHasher()
	tokenGenerator := infrastructure.NewJWTGenerator()
	authUseCase := usecase.NewAuthUseCase(mockUserRepo, passwordHasher, tokenGenerator)
	authHandler := http.NewAuthHandler(authUseCase)

	authMiddleware := middleware.NewAuthMiddleware(tokenGenerator)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Task Management API",
			"version": "1.0.0",
		})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	protected := r.Group("/")
	protected.Use(authMiddleware.RequireAuth())
	{
		protected.GET("/tasks", taskHandler.GetAllTasks)
		protected.GET("/tasks/:id", taskHandler.GetTaskByID)

		admin := protected.Group("/")
		admin.Use(authMiddleware.RequireAdmin())
		{
			admin.POST("/tasks", taskHandler.CreateTask)
			admin.PUT("/tasks/:id", taskHandler.UpdateTask)
			admin.DELETE("/tasks/:id", taskHandler.DeleteTask)
			admin.POST("/promote", authHandler.PromoteUser)
		}
	}

	return r
}

func TestRouter_PublicEndpoints(t *testing.T) {
	router := setupTestRouterWithMocks()

	t.Run("GET /", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Task Management API", response["message"])
	})

	t.Run("POST /auth/register", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"username": "testuser",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusBadRequest)
	})
}

func TestRouter_AuthenticatedEndpoints(t *testing.T) {
	router := setupTestRouterWithMocks()
	tokenGenerator := infrastructure.NewJWTGenerator()

	t.Run("GET /tasks without token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tasks", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("GET /tasks with valid token", func(t *testing.T) {
		token, _ := tokenGenerator.Generate("123", "testuser", "user")

		req := httptest.NewRequest("GET", "/tasks", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
	})
}

func TestRouter_AdminEndpoints(t *testing.T) {
	router := setupTestRouterWithMocks()
	tokenGenerator := infrastructure.NewJWTGenerator()

	t.Run("POST /tasks as user (forbidden)", func(t *testing.T) {
		token, _ := tokenGenerator.Generate("123", "testuser", "user")

		reqBody := map[string]interface{}{
			"title":       "New Task",
			"description": "Description",
			"due_date":    time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("POST /tasks as admin (authorized)", func(t *testing.T) {
		token, _ := tokenGenerator.Generate("123", "adminuser", "admin")

		reqBody := map[string]interface{}{
			"title":       "New Task",
			"description": "Description",
			"due_date":    time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusBadRequest || w.Code == http.StatusInternalServerError)
	})
}

