package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/kv4zimodo/food-tracker-bot/internal/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Не удалось загрузить конфигурацию: ", err)
	}

	ctx := context.Background()
	conn, err := database.Connect(ctx)
	if err != nil {
		log.Fatal("Connection error", err)
	}
	defer conn.Close(ctx)
	log.Print("Успешно подключились!")
}
