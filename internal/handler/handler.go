package handler

import (
	"context"
	"errors"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"github.com/kv4zimodo/food-tracker-bot/internal/models"
	"github.com/kv4zimodo/food-tracker-bot/internal/service"
)

type handler struct {
	bot             *tgbotapi.BotAPI
	serviceFood     service.FoodService
	serviceGoal     service.GoalService
	serviceMeal     service.MealService
	serviceMealItem service.MealItemService
	serviceUser     service.UserService
}

type Handler interface {
}

func NewHandler(bot *tgbotapi.BotAPI,
	serviceFood service.FoodService,
	serviceGoal service.GoalService,
	serviceMeal service.MealService,
	serviceMealItem service.MealItemService,
	serviceUser service.UserService) *handler {
	return &handler{
		bot:             bot,
		serviceFood:     serviceFood,
		serviceGoal:     serviceGoal,
		serviceMeal:     serviceMeal,
		serviceMealItem: serviceMealItem,
		serviceUser:     serviceUser,
	}
}

func (h *handler) HandleUpdate(ctx context.Context, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	switch update.Message.Text {
	case "/start":
		h.HandlerStart(ctx, update)
	}
}

func (h *handler) HandlerStart(ctx context.Context, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	telegramID := update.Message.From.ID
	userName := update.Message.From.FirstName
	chatID := update.Message.Chat.ID

	user, err := h.serviceUser.GetUserByTelegramID(ctx, telegramID)

	if errors.Is(err, pgx.ErrNoRows) {
		user, err = h.serviceUser.CreateUser(ctx, models.User{
			Name:       userName,
			TelegramID: telegramID,
		})
		if err != nil {
			log.Println("Ошибка создания пользователя:", err)
			return
		}
	} else if err != nil {
		log.Println("Ошибка получения пользователя:", err)
		return
	}

	msg := tgbotapi.NewMessage(chatID, "Привет, "+user.Name+"!")

	if _, err := h.bot.Send(msg); err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
}

func (h *handler) HandlerHelp(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	text := `Доступные операции:
	/start
	/help`

	chatID := update.Message.Chat.ID

	msg := tgbotapi.NewMessage(chatID, text)

	if _, err := h.bot.Send(msg); err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
}
