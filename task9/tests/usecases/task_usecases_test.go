package usecases

import (
	"errors"
	"task9/domain"
	"task9/tests/mocks"
	"task9/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTaskUseCase_GetAllTasks(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	taskUseCase := usecase.NewTaskUseCase(mockTaskRepo)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedTasks := []domain.Task{
			{ID: "1", Title: "Task 1", Status: "pending"},
			{ID: "2", Title: "Task 2", Status: "completed"},
		}

		mockTaskRepo.On("GetAll").Return(expectedTasks, nil)

		tasks, err := taskUseCase.GetAllTasks()

		assert.NoError(t, err)
		assert.Len(t, tasks, 2)
		assert.Equal(t, "Task 1", tasks[0].Title)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockTaskRepo.On("GetAll").Return([]domain.Task{}, nil)

		tasks, err := taskUseCase.GetAllTasks()

		assert.NoError(t, err)
		assert.Empty(t, tasks)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockTaskRepo.On("GetAll").Return([]domain.Task{}, errors.New("database error"))

		_, err := taskUseCase.GetAllTasks()

		assert.Error(t, err)
		mockTaskRepo.AssertExpectations(t)
	})
}

func TestTaskUseCase_GetTaskByID(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	taskUseCase := usecase.NewTaskUseCase(mockTaskRepo)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedTask := domain.Task{
			ID:     "123",
			Title:  "Test Task",
			Status: "pending",
		}

		mockTaskRepo.On("GetByID", "123").Return(expectedTask, nil)

		task, err := taskUseCase.GetTaskByID("123")

		assert.NoError(t, err)
		assert.Equal(t, "Test Task", task.Title)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("task not found", func(t *testing.T) {
		mockTaskRepo.On("GetByID", "999").Return(domain.Task{}, errors.New("task not found"))

		_, err := taskUseCase.GetTaskByID("999")

		assert.Error(t, err)
		mockTaskRepo.AssertExpectations(t)
	})
}

func TestTaskUseCase_CreateTask(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	taskUseCase := usecase.NewTaskUseCase(mockTaskRepo)

	t.Run("successful creation with default status", func(t *testing.T) {
		req := domain.CreateTaskRequest{
			Title:       "New Task",
			Description: "Description",
			DueDate:     time.Now().Add(24 * time.Hour),
		}

		mockTaskRepo.On("Create", mock.AnythingOfType("domain.Task")).Return(domain.Task{
			ID:     "123",
			Title:  "New Task",
			Status: "pending",
		}, nil)

		task, err := taskUseCase.CreateTask(req)

		assert.NoError(t, err)
		assert.Equal(t, "pending", task.Status)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("successful creation with custom status", func(t *testing.T) {
		req := domain.CreateTaskRequest{
			Title:       "New Task",
			Description: "Description",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "in_progress",
		}

		mockTaskRepo.On("Create", mock.AnythingOfType("domain.Task")).Return(domain.Task{
			ID:     "123",
			Title:  "New Task",
			Status: "in_progress",
		}, nil)

		task, err := taskUseCase.CreateTask(req)

		assert.NoError(t, err)
		assert.Equal(t, "in_progress", task.Status)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("invalid status", func(t *testing.T) {
		req := domain.CreateTaskRequest{
			Title:       "New Task",
			Description: "Description",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "invalid_status",
		}

		_, err := taskUseCase.CreateTask(req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid status")
	})
}

func TestTaskUseCase_UpdateTask(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	taskUseCase := usecase.NewTaskUseCase(mockTaskRepo)

	t.Run("successful update", func(t *testing.T) {
		existingTask := domain.Task{
			ID:     "123",
			Title:  "Old Title",
			Status: "pending",
		}

		req := domain.UpdateTaskRequest{
			Title:  "New Title",
			Status: "completed",
		}

		mockTaskRepo.On("GetByID", "123").Return(existingTask, nil)
		mockTaskRepo.On("Update", "123", mock.AnythingOfType("domain.Task")).Return(domain.Task{
			ID:     "123",
			Title:  "New Title",
			Status: "completed",
		}, nil)

		task, err := taskUseCase.UpdateTask("123", req)

		assert.NoError(t, err)
		assert.Equal(t, "New Title", task.Title)
		assert.Equal(t, "completed", task.Status)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("task not found", func(t *testing.T) {
		req := domain.UpdateTaskRequest{
			Title: "New Title",
		}

		mockTaskRepo.On("GetByID", "999").Return(domain.Task{}, errors.New("task not found"))

		_, err := taskUseCase.UpdateTask("999", req)

		assert.Error(t, err)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("invalid status", func(t *testing.T) {
		existingTask := domain.Task{
			ID:     "123",
			Title:  "Task",
			Status: "pending",
		}

		req := domain.UpdateTaskRequest{
			Status: "invalid_status",
		}

		mockTaskRepo.On("GetByID", "123").Return(existingTask, nil)

		_, err := taskUseCase.UpdateTask("123", req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid status")
	})
}

func TestTaskUseCase_DeleteTask(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	taskUseCase := usecase.NewTaskUseCase(mockTaskRepo)

	t.Run("successful deletion", func(t *testing.T) {
		mockTaskRepo.On("Delete", "123").Return(nil)

		err := taskUseCase.DeleteTask("123")

		assert.NoError(t, err)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("task not found", func(t *testing.T) {
		mockTaskRepo.On("Delete", "999").Return(errors.New("task not found"))

		err := taskUseCase.DeleteTask("999")

		assert.Error(t, err)
		mockTaskRepo.AssertExpectations(t)
	})
}

