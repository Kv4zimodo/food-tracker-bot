package service

import (
	"context"
	"errors"

	"github.com/kv4zimodo/food-tracker-bot/internal/models"
	"github.com/kv4zimodo/food-tracker-bot/internal/repository"
)

type mealService struct {
	repository repository.MealRepository
}

type MealService interface {
	CreateMeal(ctx context.Context, meal models.Meal) (*models.Meal, error)
	DeleteMeal(ctx context.Context, id int64) error
	GetMealByID(ctx context.Context, id int64) (*models.Meal, error)
	GetAllMeals(ctx context.Context) ([]models.Meal, error)
	UpdateMeal(ctx context.Context, id int64, meal models.Meal) (*models.Meal, error)
}

func NewMealService(repository repository.MealRepository) MealService {
	return &mealService{
		repository: repository,
	}
}

func validateMeal(meal models.Meal) error {
	if meal.Category != models.Breakfast &&
		meal.Category != models.Lunch &&
		meal.Category != models.Dinner &&
		meal.Category != models.Snack {
		return errors.New("Некорректная категория приема пищи")
	}
	return nil
}

func (m *mealService) CreateMeal(ctx context.Context, meal models.Meal) (*models.Meal, error) {
	if err := validateMeal(meal); err != nil {
		return nil, err
	}
	return m.repository.CreateMeal(ctx, meal)
}

func (m *mealService) DeleteMeal(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("Некорректный ID")
	}
	return m.repository.DeleteMeal(ctx, id)
}

func (m *mealService) GetMealByID(ctx context.Context, id int64) (*models.Meal, error) {
	if id <= 0 {
		return nil, errors.New("Некорректный ID")
	}
	return m.repository.GetMealByID(ctx, id)
}

func (m *mealService) GetAllMeals(ctx context.Context) ([]models.Meal, error) {
	return m.repository.GetAllMeals(ctx)
}

func (m *mealService) UpdateMeal(ctx context.Context, id int64, meal models.Meal) (*models.Meal, error) {
	if err := validateMeal(meal); err != nil {
		return nil, err
	}
	if id <= 0 {
		return nil, errors.New("Некорректный ID")
	}
	return m.repository.UpdateMeal(ctx, id, meal)
}
