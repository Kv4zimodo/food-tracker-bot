package service

import (
	"context"
	"errors"

	"github.com/kv4zimodo/food-tracker-bot/internal/models"
	"github.com/kv4zimodo/food-tracker-bot/internal/repository"
)

type mealItemService struct {
	repository repository.MealItemRepository
}

type MealItemService interface {
	CreateMealItem(ctx context.Context, mealItem models.MealItem) (*models.MealItem, error)
	DeleteMealItem(ctx context.Context, id int64) error
	GetMealItemByID(ctx context.Context, id int64) (*models.MealItem, error)
	GetAllMealItems(ctx context.Context) ([]models.MealItem, error)
	UpdateMealItem(ctx context.Context, id int64, mealItem models.MealItem) (*models.MealItem, error)
}

func NewMealItemService(repository repository.MealItemRepository) MealItemService {
	return &mealItemService{
		repository: repository,
	}
}

func (m *mealItemService) CreateMealItem(ctx context.Context, mealItem models.MealItem) (*models.MealItem, error) {
	if mealItem.Weight <= 0 {
		return nil, errors.New("Вес должен быть больше нуля")
	}
	return m.repository.CreateMealItem(ctx, mealItem)
}

func (m *mealItemService) DeleteMealItem(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("Некорректный ID")
	}
	return m.repository.DeleteMealItem(ctx, id)
}

func (m *mealItemService) GetMealItemByID(ctx context.Context, id int64) (*models.MealItem, error) {
	if id <= 0 {
		return nil, errors.New("Некорректный ID")
	}
	return m.repository.GetMealItemByID(ctx, id)
}

func (m *mealItemService) GetAllMealItems(ctx context.Context) ([]models.MealItem, error) {
	return m.repository.GetAllMealItems(ctx)
}

func (m *mealItemService) UpdateMealItem(ctx context.Context, id int64, mealItem models.MealItem) (*models.MealItem, error) {
	if mealItem.Weight <= 0 {
		return nil, errors.New("Вес должен быть больше нуля")
	}
	if id <= 0 {
		return nil, errors.New("Некорректный ID")
	}
	return m.repository.UpdateMealItem(ctx, id, mealItem)
}
