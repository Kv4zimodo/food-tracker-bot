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

func TestCreateMeal(t *testing.T) {
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

	meal := models.Meal{
		UserID:   1,
		Category: "breakfast",
	}

	repo := NewMealRepository(conn)

	createdMeal, err := repo.CreateMeal(ctx, meal)
	if err != nil {
		t.Fatal(err)
	}

	if createdMeal.UserID != meal.UserID {
		t.Errorf("expected %d, got %d", meal.UserID, createdMeal.UserID)
	}

	if createdMeal.Category != meal.Category {
		t.Errorf("expected %s, got %s", meal.Category, createdMeal.Category)
	}

	t.Cleanup(func() {
		_ = repo.DeleteMeal(ctx, createdMeal.ID)
	})
}

func TestDeleteMeal(t *testing.T) {
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

	meal := models.Meal{
		UserID:   1,
		Category: "breakfast",
	}

	repo := NewMealRepository(conn)

	createdMeal, err := repo.CreateMeal(ctx, meal)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.DeleteMeal(ctx, createdMeal.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.GetMealByID(ctx, createdMeal.ID)
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("expected pgx.ErrNoRows, got %v", err)
	}

	t.Cleanup(func() {
		_ = repo.DeleteMeal(ctx, createdMeal.ID)
	})
}

func TestGetMealByID(t *testing.T) {
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

	meal := models.Meal{
		UserID:   1,
		Category: "breakfast",
	}

	repo := NewMealRepository(conn)

	createdMeal, err := repo.CreateMeal(ctx, meal)
	if err != nil {
		t.Fatal(err)
	}

	foundMeal, err := repo.GetMealByID(ctx, createdMeal.ID)
	if err != nil {
		t.Fatal(err)
	}

	if foundMeal.ID != createdMeal.ID {
		t.Errorf("expected %d, got %d", createdMeal.ID, foundMeal.ID)
	}

	if foundMeal.UserID != createdMeal.UserID {
		t.Errorf("expected %d, got %d", createdMeal.UserID, foundMeal.UserID)
	}

	if foundMeal.Category != createdMeal.Category {
		t.Errorf("expected %s, got %s", createdMeal.Category, foundMeal.Category)
	}

	t.Cleanup(func() {
		_ = repo.DeleteMeal(ctx, createdMeal.ID)
	})
}

func TestFoundAllMeals(t *testing.T) {
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

	meal := models.Meal{
		UserID:   1,
		Category: "breakfast",
	}

	repo := NewMealRepository(conn)

	createdMeal, err := repo.CreateMeal(ctx, meal)
	if err != nil {
		t.Fatal(err)
	}

	createdMeals, err := repo.GetAllMeals(ctx)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, testMeals := range createdMeals {
		if testMeals.ID == createdMeal.ID {
			found = true
			break
		}
	}

	if !found {
		t.Error("created meal was not found")
	}

	t.Cleanup(func() {
		_ = repo.DeleteMeal(ctx, createdMeal.ID)
	})
}

func TestUpdateMeal(t *testing.T) {
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

	meal1 := models.Meal{
		UserID:   1,
		Category: "breakfast",
	}

	meal2 := models.Meal{
		UserID:   1,
		Category: "lunch",
	}

	repo := NewMealRepository(conn)

	createdMeal, err := repo.CreateMeal(ctx, meal1)
	if err != nil {
		t.Fatal(err)
	}

	updatedMeal, err := repo.UpdateMeal(ctx, createdMeal.ID, meal2)
	if err != nil {
		t.Fatal(err)
	}

	if updatedMeal.UserID != meal2.UserID {
		t.Errorf("expected %d, got %d", meal2.UserID, updatedMeal.UserID)
	}

	if updatedMeal.Category != meal2.Category {
		t.Errorf("expected %s, got %s", meal2.Category, updatedMeal.Category)
	}

	t.Cleanup(func() {
		_ = repo.DeleteMeal(ctx, createdMeal.ID)
	})
}
