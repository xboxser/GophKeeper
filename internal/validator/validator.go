package validator

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

// https://habr.com/ru/companies/first/articles/927460/
var Validate *validator.Validate
var onceValidator sync.Once

// Инициализируем валидатор один раз на все приложение
func Init() {
	onceValidator.Do(func() {
		Validate = validator.New()
	})
}
