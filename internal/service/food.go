package service

import (
	"context"
	"errors"
	"strings"

	"github.com/kv4zimodo/food-tracker-bot/internal/models"
	"github.com/kv4zimodo/food-tracker-bot/internal/repository"
)

type foodService struct {
	repository repository.FoodRepository
}

type FoodService interface {
	CreateFood(ctx context.Context, food models.Food) (*models.Food, error)
	DeleteFood(ctx context.Context, id int64) error
	GetFoodByID(ctx context.Context, id int64) (*models.Food, error)
	GetAllFoods(ctx context.Context) ([]models.Food, error)
	UpdateFood(ctx context.Context, id int64, food models.Food) (*models.Food, error)
}

func NewFoodService(repository repository.FoodRepository) FoodService {
	return &foodService{
		repository: repository,
	}
}

func validateFood(food models.Food) error {
	name := strings.TrimSpace(food.Name)
	if name == "" {
		return errors.New("Название продукта не может быть пустым")
	}
	if food.Calories < 0 {
		return errors.New("Калории не могут быть меньше нуля")
	}
	if food.Calories >= 20000 {
		return errors.New("Куда тебе столько?!")
	}
	if food.Protein < 0 {
		return errors.New("Белки не могут быть меньше нуля")
	}
	if food.Fat < 0 {
		return errors.New("Жиры не могут быть меньше нуля")
	}
	if food.Carbs < 0 {
		return errors.New("Углеводы не могут быть меньше нуля")
	}
	return nil
}

func (f *foodService) CreateFood(ctx context.Context, food models.Food) (*models.Food, error) {
	if err := validateFood(food); err != nil {
		return nil, err
	}
	return f.repository.CreateFood(ctx, food)
}

func (f *foodService) DeleteFood(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("Некорректный ID")
	}
	return f.repository.DeleteFood(ctx, id)
}

func (f *foodService) GetFoodByID(ctx context.Context, id int64) (*models.Food, error) {
	if id <= 0 {
		return nil, errors.New("Некорректный ID")
	}
	return f.repository.GetFoodByID(ctx, id)
}

func (f *foodService) GetAllFoods(ctx context.Context) ([]models.Food, error) {
	return f.repository.GetAllFoods(ctx)
}

func (f *foodService) UpdateFood(ctx context.Context, id int64, food models.Food) (*models.Food, error) {
	if id <= 0 {
		return nil, errors.New("Некорректный ID")
	}
	if err := validateFood(food); err != nil {
		return nil, err
	}
	return f.repository.UpdateFood(ctx, id, food)
}
