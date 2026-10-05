package handler

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	if update.CallbackQuery != nil {
		h.MealCallBack(update)
		return
	}

	if update.Message == nil {
		return
	}

	switch update.Message.Text {
	case "/start":
		h.HandlerStart(ctx, update)
	case "/help":
		h.HandlerHelp(update)
	case "Добавить еду":
		h.AddFood(update)
	case "Сегодня":
		h.Today(ctx, update)
	default:
		state := states[update.Message.From.ID]

		if state.State == "waiting_food" {
			h.FoodInput(update)
			return
		}

		if state.State == "waiting_weight" {
			h.WeightInput(ctx, update)
			return
		}
	}
}
