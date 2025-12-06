package repository

import "task8/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error
	FindByUsername(username string) (*entity.User, error)
	FindByID(id string) (*entity.User, error)
	Update(user *entity.User) error
	Count() (int64, error)
}


