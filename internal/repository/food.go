package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kv4zimodo/food-tracker-bot/internal/models"
)

type FoodRepository struct {
	db *pgx.Conn
}

func NewFoodRepository(db *pgx.Conn) *FoodRepository {
	return &FoodRepository{
		db: db,
	}
}

func (f *FoodRepository) CreateFood(food models.Food, ctx context.Context) (*models.Food, error) {
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
func (f *FoodRepository) DeleteFood(id int64, ctx context.Context) error {
	query := `DELETE FROM foods 
	WHERE id = $1
	`
	_, err := f.db.Exec(ctx, query, id)
	return err

}

func (f *FoodRepository) GetFoodByID(id int64, ctx context.Context) (*models.Food, error) {
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

func (f *FoodRepository) GetAllFoods(ctx context.Context) ([]models.Food, error) {
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

func (f *FoodRepository) UpdateFood(ctx context.Context, id int64, food models.Food) (*models.Food, error) {
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
