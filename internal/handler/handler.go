package handler

import (
	"context"
	"strings"

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
		if strings.HasPrefix(update.CallbackQuery.Data, "delete_meal_item_") {
			h.DeleteMealItemCallback(ctx, update)
			return
		}

		switch update.CallbackQuery.Data {
		case "goal_mode_calories", "goal_mode_bju":
			h.GoalModeCallback(update)
		default:
			h.MealCallBack(update)
		}
		return
	}

	if update.Message == nil {
		return
	}

	switch update.Message.Text {
	case "/start":
		delete(states, update.Message.From.ID)
		h.HandlerStart(ctx, update)
	case "/help":
		h.HandlerHelp(update)
	case "Добавить еду":
		h.AddFood(update)
	case "Сегодня":
		h.Today(ctx, update)
	case "Моя цель":
		h.Goal(ctx, update)
	case "Задать цель":
		h.SetGoal(update)
	default:
		state := states[update.Message.From.ID]

		if state.State == "waiting_food" {
			h.FoodInput(ctx, update)
			return
		}

		if state.State == "waiting_weight" {
			h.WeightInput(ctx, update)
			return
		}
		if state.State == "waiting_goal_calories" ||
			state.State == "waiting_goal_protein" ||
			state.State == "waiting_goal_fat" ||
			state.State == "waiting_goal_carbs" {
			h.GoalInput(ctx, update)
			return
		}
	}
}
