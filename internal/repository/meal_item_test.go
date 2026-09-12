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

func TestCreateMealItem(t *testing.T) {
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

	user := models.User{
		Name: "test_user",
	}

	userRepo := NewUserRepository(conn)

	createdUser, err := userRepo.CreateUser(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	meal := models.Meal{
		UserID:   createdUser.ID,
		Category: "breakfast",
	}

	mealRepo := NewMealRepository(conn)

	createdMeal, err := mealRepo.CreateMeal(ctx, meal)
	if err != nil {
		t.Fatal(err)
	}

	food := models.Food{
		Name:     "Chicken",
		Calories: 165,
		Protein:  31,
		Fat:      3.6,
		Carbs:    0,
	}

	foodRepo := NewFoodRepository(conn)

	createdFood, err := foodRepo.CreateFood(ctx, food)
	if err != nil {
		t.Fatal(err)
	}

	mealItem := models.MealItem{
		MealID: createdMeal.ID,
		FoodID: createdFood.ID,
		Weight: 150,
	}

	repo := NewMealItemRepository(conn)

	createdMealItem, err := repo.CreateMealItem(ctx, mealItem)
	if err != nil {
		t.Fatal(err)
	}

	if createdMealItem.MealID != mealItem.MealID {
		t.Errorf("expected %d, got %d", mealItem.MealID, createdMealItem.MealID)
	}

	if createdMealItem.FoodID != mealItem.FoodID {
		t.Errorf("expected %d, got %d", mealItem.FoodID, createdMealItem.FoodID)
	}

	if createdMealItem.Weight != mealItem.Weight {
		t.Errorf("expected %f, got %f", mealItem.Weight, createdMealItem.Weight)
	}

	t.Cleanup(func() {
		_ = repo.DeleteMealItem(ctx, createdMealItem.ID)
		_ = foodRepo.DeleteFood(ctx, createdFood.ID)
		_ = mealRepo.DeleteMeal(ctx, createdMeal.ID)
		_ = userRepo.DeleteUser(ctx, createdUser.ID)
	})
}

func TestDeleteMealItem(t *testing.T) {
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

	user := models.User{
		Name: "test_user",
	}

	userRepo := NewUserRepository(conn)

	createdUser, err := userRepo.CreateUser(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	meal := models.Meal{
		UserID:   createdUser.ID,
		Category: "breakfast",
	}

	mealRepo := NewMealRepository(conn)

	createdMeal, err := mealRepo.CreateMeal(ctx, meal)
	if err != nil {
		t.Fatal(err)
	}

	food := models.Food{
		Name:     "Chicken",
		Calories: 165,
		Protein:  31,
		Fat:      3.6,
		Carbs:    0,
	}

	foodRepo := NewFoodRepository(conn)

	createdFood, err := foodRepo.CreateFood(ctx, food)
	if err != nil {
		t.Fatal(err)
	}

	mealItem := models.MealItem{
		MealID: createdMeal.ID,
		FoodID: createdFood.ID,
		Weight: 150,
	}

	repo := NewMealItemRepository(conn)

	createdMealItem, err := repo.CreateMealItem(ctx, mealItem)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.DeleteMealItem(ctx, createdMealItem.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.GetMealItemByID(ctx, createdMealItem.ID)
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("expected pgx.ErrNoRows, got %v", err)
	}

	t.Cleanup(func() {
		_ = repo.DeleteMealItem(ctx, createdMealItem.ID)
		_ = foodRepo.DeleteFood(ctx, createdFood.ID)
		_ = mealRepo.DeleteMeal(ctx, createdMeal.ID)
		_ = userRepo.DeleteUser(ctx, createdUser.ID)
	})
}

func TestGetMealItemByID(t *testing.T) {
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

	user := models.User{
		Name: "test_user",
	}

	userRepo := NewUserRepository(conn)

	createdUser, err := userRepo.CreateUser(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	meal := models.Meal{
		UserID:   createdUser.ID,
		Category: "breakfast",
	}

	mealRepo := NewMealRepository(conn)

	createdMeal, err := mealRepo.CreateMeal(ctx, meal)
	if err != nil {
		t.Fatal(err)
	}

	food := models.Food{
		Name:     "Chicken",
		Calories: 165,
		Protein:  31,
		Fat:      3.6,
		Carbs:    0,
	}

	foodRepo := NewFoodRepository(conn)

	createdFood, err := foodRepo.CreateFood(ctx, food)
	if err != nil {
		t.Fatal(err)
	}

	mealItem := models.MealItem{
		MealID: createdMeal.ID,
		FoodID: createdFood.ID,
		Weight: 150,
	}

	repo := NewMealItemRepository(conn)

	createdMealItem, err := repo.CreateMealItem(ctx, mealItem)
	if err != nil {
		t.Fatal(err)
	}

	foundMealItem, err := repo.GetMealItemByID(ctx, createdMealItem.ID)
	if err != nil {
		t.Fatal(err)
	}

	if foundMealItem.ID != createdMealItem.ID {
		t.Errorf("expected %d, got %d", createdMealItem.ID, foundMealItem.ID)
	}

	if foundMealItem.MealID != createdMealItem.MealID {
		t.Errorf("expected %d, got %d", createdMealItem.MealID, foundMealItem.MealID)
	}

	if foundMealItem.FoodID != createdMealItem.FoodID {
		t.Errorf("expected %d, got %d", createdMealItem.FoodID, foundMealItem.FoodID)
	}

	if foundMealItem.Weight != createdMealItem.Weight {
		t.Errorf("expected %f, got %f", createdMealItem.Weight, foundMealItem.Weight)
	}

	t.Cleanup(func() {
		_ = repo.DeleteMealItem(ctx, createdMealItem.ID)
		_ = foodRepo.DeleteFood(ctx, createdFood.ID)
		_ = mealRepo.DeleteMeal(ctx, createdMeal.ID)
		_ = userRepo.DeleteUser(ctx, createdUser.ID)
	})
}

func TestFoundAllMealItems(t *testing.T) {
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

	user := models.User{
		Name: "test_user",
	}

	userRepo := NewUserRepository(conn)

	createdUser, err := userRepo.CreateUser(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	meal := models.Meal{
		UserID:   createdUser.ID,
		Category: "breakfast",
	}

	mealRepo := NewMealRepository(conn)

	createdMeal, err := mealRepo.CreateMeal(ctx, meal)
	if err != nil {
		t.Fatal(err)
	}

	food := models.Food{
		Name:     "Chicken",
		Calories: 165,
		Protein:  31,
		Fat:      3.6,
		Carbs:    0,
	}

	foodRepo := NewFoodRepository(conn)

	createdFood, err := foodRepo.CreateFood(ctx, food)
	if err != nil {
		t.Fatal(err)
	}

	mealItem := models.MealItem{
		MealID: createdMeal.ID,
		FoodID: createdFood.ID,
		Weight: 150,
	}

	repo := NewMealItemRepository(conn)

	createdMealItem, err := repo.CreateMealItem(ctx, mealItem)
	if err != nil {
		t.Fatal(err)
	}

	createdMealItems, err := repo.GetAllMealItems(ctx)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, testMealItems := range createdMealItems {
		if testMealItems.ID == createdMealItem.ID {
			found = true
			break
		}
	}

	if !found {
		t.Error("created meal item was not found")
	}

	t.Cleanup(func() {
		_ = repo.DeleteMealItem(ctx, createdMealItem.ID)
		_ = foodRepo.DeleteFood(ctx, createdFood.ID)
		_ = mealRepo.DeleteMeal(ctx, createdMeal.ID)
		_ = userRepo.DeleteUser(ctx, createdUser.ID)
	})
}

func TestUpdateMealItem(t *testing.T) {
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

	user := models.User{
		Name: "test_user",
	}

	userRepo := NewUserRepository(conn)

	createdUser, err := userRepo.CreateUser(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	meal := models.Meal{
		UserID:   createdUser.ID,
		Category: "breakfast",
	}

	mealRepo := NewMealRepository(conn)

	createdMeal, err := mealRepo.CreateMeal(ctx, meal)
	if err != nil {
		t.Fatal(err)
	}

	food1 := models.Food{
		Name:     "Chicken",
		Calories: 165,
		Protein:  31,
		Fat:      3.6,
		Carbs:    0,
	}

	food2 := models.Food{
		Name:     "Rice",
		Calories: 130,
		Protein:  2.7,
		Fat:      0.3,
		Carbs:    28,
	}

	foodRepo := NewFoodRepository(conn)

	createdFood1, err := foodRepo.CreateFood(ctx, food1)
	if err != nil {
		t.Fatal(err)
	}

	createdFood2, err := foodRepo.CreateFood(ctx, food2)
	if err != nil {
		t.Fatal(err)
	}

	mealItem1 := models.MealItem{
		MealID: createdMeal.ID,
		FoodID: createdFood1.ID,
		Weight: 150,
	}

	mealItem2 := models.MealItem{
		MealID: createdMeal.ID,
		FoodID: createdFood2.ID,
		Weight: 200,
	}

	repo := NewMealItemRepository(conn)

	createdMealItem, err := repo.CreateMealItem(ctx, mealItem1)
	if err != nil {
		t.Fatal(err)
	}

	updatedMealItem, err := repo.UpdateMealItem(ctx, createdMealItem.ID, mealItem2)
	if err != nil {
		t.Fatal(err)
	}

	if updatedMealItem.MealID != mealItem2.MealID {
		t.Errorf("expected %d, got %d", mealItem2.MealID, updatedMealItem.MealID)
	}

	if updatedMealItem.FoodID != mealItem2.FoodID {
		t.Errorf("expected %d, got %d", mealItem2.FoodID, updatedMealItem.FoodID)
	}

	if updatedMealItem.Weight != mealItem2.Weight {
		t.Errorf("expected %f, got %f", mealItem2.Weight, updatedMealItem.Weight)
	}

	t.Cleanup(func() {
		_ = repo.DeleteMealItem(ctx, createdMealItem.ID)
		_ = foodRepo.DeleteFood(ctx, createdFood1.ID)
		_ = foodRepo.DeleteFood(ctx, createdFood2.ID)
		_ = mealRepo.DeleteMeal(ctx, createdMeal.ID)
		_ = userRepo.DeleteUser(ctx, createdUser.ID)
	})
}
