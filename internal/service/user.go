package service

import (
	"context"
	"errors"
	"strings"

	"github.com/kv4zimodo/food-tracker-bot/internal/models"
	"github.com/kv4zimodo/food-tracker-bot/internal/repository"
)

type userService struct {
	repository repository.UserRepository
}

type UserService interface {
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id int64) error
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	UpdateUser(ctx context.Context, id int64, user models.User) (*models.User, error)
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{
		repository: repository,
	}
}

func (u *userService) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	name := strings.TrimSpace(user.Name)
	if name == "" {
		return nil, errors.New("Имя пользователя не может быть пустым")
	}
	return u.repository.CreateUser(ctx, user)
}

func (u *userService) DeleteUser(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("Неккоректный ID")
	}
	return u.repository.DeleteUser(ctx, id)
}

func (u *userService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	if id <= 0 {
		return nil, errors.New("Неккоректный ID")
	}
	return u.repository.GetUserByID(ctx, id)
}

func (u *userService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return u.repository.GetAllUsers(ctx)
}

func (u *userService) UpdateUser(ctx context.Context, id int64, user models.User) (*models.User, error) {
	name := strings.TrimSpace(user.Name)
	if name == "" {
		return nil, errors.New("Имя пользователя не может быть пустым")
	}
	if id <= 0 {
		return nil, errors.New("Некорректный ID")
	}
	return u.repository.UpdateUser(ctx, id, user)
}
