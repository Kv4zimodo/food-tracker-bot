package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kv4zimodo/food-tracker-bot/internal/models"
)

type userRepository struct {
	db *pgx.Conn
}

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id int64) error
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	UpdateUser(ctx context.Context, id int64, user models.User) (*models.User, error)
}

func NewUserRepository(db *pgx.Conn) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	query := `INSERT INTO users (user_name)
	VALUES ($1)
	RETURNING id
	`
	err := u.db.QueryRow(ctx, query,
		user.Name).Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users
	WHERE id = $1
	`
	_, err := u.db.Exec(ctx, query, id)
	return err
}

func (u *userRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	query := `SELECT id, user_name
	FROM users
	WHERE id = $1
	`
	var user models.User
	err := u.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, user_name
	FROM users
	`
	var users []models.User
	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userRepository) UpdateUser(ctx context.Context, id int64, user models.User) (*models.User, error) {
	query := `UPDATE users
	SET user_name = $1
	WHERE id = $2
	RETURNING id, user_name
	`
	err := u.db.QueryRow(ctx, query,
		user.Name,
		id).Scan(
		&user.ID,
		&user.Name)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
