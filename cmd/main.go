package main

import (
	"context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/kv4zimodo/food-tracker-bot/internal/database"
	"github.com/kv4zimodo/food-tracker-bot/internal/handler"
	"github.com/kv4zimodo/food-tracker-bot/internal/repository"
	"github.com/kv4zimodo/food-tracker-bot/internal/service"
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

	repoFood := repository.NewFoodRepository(conn)
	serviceFood := service.NewFoodService(repoFood)

	repoUser := repository.NewUserRepository(conn)
	serviceUser := service.NewUserService(repoUser)

	repoGoal := repository.NewGoalRepository(conn)
	serviceGoal := service.NewGoalService(repoGoal)

	repoMeal := repository.NewMealRepository(conn)
	serviceMeal := service.NewMealService(repoMeal)

	repoMealItem := repository.NewMealItemRepository(conn)
	serviceMealItem := service.NewMealItemService(repoMealItem)

	token := os.Getenv("BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Ошибка Бота: ", err)
	}

	h := handler.NewHandler(bot,
		serviceFood, serviceGoal,
		serviceMeal,
		serviceMealItem,
		serviceUser)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		h.HandleUpdate(ctx, update)
	}
}

//  ЗАКОММИТЬ
