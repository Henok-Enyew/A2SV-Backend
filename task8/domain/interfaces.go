package domain

type TaskRepository interface {
	GetAll() ([]Task, error)
	GetByID(id string) (Task, error)
	Create(task Task) (Task, error)
	Update(id string, task Task) (Task, error)
	Delete(id string) error
}

type UserRepository interface {
	Create(user User) (User, error)
	GetByUsername(username string) (User, error)
	GetByID(id string) (User, error)
	UpdateRole(username string, role string) error
	IsFirstUser() (bool, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) bool
}

type TokenGenerator interface {
	Generate(userID, username, role string) (string, error)
	Validate(tokenString string) (map[string]interface{}, error)
}

