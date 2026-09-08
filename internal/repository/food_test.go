package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/kv4zimodo/food-tracker-bot/internal/database"
	"github.com/kv4zimodo/food-tracker-bot/internal/models"
)

func TestCreatedFood(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	conn, err := database.Connect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		conn.Close(ctx)
	})

	food := models.Food{
		Name:     "sausage",
		Calories: 210,
		Protein:  12,
		Fat:      15,
		Carbs:    25,
	}
	repo := NewFoodRepository(conn)
	createdFood, err := repo.CreateFood(food, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if createdFood.ID == 0 {
		t.Error("food ID must be than 0")
	}
	if createdFood.Name != food.Name {
		t.Errorf("expected %s, got %s", food.Name, createdFood.Name)
	}

	if createdFood.Calories != food.Calories {
		t.Errorf("expected %f, got %f", food.Calories, createdFood.Calories)
	}

	if createdFood.Protein != food.Protein {
		t.Errorf("expected %f, got %f", food.Protein, createdFood.Protein)
	}

	if createdFood.Fat != food.Fat {
		t.Errorf("expected %f, got %f", food.Fat, createdFood.Fat)
	}

	if createdFood.Carbs != food.Carbs {
		t.Errorf("expected %f, got %f", food.Carbs, createdFood.Carbs)
	}
	t.Cleanup(func() {
		_ = repo.DeleteFood(createdFood.ID, ctx)
	})
}

func TestDeletedFood(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	conn, err := database.Connect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		conn.Close(ctx)
	})

	food := models.Food{
		Name:     "sausage",
		Calories: 210,
		Protein:  12,
		Fat:      15,
		Carbs:    25,
	}
	repo := NewFoodRepository(conn)
	createdFood, err := repo.CreateFood(food, ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.DeleteFood(createdFood.ID, ctx)
	if err != nil {
		t.Error(err)
	}
	_, err = repo.GetFoodByID(createdFood.ID, ctx)
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("expected pgx.ErrNoRows, got %v", err)
	}

}

func TestGetFoodByID(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	conn, err := database.Connect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		conn.Close(ctx)
	})

	food := models.Food{
		Name:     "sausage",
		Calories: 210,
		Protein:  12,
		Fat:      15,
		Carbs:    25,
	}
	repo := NewFoodRepository(conn)
	createdFood, err := repo.CreateFood(food, ctx)
	if err != nil {
		t.Fatal(err)
	}
	foundFood, err := repo.GetFoodByID(createdFood.ID, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if foundFood.Name != food.Name {
		t.Errorf("expected %s, got %s", food.Name, foundFood.Name)
	}
	if foundFood.Calories != food.Calories {
		t.Errorf("expected %f, got %f", food.Calories, foundFood.Calories)
	}
	if foundFood.Protein != food.Protein {
		t.Errorf("expected %f, got %f", food.Protein, foundFood.Protein)
	}
	if foundFood.Fat != food.Fat {
		t.Errorf("expected %f, got %f", food.Fat, foundFood.Fat)
	}
	if foundFood.Carbs != food.Carbs {
		t.Errorf("expected %f, got %f", food.Carbs, foundFood.Carbs)
	}
}

func TestFoundAllFoods(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	conn, err := database.Connect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		conn.Close(ctx)
	})

	food := models.Food{
		Name:     "sausage",
		Calories: 210,
		Protein:  12,
		Fat:      15,
		Carbs:    25,
	}

	repo := NewFoodRepository(conn)
	createdFood, err := repo.CreateFood(food, ctx)
	if err != nil {
		t.Fatal(err)
	}

	createdFoods, err := repo.GetAllFoods(ctx)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, testFoods := range createdFoods {
		if testFoods.ID == createdFood.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("created food was not found")
	}
	t.Cleanup(func() {
		_ = repo.DeleteFood(createdFood.ID, ctx)

	})
}

func TestUpdateFood(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	conn, err := database.Connect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		conn.Close(ctx)
	})

	food := models.Food{
		Name:     "sausage",
		Calories: 210,
		Protein:  12,
		Fat:      15,
		Carbs:    25,
	}
	food2 := models.Food{
		Name:     "chicken",
		Calories: 190,
		Protein:  21,
		Fat:      4,
		Carbs:    0,
	}
	repo := NewFoodRepository(conn)
	createdFood, err := repo.CreateFood(food, ctx)
	if err != nil {
		t.Fatal(err)
	}
	updateFood, err := repo.UpdateFood(ctx, createdFood.ID, food2)
	if err != nil {
		t.Fatal(err)
	}
	if updateFood.Name != food2.Name {
		t.Errorf("expected %s, got %s", food2.Name, updateFood.Name)
	}
	if updateFood.Calories != food2.Calories {
		t.Errorf("expected %f, got %f", food2.Calories, updateFood.Calories)
	}
	if updateFood.Protein != food2.Protein {
		t.Errorf("expected %f, got %f", food2.Protein, updateFood.Protein)
	}
	if updateFood.Fat != food2.Fat {
		t.Errorf("expected %f, got %f", food2.Fat, updateFood.Fat)
	}
	if updateFood.Carbs != food2.Carbs {
		t.Errorf("expected %f, got %f", food2.Carbs, updateFood.Carbs)
	}
	t.Cleanup(func() {
		_ = repo.DeleteFood(createdFood.ID, ctx)
	})
}
