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

func TestCreateUser(t *testing.T) {
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

	repo := NewUserRepository(conn)

	createdUser, err := repo.CreateUser(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	if createdUser.Name != user.Name {
		t.Errorf("expected %s, got %s", user.Name, createdUser.Name)
	}

	t.Cleanup(func() {
		_ = repo.DeleteUser(ctx, createdUser.ID)
	})
}

func TestDeleteUser(t *testing.T) {
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

	repo := NewUserRepository(conn)

	createdUser, err := repo.CreateUser(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.DeleteUser(ctx, createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.GetUserByID(ctx, createdUser.ID)
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("expected pgx.ErrNoRows, got %v", err)
	}

	t.Cleanup(func() {
		_ = repo.DeleteUser(ctx, createdUser.ID)
	})
}

func TestGetUserByID(t *testing.T) {
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

	repo := NewUserRepository(conn)

	createdUser, err := repo.CreateUser(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	foundUser, err := repo.GetUserByID(ctx, createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}

	if foundUser.ID != createdUser.ID {
		t.Errorf("expected %d, got %d", createdUser.ID, foundUser.ID)
	}

	if foundUser.Name != createdUser.Name {
		t.Errorf("expected %s, got %s", createdUser.Name, foundUser.Name)
	}

	t.Cleanup(func() {
		_ = repo.DeleteUser(ctx, createdUser.ID)
	})
}

func TestFoundAllUsers(t *testing.T) {
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

	repo := NewUserRepository(conn)

	createdUser, err := repo.CreateUser(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	createdUsers, err := repo.GetAllUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, testUsers := range createdUsers {
		if testUsers.ID == createdUser.ID {
			found = true
			break
		}
	}

	if !found {
		t.Error("created user was not found")
	}

	t.Cleanup(func() {
		_ = repo.DeleteUser(ctx, createdUser.ID)
	})
}

func TestUpdateUser(t *testing.T) {
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

	user1 := models.User{
		Name: "test_user",
	}

	user2 := models.User{
		Name: "updated_user",
	}

	repo := NewUserRepository(conn)

	createdUser, err := repo.CreateUser(ctx, user1)
	if err != nil {
		t.Fatal(err)
	}

	updatedUser, err := repo.UpdateUser(ctx, createdUser.ID, user2)
	if err != nil {
		t.Fatal(err)
	}

	if updatedUser.Name != user2.Name {
		t.Errorf("expected %s, got %s", user2.Name, updatedUser.Name)
	}

	t.Cleanup(func() {
		_ = repo.DeleteUser(ctx, createdUser.ID)
	})
}
