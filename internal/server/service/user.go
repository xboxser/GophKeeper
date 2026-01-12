package service

import (
	"context"
	"fmt"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(context.Context, model.APIUser) (string, error)
	Login(context.Context, model.APIUser) (string, error)
	GetUserForID(context.Context, int) (model.User, error)
}

type userService struct {
	UserRepository repository.UserRepository
	TokenService   TokenService
	SessionService SessionService
}

func NewUserService(userRepository repository.UserRepository, tokenService TokenService, sessionService SessionService) *userService {
	return &userService{
		UserRepository: userRepository,
		TokenService:   tokenService,
		SessionService: sessionService,
	}
}

// Register - сервис регистрации пользователя
func (u *userService) Register(ctx context.Context, apiUser model.APIUser) (string, error) {
	user, err := u.UserRepository.GetUserForLogin(ctx, apiUser.Login)

	if err != nil {
		return "", err
	}

	// Пользователь уже зарегистрирован
	if user.ID != 0 {
		return "", model.ErrLoginBusy
	}

	userID, err := u.UserRepository.RegisterUser(ctx, apiUser)
	if err != nil {
		return "", fmt.Errorf("error registering user: %w", err)
	}

	tokenAuth := model.TokenAuth{
		UserID: userID,
	}

	session, err := u.SessionService.AddSession(ctx, tokenAuth)
	if err != nil {
		return "", err
	}

	token, err := u.TokenService.BuildJWTString(session)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *userService) Login(ctx context.Context, apiUser model.APIUser) (string, error) {
	user, err := u.UserRepository.GetUserForLogin(ctx, apiUser.Login)

	if err != nil {
		return "", err
	}

	if user.ID == 0 {
		return "", model.ErrLoginIncorrect
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(apiUser.Password))
	if err != nil {
		return "", model.ErrLoginIncorrect
	}

	tokenAuth := model.TokenAuth{
		UserID: user.ID,
	}

	session, err := u.SessionService.AddSession(ctx, tokenAuth)
	if err != nil {
		return "", err
	}

	token, err := u.TokenService.BuildJWTString(session)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *userService) GetUserForID(ctx context.Context, userID int) (model.User, error) {
	user, err := u.UserRepository.GetUserForID(ctx, userID)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
