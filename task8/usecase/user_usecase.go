package usecase

import (
	"errors"
	"task8/domain/entity"
	"task8/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UserUseCase) Register(username, password string) (*entity.User, error) {
	existing, _ := uc.userRepo.FindByUsername(username)
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	count, err := uc.userRepo.Count()
	if err != nil {
		return nil, err
	}

	role := "user"
	if count == 0 {
		role = "admin"
	}

	user := entity.NewUser("", username, string(hashedPassword), role)
	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (uc *UserUseCase) Login(username, password string) (*entity.User, error) {
	user, err := uc.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	user.Password = ""
	return user, nil
}

func (uc *UserUseCase) PromoteUser(username string) error {
	user, err := uc.userRepo.FindByUsername(username)
	if err != nil {
		return errors.New("user not found")
	}

	user.PromoteToAdmin()
	return uc.userRepo.Update(user)
}

func (uc *UserUseCase) GetUserByID(id string) (*entity.User, error) {
	user, err := uc.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}


