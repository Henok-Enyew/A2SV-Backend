package domain

import "time"

type Task struct {
	ID          string
	Title       string
	Description string
	DueDate     time.Time
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type User struct {
	ID       string
	Username string
	Password string
	Role     string
}

type RegisterRequest struct {
	Username string
	Password string
}

type LoginRequest struct {
	Username string
	Password string
}

type CreateTaskRequest struct {
	Title       string
	Description string
	DueDate     time.Time
	Status      string
}

type UpdateTaskRequest struct {
	Title       string
	Description string
	DueDate     time.Time
	Status      string
}

type PromoteRequest struct {
	Username string
}

