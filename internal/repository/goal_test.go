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

func TestCreateGoal(t *testing.T) {
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
	goal := models.Goal{
		UserID:   1,
		Calories: 2200,
		Protein:  90,
		Fat:      60,
		Carbs:    250,
	}
	repo := NewGoalRepository(conn)
	createdGoal, err := repo.CreateGoal(ctx, goal)
	if err != nil {
		t.Fatal(err)
	}
	if createdGoal.UserID != goal.UserID {
		t.Errorf("expected %d, got %d", goal.UserID, createdGoal.UserID)
	}
	if createdGoal.Calories != goal.Calories {
		t.Errorf("expected %f, got %f", goal.Calories, createdGoal.Calories)
	}
	if createdGoal.Protein != goal.Protein {
		t.Errorf("expected %f, got %f", goal.Protein, createdGoal.Protein)
	}
	if createdGoal.Fat != goal.Fat {
		t.Errorf("expected %f, got %f", goal.Fat, createdGoal.Fat)
	}
	if createdGoal.Carbs != goal.Carbs {
		t.Errorf("expected %f, got %f", goal.Carbs, createdGoal.Carbs)
	}
	t.Cleanup(func() {
		_ = repo.DeleteGoal(ctx, createdGoal.UserID)
	})
}

func TestDeleteGoal(t *testing.T) {
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
	goal := models.Goal{
		UserID:   1,
		Calories: 2200,
		Protein:  90,
		Fat:      60,
		Carbs:    250,
	}
	repo := NewGoalRepository(conn)
	createdGoal, err := repo.CreateGoal(ctx, goal)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.DeleteGoal(ctx, createdGoal.UserID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.GetGoalByUserID(ctx, createdGoal.UserID)
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("expected pgx.ErrNoRows, got %v", err)
	}
	t.Cleanup(func() {
		_ = repo.DeleteGoal(ctx, createdGoal.UserID)
	})
}

func TestGetGoalByUserID(t *testing.T) {
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
	goal := models.Goal{
		UserID:   1,
		Calories: 2200,
		Protein:  90,
		Fat:      60,
		Carbs:    250,
	}
	repo := NewGoalRepository(conn)
	createdGoal, err := repo.CreateGoal(ctx, goal)
	if err != nil {
		t.Fatal(err)
	}
	foundGoal, err := repo.GetGoalByUserID(ctx, createdGoal.UserID)
	if err != nil {
		t.Fatal(err)
	}
	if foundGoal.Calories != createdGoal.Calories {
		t.Errorf("expected %f, got %f", createdGoal.Calories, foundGoal.Calories)
	}
	if foundGoal.Protein != createdGoal.Protein {
		t.Errorf("expected %f, got %f", createdGoal.Protein, foundGoal.Protein)
	}
	if foundGoal.Fat != createdGoal.Fat {
		t.Errorf("expected %f, got %f", createdGoal.Fat, foundGoal.Fat)
	}
	if foundGoal.Carbs != createdGoal.Carbs {
		t.Errorf("expected %f, got %f", createdGoal.Carbs, foundGoal.Carbs)
	}
	t.Cleanup(func() {
		_ = repo.DeleteGoal(ctx, createdGoal.UserID)
	})
}

func TestUpdateGoal(t *testing.T) {
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
	goal1 := models.Goal{
		UserID:   1,
		Calories: 2200,
		Protein:  90,
		Fat:      60,
		Carbs:    250,
	}
	goal2 := models.Goal{
		UserID:   1,
		Calories: 2500,
		Protein:  10,
		Fat:      80,
		Carbs:    350,
	}
	repo := NewGoalRepository(conn)
	createdGoal, err := repo.CreateGoal(ctx, goal1)
	if err != nil {
		t.Fatal(err)
	}
	updatedGoal, err := repo.UpdateGoal(ctx, createdGoal.UserID, goal2)
	if err != nil {
		t.Fatal(err)
	}
	if updatedGoal.Calories != goal2.Calories {
		t.Errorf("expected %f, got %f", goal2.Calories, updatedGoal.Calories)
	}
	if updatedGoal.Protein != goal2.Protein {
		t.Errorf("expected %f, got %f", goal2.Protein, updatedGoal.Protein)
	}
	if updatedGoal.Fat != goal2.Fat {
		t.Errorf("expected %f, got %f", goal2.Fat, updatedGoal.Fat)
	}
	if updatedGoal.Carbs != goal2.Carbs {
		t.Errorf("expected %f, got %f", goal2.Carbs, updatedGoal.Carbs)
	}
	t.Cleanup(func() {
		_ = repo.DeleteGoal(ctx, createdGoal.UserID)
	})
}
