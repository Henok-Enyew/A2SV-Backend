package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"task9/delivery/http"
	"task9/domain"
	"task9/tests/mocks"
	"task9/usecase"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestTaskHandler_GetAllTasks(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	taskUseCase := usecase.NewTaskUseCase(mockTaskRepo)
	taskHandler := http.NewTaskHandler(taskUseCase)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedTasks := []domain.Task{
			{ID: "1", Title: "Task 1", Status: "pending"},
			{ID: "2", Title: "Task 2", Status: "completed"},
		}

		mockTaskRepo.On("GetAll").Return(expectedTasks, nil)

		router := setupTestRouter()
		router.GET("/tasks", taskHandler.GetAllTasks)

		req := httptest.NewRequest("GET", "/tasks", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "success", response["status"])
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockTaskRepo.On("GetAll").Return([]domain.Task{}, errors.New("database error"))

		router := setupTestRouter()
		router.GET("/tasks", taskHandler.GetAllTasks)

		req := httptest.NewRequest("GET", "/tasks", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockTaskRepo.AssertExpectations(t)
	})
}

func TestTaskHandler_GetTaskByID(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	taskUseCase := usecase.NewTaskUseCase(mockTaskRepo)
	taskHandler := http.NewTaskHandler(taskUseCase)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedTask := domain.Task{
			ID:     "123",
			Title:  "Test Task",
			Status: "pending",
		}

		mockTaskRepo.On("GetByID", "123").Return(expectedTask, nil)

		router := setupTestRouter()
		router.GET("/tasks/:id", taskHandler.GetTaskByID)

		req := httptest.NewRequest("GET", "/tasks/123", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "success", response["status"])
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("task not found", func(t *testing.T) {
		mockTaskRepo.On("GetByID", "999").Return(domain.Task{}, errors.New("task not found"))

		router := setupTestRouter()
		router.GET("/tasks/:id", taskHandler.GetTaskByID)

		req := httptest.NewRequest("GET", "/tasks/999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("invalid ID format", func(t *testing.T) {
		mockTaskRepo.On("GetByID", "invalid").Return(domain.Task{}, errors.New("invalid task ID format"))

		router := setupTestRouter()
		router.GET("/tasks/:id", taskHandler.GetTaskByID)

		req := httptest.NewRequest("GET", "/tasks/invalid", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockTaskRepo.AssertExpectations(t)
	})
}

func TestTaskHandler_CreateTask(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	taskUseCase := usecase.NewTaskUseCase(mockTaskRepo)
	taskHandler := http.NewTaskHandler(taskUseCase)

	t.Run("successful creation", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title":       "New Task",
			"description": "Description",
			"due_date":    time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			"status":      "pending",
		}

		jsonBody, _ := json.Marshal(reqBody)

		mockTaskRepo.On("Create", mock.AnythingOfType("domain.Task")).Return(domain.Task{
			ID:     "123",
			Title:  "New Task",
			Status: "pending",
		}, nil)

		router := setupTestRouter()
		router.POST("/tasks", taskHandler.CreateTask)

		req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		router := setupTestRouter()
		router.POST("/tasks", taskHandler.CreateTask)

		req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing required fields", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"description": "Description",
		}

		jsonBody, _ := json.Marshal(reqBody)

		router := setupTestRouter()
		router.POST("/tasks", taskHandler.CreateTask)

		req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAuthHandler_Register(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	authHandler := setupAuthHandler(mockUserRepo)

	t.Run("successful registration", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"username": "newuser",
			"password": "password123",
		}

		jsonBody, _ := json.Marshal(reqBody)

		mockUserRepo.On("IsFirstUser").Return(true, nil)
		mockUserRepo.On("Create", mock.AnythingOfType("domain.User")).Return(domain.User{
			ID:       "123",
			Username: "newuser",
			Role:     "admin",
		}, nil)

		router := setupTestRouter()
		router.POST("/auth/register", authHandler.Register)

		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("password too short", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"username": "user",
			"password": "short",
		}

		jsonBody, _ := json.Marshal(reqBody)

		router := setupTestRouter()
		router.POST("/auth/register", authHandler.Register)

		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	authHandler := setupAuthHandler(mockUserRepo)

	t.Run("successful login", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"username": "testuser",
			"password": "password123",
		}

		jsonBody, _ := json.Marshal(reqBody)

		hashedPassword, _ := setupPasswordHasher().Hash("password123")
		mockUserRepo.On("GetByUsername", "testuser").Return(domain.User{
			ID:       "123",
			Username: "testuser",
			Password: hashedPassword,
			Role:     "user",
		}, nil)

		router := setupTestRouter()
		router.POST("/auth/login", authHandler.Login)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "success", response["status"])
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"username": "testuser",
			"password": "wrongpassword",
		}

		jsonBody, _ := json.Marshal(reqBody)

		hashedPassword, _ := setupPasswordHasher().Hash("correctpassword")
		mockUserRepo.On("GetByUsername", "testuser").Return(domain.User{
			ID:       "123",
			Username: "testuser",
			Password: hashedPassword,
			Role:     "user",
		}, nil)

		router := setupTestRouter()
		router.POST("/auth/login", authHandler.Login)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockUserRepo.AssertExpectations(t)
	})
}

func setupAuthHandler(mockUserRepo *mocks.MockUserRepository) *http.AuthHandler {
	passwordHasher := setupPasswordHasher()
	tokenGenerator := setupTokenGenerator()
	authUseCase := usecase.NewAuthUseCase(mockUserRepo, passwordHasher, tokenGenerator)
	return http.NewAuthHandler(authUseCase)
}

func setupPasswordHasher() domain.PasswordHasher {
	hasher := &mockPasswordHasher{}
	return hasher
}

func setupTokenGenerator() domain.TokenGenerator {
	generator := &mockTokenGenerator{}
	return generator
}

type mockPasswordHasher struct{}

func (m *mockPasswordHasher) Hash(password string) (string, error) {
	return "hashed_" + password, nil
}

func (m *mockPasswordHasher) Compare(hashedPassword, password string) bool {
	return hashedPassword == "hashed_"+password
}

type mockTokenGenerator struct{}

func (m *mockTokenGenerator) Generate(userID, username, role string) (string, error) {
	return "mock_token_" + userID, nil
}

func (m *mockTokenGenerator) Validate(tokenString string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"user_id":  "123",
		"username": "testuser",
		"role":     "user",
	}, nil
}

