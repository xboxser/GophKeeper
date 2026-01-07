package repository

type UserRepository interface {
}

type userRepository struct {
}

func NewUserRepository() *userRepository {
	return &userRepository{}
}
