package repository

import "task9/domain/entity"

type TaskRepository interface {
	Create(task *entity.Task) error
	FindByID(id string) (*entity.Task, error)
	FindAll() ([]*entity.Task, error)
	Update(task *entity.Task) error
	Delete(id string) error
}


