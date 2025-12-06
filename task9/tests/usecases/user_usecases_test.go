package usecases

import (
	"errors"
	"task9/domain"
	"task9/infrastructure"
	"task9/tests/mocks"
	"task9/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthUseCase_Register(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	passwordHasher := infrastructure.NewBcryptHasher()
	tokenGenerator := infrastructure.NewJWTGenerator()

	authUseCase := usecase.NewAuthUseCase(mockUserRepo, passwordHasher, tokenGenerator)

	t.Run("successful registration - first user becomes admin", func(t *testing.T) {
		req := domain.RegisterRequest{
			Username: "firstuser",
			Password: "password123",
		}

		mockUserRepo.On("IsFirstUser").Return(true, nil)
		mockUserRepo.On("Create", mock.AnythingOfType("domain.User")).Return(domain.User{
			ID:       "123",
			Username: "firstuser",
			Role:     "admin",
		}, nil)

		user, err := authUseCase.Register(req)

		assert.NoError(t, err)
		assert.Equal(t, "admin", user.Role)
		assert.Equal(t, "firstuser", user.Username)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("successful registration - regular user", func(t *testing.T) {
		req := domain.RegisterRequest{
			Username: "regularuser",
			Password: "password123",
		}

		mockUserRepo.On("IsFirstUser").Return(false, nil)
		mockUserRepo.On("Create", mock.AnythingOfType("domain.User")).Return(domain.User{
			ID:       "456",
			Username: "regularuser",
			Role:     "user",
		}, nil)

		user, err := authUseCase.Register(req)

		assert.NoError(t, err)
		assert.Equal(t, "user", user.Role)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("password too short", func(t *testing.T) {
		req := domain.RegisterRequest{
			Username: "user",
			Password: "short",
		}

		_, err := authUseCase.Register(req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must be at least 6 characters")
	})

	t.Run("username already exists", func(t *testing.T) {
		req := domain.RegisterRequest{
			Username: "existinguser",
			Password: "password123",
		}

		mockUserRepo.On("IsFirstUser").Return(false, nil)
		mockUserRepo.On("Create", mock.AnythingOfType("domain.User")).Return(domain.User{}, errors.New("username already exists"))

		_, err := authUseCase.Register(req)

		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestAuthUseCase_Login(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	passwordHasher := infrastructure.NewBcryptHasher()
	tokenGenerator := infrastructure.NewJWTGenerator()

	authUseCase := usecase.NewAuthUseCase(mockUserRepo, passwordHasher, tokenGenerator)

	t.Run("successful login", func(t *testing.T) {
		req := domain.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}

		hashedPassword, _ := passwordHasher.Hash("password123")

		mockUserRepo.On("GetByUsername", "testuser").Return(domain.User{
			ID:       "123",
			Username: "testuser",
			Password: hashedPassword,
			Role:     "user",
		}, nil)

		token, user, err := authUseCase.Login(req)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, "testuser", user.Username)
		assert.Empty(t, user.Password)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		req := domain.LoginRequest{
			Username: "nonexistent",
			Password: "password123",
		}

		mockUserRepo.On("GetByUsername", "nonexistent").Return(domain.User{}, errors.New("user not found"))

		_, _, err := authUseCase.Login(req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid credentials")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("wrong password", func(t *testing.T) {
		req := domain.LoginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		}

		hashedPassword, _ := passwordHasher.Hash("correctpassword")

		mockUserRepo.On("GetByUsername", "testuser").Return(domain.User{
			ID:       "123",
			Username: "testuser",
			Password: hashedPassword,
			Role:     "user",
		}, nil)

		_, _, err := authUseCase.Login(req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid credentials")
		mockUserRepo.AssertExpectations(t)
	})
}

func TestAuthUseCase_PromoteUser(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	passwordHasher := infrastructure.NewBcryptHasher()
	tokenGenerator := infrastructure.NewJWTGenerator()

	authUseCase := usecase.NewAuthUseCase(mockUserRepo, passwordHasher, tokenGenerator)

	t.Run("successful promotion", func(t *testing.T) {
		mockUserRepo.On("UpdateRole", "testuser", "admin").Return(nil)

		err := authUseCase.PromoteUser("testuser")

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockUserRepo.On("UpdateRole", "nonexistent", "admin").Return(errors.New("user not found"))

		err := authUseCase.PromoteUser("nonexistent")

		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}

