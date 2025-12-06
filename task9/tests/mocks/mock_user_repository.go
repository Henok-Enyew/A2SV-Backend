package mocks

import (
	"task9/domain"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (domain.User, error) {
	args := m.Called(username)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id string) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateRole(username string, role string) error {
	args := m.Called(username, role)
	return args.Error(0)
}

func (m *MockUserRepository) IsFirstUser() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

