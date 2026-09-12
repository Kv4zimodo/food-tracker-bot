package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kv4zimodo/food-tracker-bot/internal/models"
)

type goalRepository struct {
	db *pgx.Conn
}

type GoalRepository interface {
	CreateGoal(ctx context.Context, goal models.Goal) (*models.Goal, error)
	DeleteGoal(ctx context.Context, userID int64) error
	GetGoalByUserID(ctx context.Context, userID int64) (*models.Goal, error)
	UpdateGoal(ctx context.Context, userID int64, goal models.Goal) (*models.Goal, error)
}

func NewGoalRepository(db *pgx.Conn) GoalRepository {
	return &goalRepository{
		db: db,
	}
}

func (g *goalRepository) CreateGoal(ctx context.Context, goal models.Goal) (*models.Goal, error) {
	query := `INSERT INTO goals (user_id, calories, protein, fats, carbs)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING user_id, calories, protein, fats, carbs
	`
	err := g.db.QueryRow(ctx, query,
		goal.UserID,
		goal.Calories,
		goal.Protein,
		goal.Fat,
		goal.Carbs).Scan(&goal.UserID,
		&goal.Calories,
		&goal.Protein,
		&goal.Fat,
		&goal.Carbs)
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

func (g *goalRepository) DeleteGoal(ctx context.Context, userID int64) error {
	query := `DELETE FROM goals
	WHERE user_id = $1
	`
	_, err := g.db.Exec(ctx, query, userID)
	return err
}

func (g *goalRepository) GetGoalByUserID(ctx context.Context, userID int64) (*models.Goal, error) {
	query := `SELECT user_id, calories, protein, fats, carbs
	FROM goals
	WHERE user_id = $1
	`
	var goal models.Goal
	err := g.db.QueryRow(ctx, query, userID).Scan(&goal.UserID,
		&goal.Calories,
		&goal.Protein,
		&goal.Fat,
		&goal.Carbs)
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

func (g *goalRepository) UpdateGoal(ctx context.Context, userID int64, goal models.Goal) (*models.Goal, error) {
	query := `UPDATE goals
	SET calories = $1, protein = $2, fats = $3, carbs = $4
	WHERE user_id = $5
	RETURNING user_id, calories, protein, fats, carbs
	`
	err := g.db.QueryRow(ctx, query,
		goal.Calories,
		goal.Protein,
		goal.Fat,
		goal.Carbs,
		userID).Scan(&goal.UserID,
		&goal.Calories,
		&goal.Protein,
		&goal.Fat,
		&goal.Carbs)
	if err != nil {
		return nil, err
	}
	return &goal, nil
}
