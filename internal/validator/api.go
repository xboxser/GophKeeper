package validator

import (
	"errors"
	"gophkeeper/internal/model"

	"github.com/go-playground/validator/v10"
)

// ValidateModelUserAPI - Проверяем обязательные поля APIUser
func ValidateModelUserAPI(user model.APIUser) error {
	if Validate == nil {
		Init()
	}
	err := Validate.Struct(user)
	if err != nil {
		// Скрываем ошибки формата "Key: 'ApiUser.Login' Error:Field validation for 'Login' failed on the 'required' tag"
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				switch fieldError.Field() {
				case "Login":
					return errors.New("login is required")
				case "Password":
					return errors.New("password is required")
				}
			}
		} else {
			return errors.New("invalid input data")
		}

	}
	return nil
}

// ValidateModelCredentialAPI - Проверяем обязательные поля Credential
func ValidateModelCredentialAPI(credential model.Credential) error {
	if Validate == nil {
		Init()
	}
	err := Validate.Struct(credential)
	if err != nil {
		// Скрываем ошибки формата "Key: 'Credential.Login' Error:Field validation for 'Login' failed on the 'required' tag"
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				switch fieldError.Field() {
				case "Login":
					return errors.New("login is required")
				case "Password":
					return errors.New("password is required")
				}
			}
		} else {
			return errors.New("invalid input data")
		}

	}
	return nil
}
