package usecase

import (
	"errors"
	"task9/domain"
	"time"
)

type TaskUseCase struct {
	taskRepo domain.TaskRepository
}

func NewTaskUseCase(taskRepo domain.TaskRepository) *TaskUseCase {
	return &TaskUseCase{taskRepo: taskRepo}
}

func (uc *TaskUseCase) GetAllTasks() ([]domain.Task, error) {
	return uc.taskRepo.GetAll()
}

func (uc *TaskUseCase) GetTaskByID(id string) (domain.Task, error) {
	return uc.taskRepo.GetByID(id)
}

func (uc *TaskUseCase) CreateTask(req domain.CreateTaskRequest) (domain.Task, error) {
	status := req.Status
	if status == "" {
		status = "pending"
	}

	if status != "pending" && status != "in_progress" && status != "completed" {
		return domain.Task{}, errors.New("invalid status")
	}

	now := time.Now()
	task := domain.Task{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return uc.taskRepo.Create(task)
}

func (uc *TaskUseCase) UpdateTask(id string, req domain.UpdateTaskRequest) (domain.Task, error) {
	existingTask, err := uc.taskRepo.GetByID(id)
	if err != nil {
		return domain.Task{}, err
	}

	if req.Title != "" {
		existingTask.Title = req.Title
	}
	if req.Description != "" {
		existingTask.Description = req.Description
	}
	if !req.DueDate.IsZero() {
		existingTask.DueDate = req.DueDate
	}
	if req.Status != "" {
		if req.Status != "pending" && req.Status != "in_progress" && req.Status != "completed" {
			return domain.Task{}, errors.New("invalid status")
		}
		existingTask.Status = req.Status
	}

	existingTask.UpdatedAt = time.Now()

	return uc.taskRepo.Update(id, existingTask)
}

func (uc *TaskUseCase) DeleteTask(id string) error {
	return uc.taskRepo.Delete(id)
}
