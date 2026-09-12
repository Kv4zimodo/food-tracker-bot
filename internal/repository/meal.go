package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kv4zimodo/food-tracker-bot/internal/models"
)

type mealRepository struct {
	db *pgx.Conn
}

type MealRepository interface {
	CreateMeal(ctx context.Context, meal models.Meal) (*models.Meal, error)
	DeleteMeal(ctx context.Context, id int64) error
	GetMealByID(ctx context.Context, id int64) (*models.Meal, error)
	GetAllMeals(ctx context.Context) ([]models.Meal, error)
	UpdateMeal(ctx context.Context, id int64, meal models.Meal) (*models.Meal, error)
}

func NewMealRepository(db *pgx.Conn) MealRepository {
	return &mealRepository{
		db: db,
	}
}

func (m *mealRepository) CreateMeal(ctx context.Context, meal models.Meal) (*models.Meal, error) {
	query := `INSERT INTO meals (user_id, category)
	VALUES ($1, $2)
	RETURNING id
	`
	err := m.db.QueryRow(ctx, query,
		meal.UserID,
		meal.Category).Scan(&meal.ID)
	if err != nil {
		return nil, err
	}
	return &meal, nil
}

func (m *mealRepository) DeleteMeal(ctx context.Context, id int64) error {
	query := `DELETE FROM meals
	WHERE id = $1
	`
	_, err := m.db.Exec(ctx, query, id)
	return err
}

func (m *mealRepository) GetMealByID(ctx context.Context, id int64) (*models.Meal, error) {
	query := `SELECT id, user_id, category
	FROM meals
	WHERE id = $1
	`
	var meal models.Meal
	err := m.db.QueryRow(ctx, query, id).Scan(
		&meal.ID,
		&meal.UserID,
		&meal.Category)
	if err != nil {
		return nil, err
	}
	return &meal, nil
}

func (m *mealRepository) GetAllMeals(ctx context.Context) ([]models.Meal, error) {
	query := `SELECT id, user_id, category
	FROM meals
	`
	var meals []models.Meal
	rows, err := m.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var meal models.Meal
		err := rows.Scan(&meal.ID, &meal.UserID, &meal.Category)
		if err != nil {
			return nil, err
		}
		meals = append(meals, meal)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return meals, nil
}

func (m *mealRepository) UpdateMeal(ctx context.Context, id int64, meal models.Meal) (*models.Meal, error) {
	query := `UPDATE meals
	SET user_id = $1, category = $2
	WHERE id = $3
	RETURNING id, user_id, category
	`
	err := m.db.QueryRow(ctx, query,
		meal.UserID,
		meal.Category,
		id).Scan(
		&meal.ID,
		&meal.UserID,
		&meal.Category)
	if err != nil {
		return nil, err
	}
	return &meal, nil
}
