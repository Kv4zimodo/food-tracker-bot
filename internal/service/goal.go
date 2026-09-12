package service

import (
	"context"
	"errors"

	"github.com/kv4zimodo/food-tracker-bot/internal/models"
	"github.com/kv4zimodo/food-tracker-bot/internal/repository"
)

type goalService struct {
	repository repository.GoalRepository
}

type GoalService interface {
	CreateGoal(ctx context.Context, goal models.Goal) (*models.Goal, error)
	DeleteGoal(ctx context.Context, userID int64) error
	GetGoalByUserID(ctx context.Context, userID int64) (*models.Goal, error)
	UpdateGoal(ctx context.Context, userID int64, goal models.Goal) (*models.Goal, error)
}

func NewGoalService(repository repository.GoalRepository) GoalService {
	return &goalService{
		repository: repository,
	}
}

func validateGoal(goal models.Goal) error {
	if goal.Calories < 0 {
		return errors.New("Калории не могут быть меньше нуля")
	}
	if goal.Calories >= 20000 {
		return errors.New("Куда тебе столько?!")
	}
	if goal.Protein < 0 {
		return errors.New("Белки не могут быть меньше нуля")
	}
	if goal.Fat < 0 {
		return errors.New("Жиры не могут быть меньше нуля")
	}
	if goal.Carbs < 0 {
		return errors.New("Углеводы не могут быть меньше нуля")
	}
	return nil
}

func (g *goalService) CreateGoal(ctx context.Context, goal models.Goal) (*models.Goal, error) {
	if goal.UserID <= 0 {
		return nil, errors.New("Некорректный ID")
	}
	if err := validateGoal(goal); err != nil {
		return nil, err
	}
	return g.repository.CreateGoal(ctx, goal)
}

func (g *goalService) DeleteGoal(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return errors.New("Некорректный ID")
	}
	return g.repository.DeleteGoal(ctx, userID)
}

func (g *goalService) GetGoalByUserID(ctx context.Context, userID int64) (*models.Goal, error) {
	if userID <= 0 {
		return nil, errors.New("Некорректный ID")
	}
	return g.repository.GetGoalByUserID(ctx, userID)
}

func (g *goalService) UpdateGoal(ctx context.Context, userID int64, goal models.Goal) (*models.Goal, error) {
	if userID <= 0 {
		return nil, errors.New("Некорректный ID")
	}
	if err := validateGoal(goal); err != nil {
		return nil, err
	}
	return g.repository.UpdateGoal(ctx, userID, goal)
}
