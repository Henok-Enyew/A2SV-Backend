package entity

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

func NewTask(id, title, description string, dueDate time.Time, status string) *Task {
	now := time.Now()
	if status == "" {
		status = "pending"
	}
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (t *Task) Update(title, description string, dueDate time.Time, status string) {
	if title != "" {
		t.Title = title
	}
	if description != "" {
		t.Description = description
	}
	if !dueDate.IsZero() {
		t.DueDate = dueDate
	}
	if status != "" {
		t.Status = status
	}
	t.UpdatedAt = time.Now()
}


