package auth

import (
	"errors"
	"golang/internal/user"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(userRepository *user.UserRepository) *AuthService {
	return &AuthService{UserRepository: userRepository}
}

func (service *AuthService) Register(email, password, name string) (string, error) {
	exsistedUser, _ := service.UserRepository.FindByEmail(email)
	if exsistedUser != nil {
		return "", errors.New(ErrUserExists)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user := &user.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}
	_, err = service.UserRepository.Create(user)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

func (service *AuthService) Login(email, password string) (string, error) {
	exsistedUser, _ := service.UserRepository.FindByEmail(email)
	if exsistedUser == nil {
		return "", errors.New(ErrWrongCredetials)
	}
	err := bcrypt.CompareHashAndPassword([]byte(exsistedUser.Password), []byte(password))
	if err != nil {
		return "", errors.New(ErrWrongCredetials)
	}
	return exsistedUser.Email, nil
}
