package repositories_integration

import (
	"context"
	"os"
	"task9/domain"
	"task9/infrastructure"
	"task9/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func setupTestUserDB(t *testing.T) (*mongo.Collection, func()) {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	err := infrastructure.ConnectDB(mongoURI, "task_manager_test")
	require.NoError(t, err)

	collection := infrastructure.UserCollection

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		collection.DeleteMany(ctx, bson.M{})
		infrastructure.DisconnectDB()
	}

	return collection, cleanup
}

func TestUserRepository_Integration(t *testing.T) {
	collection, cleanup := setupTestUserDB(t)
	defer cleanup()

	userRepo := repository.NewUserRepositoryMongo(collection)

	t.Run("Create and GetByUsername user", func(t *testing.T) {
		user := domain.User{
			Username: "integration_user",
			Password: "hashed_password",
			Role:     "user",
		}

		createdUser, err := userRepo.Create(user)
		require.NoError(t, err)
		assert.NotEmpty(t, createdUser.ID)
		assert.Equal(t, "integration_user", createdUser.Username)

		retrievedUser, err := userRepo.GetByUsername("integration_user")
		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, retrievedUser.ID)
		assert.Equal(t, "integration_user", retrievedUser.Username)
	})

	t.Run("IsFirstUser - true when empty", func(t *testing.T) {
		isFirst, err := userRepo.IsFirstUser()
		require.NoError(t, err)
		assert.True(t, isFirst)
	})

	t.Run("IsFirstUser - false when users exist", func(t *testing.T) {
		user := domain.User{
			Username: "existing_user",
			Password: "hashed_password",
			Role:     "user",
		}

		_, err := userRepo.Create(user)
		require.NoError(t, err)

		isFirst, err := userRepo.IsFirstUser()
		require.NoError(t, err)
		assert.False(t, isFirst)
	})

	t.Run("UpdateRole", func(t *testing.T) {
		user := domain.User{
			Username: "user_to_promote",
			Password: "hashed_password",
			Role:     "user",
		}

		createdUser, err := userRepo.Create(user)
		require.NoError(t, err)
		assert.Equal(t, "user", createdUser.Role)

		err = userRepo.UpdateRole("user_to_promote", "admin")
		require.NoError(t, err)

		updatedUser, err := userRepo.GetByUsername("user_to_promote")
		require.NoError(t, err)
		assert.Equal(t, "admin", updatedUser.Role)
	})

	t.Run("GetByID", func(t *testing.T) {
		user := domain.User{
			Username: "user_by_id",
			Password: "hashed_password",
			Role:     "user",
		}

		createdUser, err := userRepo.Create(user)
		require.NoError(t, err)

		retrievedUser, err := userRepo.GetByID(createdUser.ID)
		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, retrievedUser.ID)
		assert.Equal(t, "user_by_id", retrievedUser.Username)
	})

	t.Run("Create duplicate username", func(t *testing.T) {
		user := domain.User{
			Username: "duplicate_user",
			Password: "hashed_password",
			Role:     "user",
		}

		_, err := userRepo.Create(user)
		require.NoError(t, err)

		_, err = userRepo.Create(user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username already exists")
	})

	t.Run("GetByUsername - user not found", func(t *testing.T) {
		_, err := userRepo.GetByUsername("nonexistent_user")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("GetByID with invalid ID", func(t *testing.T) {
		_, err := userRepo.GetByID("invalid_id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID format")
	})

	t.Run("UpdateRole - user not found", func(t *testing.T) {
		err := userRepo.UpdateRole("nonexistent_user", "admin")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})
}

