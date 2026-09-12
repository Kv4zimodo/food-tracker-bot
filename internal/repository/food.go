package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kv4zimodo/food-tracker-bot/internal/models"
)

type foodRepository struct {
	db *pgx.Conn
}

type FoodRepository interface {
	CreateFood(ctx context.Context, food models.Food) (*models.Food, error)
	DeleteFood(ctx context.Context, id int64) error
	GetFoodByID(ctx context.Context, id int64) (*models.Food, error)
	GetAllFoods(ctx context.Context) ([]models.Food, error)
	UpdateFood(ctx context.Context, id int64, food models.Food) (*models.Food, error)
}

func NewFoodRepository(db *pgx.Conn) FoodRepository {
	return &foodRepository{
		db: db,
	}
}

func (f *foodRepository) CreateFood(ctx context.Context, food models.Food) (*models.Food, error) {
	query := `INSERT INTO foods (name, calories, protein, fats, carbs)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`
	err := f.db.QueryRow(ctx, query,
		food.Name,
		food.Calories,
		food.Protein,
		food.Fat,
		food.Carbs,
	).Scan(&food.ID)
	if err != nil {
		return nil, err
	}
	return &food, nil
}
func (f *foodRepository) DeleteFood(ctx context.Context, id int64) error {
	query := `DELETE FROM foods 
	WHERE id = $1
	`
	_, err := f.db.Exec(ctx, query, id)
	return err

}

func (f *foodRepository) GetFoodByID(ctx context.Context, id int64) (*models.Food, error) {
	query := `SELECT id, name, calories, protein, fats, carbs
	FROM foods
	WHERE id = $1
	`
	var food models.Food
	err := f.db.QueryRow(ctx, query, id).Scan(&food.ID,
		&food.Name,
		&food.Calories,
		&food.Protein,
		&food.Fat,
		&food.Carbs)
	if err != nil {
		return nil, err
	}
	return &food, nil
}

func (f *foodRepository) GetAllFoods(ctx context.Context) ([]models.Food, error) {
	query := `SELECT id, name, calories, protein, fats, carbs
	FROM foods
	`
	var foods []models.Food
	rows, err := f.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var food models.Food
		err := rows.Scan(
			&food.ID,
			&food.Name,
			&food.Calories,
			&food.Protein,
			&food.Fat,
			&food.Carbs)
		if err != nil {
			return nil, err
		}
		foods = append(foods, food)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return foods, nil
}

func (f *foodRepository) UpdateFood(ctx context.Context, id int64, food models.Food) (*models.Food, error) {
	query := `UPDATE foods 
	SET name = $1, calories = $2, protein = $3, fats = $4, carbs = $5
	WHERE id = $6
	RETURNING id, name, calories, protein, fats, carbs
	`
	err := f.db.QueryRow(ctx, query,
		food.Name,
		food.Calories,
		food.Protein,
		food.Fat,
		food.Carbs,
		id).Scan(&food.ID, &food.Name, &food.Calories, &food.Protein, &food.Fat, &food.Carbs)
	if err != nil {
		return nil, err
	}
	return &food, nil
}
