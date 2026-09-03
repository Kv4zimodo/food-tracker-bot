package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context) (*pgx.Conn, error) {
	connString := os.Getenv("CONN_STRING")
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}
	return conn, nil

}
