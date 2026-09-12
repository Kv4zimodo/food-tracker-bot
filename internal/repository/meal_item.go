package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kv4zimodo/food-tracker-bot/internal/models"
)

type mealItemRepository struct {
	db *pgx.Conn
}

type MealItemRepository interface {
	CreateMealItem(ctx context.Context, mealItem models.MealItem) (*models.MealItem, error)
	DeleteMealItem(ctx context.Context, id int64) error
	GetMealItemByID(ctx context.Context, id int64) (*models.MealItem, error)
	GetAllMealItems(ctx context.Context) ([]models.MealItem, error)
	UpdateMealItem(ctx context.Context, id int64, mealItem models.MealItem) (*models.MealItem, error)
}

func NewMealItemRepository(db *pgx.Conn) MealItemRepository {
	return &mealItemRepository{
		db: db,
	}
}

func (m *mealItemRepository) CreateMealItem(ctx context.Context, mealItem models.MealItem) (*models.MealItem, error) {
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

func (m *mealItemRepository) DeleteMealItem(ctx context.Context, id int64) error {
	query := `DELETE FROM meal_items
	WHERE id = $1
	`
	_, err := m.db.Exec(ctx, query, id)
	return err
}

func (m *mealItemRepository) GetMealItemByID(ctx context.Context, id int64) (*models.MealItem, error) {
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

func (m *mealItemRepository) GetAllMealItems(ctx context.Context) ([]models.MealItem, error) {
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

func (m *mealItemRepository) UpdateMealItem(ctx context.Context, id int64, mealItem models.MealItem) (*models.MealItem, error) {
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
