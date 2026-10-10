package service

import (
	"context"
	"errors"
	"math"

	"github.com/jackc/pgx/v5"
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
	SetGoalByCalories(ctx context.Context, userID int64, calories float64) (*models.Goal, error)
	SetGoalByMacros(ctx context.Context, userID int64, protein float64, fat float64, carbs float64) (*models.Goal, error)
}

func NewGoalService(repository repository.GoalRepository) GoalService {
	return &goalService{
		repository: repository,
	}
}

func validateGoal(goal models.Goal) error {
	if math.IsNaN(goal.Calories) || math.IsInf(goal.Calories, 0) ||
		math.IsNaN(goal.Protein) || math.IsInf(goal.Protein, 0) ||
		math.IsNaN(goal.Fat) || math.IsInf(goal.Fat, 0) ||
		math.IsNaN(goal.Carbs) || math.IsInf(goal.Carbs, 0) {
		return errors.New("Значения цели должны быть обычными числами")
	}
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

func (g *goalService) SetGoalByCalories(ctx context.Context, userID int64, calories float64) (*models.Goal, error) {
	goal := models.Goal{
		UserID:   userID,
		Calories: calories,
		Protein:  calories * 0.25 / 4,
		Fat:      calories * 0.30 / 9,
		Carbs:    calories * 0.45 / 4,
	}

	return g.saveGoal(ctx, goal)
}

func (g *goalService) SetGoalByMacros(ctx context.Context, userID int64, protein float64, fat float64, carbs float64) (*models.Goal, error) {
	goal := models.Goal{
		UserID:   userID,
		Protein:  protein,
		Fat:      fat,
		Carbs:    carbs,
		Calories: 4*protein + 9*fat + 4*carbs,
	}

	return g.saveGoal(ctx, goal)
}

func (g *goalService) saveGoal(ctx context.Context, goal models.Goal) (*models.Goal, error) {
	if goal.UserID <= 0 {
		return nil, errors.New("Некорректный ID")
	}

	if err := validateGoal(goal); err != nil {
		return nil, err
	}

	_, err := g.repository.GetGoalByUserID(ctx, goal.UserID)

	if errors.Is(err, pgx.ErrNoRows) {
		return g.repository.CreateGoal(ctx, goal)
	}

	if err != nil {
		return nil, err
	}

	return g.repository.UpdateGoal(ctx, goal.UserID, goal)
}
