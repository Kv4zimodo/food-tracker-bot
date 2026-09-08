package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kv4zimodo/food-tracker-bot/internal/models"
)

type MealItemRepository struct {
	db *pgx.Conn
}

func NewMealItemRepository(db *pgx.Conn) *MealItemRepository {
	return &MealItemRepository{
		db: db,
	}
}

func (m *MealItemRepository) CreateMealItem(ctx context.Context, mealItem models.MealItem) (*models.MealItem, error) {
	query := `INSERT INTO meal_items (meal_id, food_id, weight)
	VALUES ($1, $2, $3)
	RETURNING id
	`
	err := m.db.QueryRow(ctx, query,
		mealItem.MealID,
		mealItem.FoodID,
		mealItem.Weight).Scan(&mealItem.ID)
	if err != nil {
		return nil, err
	}
	return &mealItem, nil
}

func (m *MealItemRepository) DeleteMealItem(ctx context.Context, id int64) error {
	query := `DELETE FROM meal_items
	WHERE id = $1
	`
	_, err := m.db.Exec(ctx, query, id)
	return err
}

func (m *MealItemRepository) GetMealItemByID(ctx context.Context, id int64) (*models.MealItem, error) {
	query := `SELECT id, meal_id, food_id, weight
	FROM meal_items
	WHERE id = $1
	`
	var mealItem models.MealItem
	err := m.db.QueryRow(ctx, query, id).Scan(
		&mealItem.ID,
		&mealItem.MealID,
		&mealItem.FoodID,
		&mealItem.Weight)
	if err != nil {
		return nil, err
	}
	return &mealItem, nil
}

func (m *MealItemRepository) GetAllMealItems(ctx context.Context) ([]models.MealItem, error) {
	query := `SELECT id, meal_id, food_id, weight
	FROM meal_items
	`
	var mealItems []models.MealItem
	rows, err := m.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var mealItem models.MealItem
		err := rows.Scan(
			&mealItem.ID,
			&mealItem.MealID,
			&mealItem.FoodID,
			&mealItem.Weight)
		if err != nil {
			return nil, err
		}
		mealItems = append(mealItems, mealItem)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return mealItems, nil
}

func (m *MealItemRepository) UpdateMealItem(ctx context.Context, id int64, mealItem models.MealItem) (*models.MealItem, error) {
	query := `UPDATE meal_items
	SET meal_id = $1, food_id = $2, weight = $3
	WHERE id = $4
	RETURNING id, meal_id, food_id, weight
	`
	err := m.db.QueryRow(ctx, query,
		mealItem.MealID,
		mealItem.FoodID,
		mealItem.Weight,
		id).Scan(
		&mealItem.ID,
		&mealItem.MealID,
		&mealItem.FoodID,
		&mealItem.Weight)
	if err != nil {
		return nil, err
	}
	return &mealItem, nil
}
