package handler

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kv4zimodo/food-tracker-bot/internal/models"
	"github.com/kv4zimodo/food-tracker-bot/internal/service"
)

func (h *handler) Today(ctx context.Context, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	telegramID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	user, err := h.serviceUser.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		log.Println("Ошибка получения пользователя:", err)
		h.bot.Send(tgbotapi.NewMessage(chatID, "Не удалось получить пользователя"))
		return
	}

	now := time.Now()

	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)

	meals, err := h.serviceMeal.GetMealsByUserAndDate(
		ctx,
		user.ID,
		today,
	)
	if err != nil {
		log.Println("Ошибка получения приемов пищи:", err)
		h.bot.Send(tgbotapi.NewMessage(chatID, "Не удалось получить данные за сегодня"))
		return
	}

	if len(meals) == 0 {
		h.bot.Send(tgbotapi.NewMessage(
			chatID,
			"Сегодня пока ничего не добавлено",
		))
		return
	}

	var result strings.Builder

	result.WriteString("Сегодня\n\n")

	var total service.Nutrition

	for _, meal := range meals {
		mealItems, err := h.serviceMealItem.GetMealItemsByMealID(
			ctx,
			meal.ID,
		)
		if err != nil {
			log.Println("Ошибка получения продуктов:", err)
			continue
		}

		result.WriteString(mealCategoryName(meal.Category))
		result.WriteString("\n")

		for _, item := range mealItems {
			food, err := h.serviceFood.GetFoodByID(
				ctx,
				item.FoodID,
			)
			if err != nil {
				log.Println("Ошибка получения продукта:", err)
				continue
			}

			nutrition := service.CalculateNutrition(*food, item.Weight)

			total.Calories += nutrition.Calories
			total.Protein += nutrition.Protein
			total.Fat += nutrition.Fat
			total.Carbs += nutrition.Carbs

			result.WriteString(fmt.Sprintf("%s - %.0f г\n",
				food.Name,
				item.Weight))
		}

		result.WriteString("\n")
	}

	result.WriteString("Итого:\n")
	result.WriteString(fmt.Sprintf("Калории: %.1f ккал\n"+
		"Белки: %.1f г\n"+
		"Жиры: %.1f г\n"+
		"Углеводы: %.1f г",
		total.Calories,
		total.Protein,
		total.Fat,
		total.Carbs))

	h.bot.Send(tgbotapi.NewMessage(chatID, result.String()))
}

func mealCategoryName(category models.MealCategory) string {
	switch category {
	case models.Breakfast:
		return "Завтрак"
	case models.Lunch:
		return "Обед"
	case models.Dinner:
		return "Ужин"
	case models.Snack:
		return "Перекус"
	default:
		return "Прием пищи"
	}
}
