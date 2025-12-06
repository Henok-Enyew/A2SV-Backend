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

func setupTestDB(t *testing.T) (*mongo.Collection, func()) {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	err := infrastructure.ConnectDB(mongoURI, "task_manager_test")
	require.NoError(t, err)

	collection := infrastructure.TaskCollection

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		collection.DeleteMany(ctx, bson.M{})
		infrastructure.DisconnectDB()
	}

	return collection, cleanup
}

func TestTaskRepository_Integration(t *testing.T) {
	collection, cleanup := setupTestDB(t)
	defer cleanup()

	taskRepo := repository.NewTaskRepositoryMongo(collection)

	t.Run("Create and Get task", func(t *testing.T) {
		task := domain.Task{
			Title:       "Integration Test Task",
			Description: "Test Description",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "pending",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		createdTask, err := taskRepo.Create(task)
		require.NoError(t, err)
		assert.NotEmpty(t, createdTask.ID)
		assert.Equal(t, "Integration Test Task", createdTask.Title)

		retrievedTask, err := taskRepo.GetByID(createdTask.ID)
		require.NoError(t, err)
		assert.Equal(t, createdTask.ID, retrievedTask.ID)
		assert.Equal(t, "Integration Test Task", retrievedTask.Title)
	})

	t.Run("GetAll tasks", func(t *testing.T) {
		task1 := domain.Task{
			Title:       "Task 1",
			Description: "Description 1",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "pending",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		task2 := domain.Task{
			Title:       "Task 2",
			Description: "Description 2",
			DueDate:     time.Now().Add(48 * time.Hour),
			Status:      "completed",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		_, err := taskRepo.Create(task1)
		require.NoError(t, err)

		_, err = taskRepo.Create(task2)
		require.NoError(t, err)

		tasks, err := taskRepo.GetAll()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(tasks), 2)
	})

	t.Run("Update task", func(t *testing.T) {
		task := domain.Task{
			Title:       "Original Title",
			Description: "Original Description",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "pending",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		createdTask, err := taskRepo.Create(task)
		require.NoError(t, err)

		updatedTask := createdTask
		updatedTask.Title = "Updated Title"
		updatedTask.Status = "completed"

		result, err := taskRepo.Update(createdTask.ID, updatedTask)
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", result.Title)
		assert.Equal(t, "completed", result.Status)
	})

	t.Run("Delete task", func(t *testing.T) {
		task := domain.Task{
			Title:       "Task to Delete",
			Description: "Will be deleted",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "pending",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		createdTask, err := taskRepo.Create(task)
		require.NoError(t, err)

		err = taskRepo.Delete(createdTask.ID)
		require.NoError(t, err)

		_, err = taskRepo.GetByID(createdTask.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("GetByID with invalid ID", func(t *testing.T) {
		_, err := taskRepo.GetByID("invalid_id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid task ID format")
	})

	t.Run("Update non-existent task", func(t *testing.T) {
		task := domain.Task{
			Title:  "Non-existent",
			Status: "pending",
		}

		_, err := taskRepo.Update("507f1f77bcf86cd799439011", task)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

