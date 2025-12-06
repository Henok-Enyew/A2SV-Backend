package usecase

import (
	"errors"
	"task9/domain"
)

type AuthUseCase struct {
	userRepo      domain.UserRepository
	passwordHasher domain.PasswordHasher
	tokenGenerator domain.TokenGenerator
}

func NewAuthUseCase(userRepo domain.UserRepository, passwordHasher domain.PasswordHasher, tokenGenerator domain.TokenGenerator) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
	}
}

func (uc *AuthUseCase) Register(req domain.RegisterRequest) (domain.User, error) {
	if len(req.Password) < 6 {
		return domain.User{}, errors.New("password must be at least 6 characters")
	}

	isFirst, err := uc.userRepo.IsFirstUser()
	if err != nil {
		return domain.User{}, err
	}

	role := "user"
	if isFirst {
		role = "admin"
	}

	hashedPassword, err := uc.passwordHasher.Hash(req.Password)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		Username: req.Username,
		Password: hashedPassword,
		Role:     role,
	}

	return uc.userRepo.Create(user)
}

func (uc *AuthUseCase) Login(req domain.LoginRequest) (string, domain.User, error) {
	user, err := uc.userRepo.GetByUsername(req.Username)
	if err != nil {
		return "", domain.User{}, errors.New("invalid credentials")
	}

	if !uc.passwordHasher.Compare(user.Password, req.Password) {
		return "", domain.User{}, errors.New("invalid credentials")
	}

	token, err := uc.tokenGenerator.Generate(user.ID, user.Username, user.Role)
	if err != nil {
		return "", domain.User{}, err
	}

	user.Password = ""
	return token, user, nil
}

func (uc *AuthUseCase) PromoteUser(username string) error {
	return uc.userRepo.UpdateRole(username, "admin")
}

