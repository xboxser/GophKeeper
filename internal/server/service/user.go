package service

import (
	"context"
	"fmt"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/repository"
)

type UserService interface {
	Register(context.Context, model.APIUser) (int, error)
}

type userService struct {
	UserRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *userService {
	return &userService{
		UserRepository: userRepository,
	}
}

// Register - сервис регистрации пользователя
func (u *userService) Register(ctx context.Context, apiUser model.APIUser) (int, error) {
	user, err := u.UserRepository.GetUserForLogin(ctx, apiUser.Login)

	if err != nil {
		return 0, err
	}

	// Пользователь уже зарегистрирован
	if user.ID != 0 {
		return 0, model.ErrLoginBusy
	}

	userID, err := u.UserRepository.RegisterUser(ctx, apiUser.Login, apiUser.Password)
	if err != nil {
		return 0, fmt.Errorf("error registering user: %w", err)
	}

	return userID, nil
}
