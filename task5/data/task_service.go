package data

import (
	"errors"
	"sync"
	"task5/models"
	"time"
)

type TaskService struct {
	tasks map[int]models.Task
	mutex sync.RWMutex
	nextID int
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (ts *TaskService) GetAllTasks() []models.Task {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()

	tasks := make([]models.Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (ts *TaskService) GetTaskByID(id int) (models.Task, error) {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()

	task, exists := ts.tasks[id]
	if !exists {
		return models.Task{}, errors.New("task not found")
	}
	return task, nil
}

func (ts *TaskService) CreateTask(req models.CreateTaskRequest) models.Task {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	now := time.Now()
	status := req.Status
	if status == "" {
		status = "pending"
	}

	task := models.Task{
		ID:          ts.nextID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	ts.tasks[ts.nextID] = task
	ts.nextID++

	return task
}

func (ts *TaskService) UpdateTask(id int, req models.UpdateTaskRequest) (models.Task, error) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	task, exists := ts.tasks[id]
	if !exists {
		return models.Task{}, errors.New("task not found")
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if !req.DueDate.IsZero() {
		task.DueDate = req.DueDate
	}
	if req.Status != "" {
		task.Status = req.Status
	}

	task.UpdatedAt = time.Now()
	ts.tasks[id] = task

	return task, nil
}

func (ts *TaskService) DeleteTask(id int) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	if _, exists := ts.tasks[id]; !exists {
		return errors.New("task not found")
	}

	delete(ts.tasks, id)
	return nil
}

