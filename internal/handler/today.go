package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
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

	goal, err := h.serviceGoal.GetGoalByUserID(ctx, user.ID)
	hasGoal := true

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			hasGoal = false
		} else {
			log.Println("Ошибка получения цели:", err)
			h.bot.Send(tgbotapi.NewMessage(
				chatID,
				"Не удалось получить дневную цель",
			))
			return
		}
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
	}

	var result strings.Builder

	result.WriteString("Сегодня\n\n")

	var total service.Nutrition

	var deleteButtons [][]tgbotapi.InlineKeyboardButton

	for _, meal := range meals {
		mealItems, err := h.serviceMealItem.GetMealItemsByMealID(
			ctx,
			meal.ID,
		)
		if err != nil {
			log.Println("Ошибка получения продуктов:", err)
			h.bot.Send(tgbotapi.NewMessage(
				chatID,
				"Не удалось рассчитать КБЖУ за сегодня",
			))
			return
		}

		result.WriteString(mealCategoryName(meal.Category))
		result.WriteString("\n")

		for _, item := range mealItems {
			food, err := h.serviceFood.GetFoodByID(
				ctx,
				item.FoodID,
			)
			if err != nil {
				log.Println("Ошибка получения продуктов:", err)
				h.bot.Send(tgbotapi.NewMessage(
					chatID,
					"Не удалось рассчитать КБЖУ за сегодня",
				))
				return
			}

			nutrition := service.CalculateNutrition(*food, item.Weight)

			total.Calories += nutrition.Calories
			total.Protein += nutrition.Protein
			total.Fat += nutrition.Fat
			total.Carbs += nutrition.Carbs

			result.WriteString(fmt.Sprintf("%s - %.0f г\n", food.Name, item.Weight))
			editButton := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("Изменить: %s", food.Name),
				fmt.Sprintf("edit_meal_item_%d", item.ID),
			)

			deleteButton := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("Удалить: %s", food.Name),
				fmt.Sprintf("delete_meal_item_%d", item.ID),
			)

			deleteButtons = append(
				deleteButtons,
				tgbotapi.NewInlineKeyboardRow(
					editButton,
					deleteButton,
				),
			)
		}

		result.WriteString("\n")
	}

	result.WriteString("Итого за день:\n")
	result.WriteString(fmt.Sprintf(
		"Калории: %.1f ккал\nБелки: %.1f г\nЖиры: %.1f г\nУглеводы: %.1f г\n",
		total.Calories,
		total.Protein,
		total.Fat,
		total.Carbs,
	))

	if hasGoal {
		result.WriteString("\nДневная цель:\n")
		result.WriteString(fmt.Sprintf(
			"Калории: %.1f / %.1f ккал\nБелки: %.1f / %.1f г\nЖиры: %.1f / %.1f г\nУглеводы: %.1f / %.1f г\n",
			total.Calories, goal.Calories,
			total.Protein, goal.Protein,
			total.Fat, goal.Fat,
			total.Carbs, goal.Carbs,
		))

		result.WriteString("\n")
		result.WriteString("\nОсталось до цели:\n")

		result.WriteString(fmt.Sprintf(
			"Калории: %s\nБелки: %s\nЖиры: %s\nУглеводы: %s\n",
			formatDifference(goal.Calories-total.Calories, "ккал"),
			formatDifference(goal.Protein-total.Protein, "г"),
			formatDifference(goal.Fat-total.Fat, "г"),
			formatDifference(goal.Carbs-total.Carbs, "г"),
		))
	} else {
		result.WriteString("\nДневная цель не задана, нажми «Задать цель»")
	}

	msg := tgbotapi.NewMessage(chatID, result.String())

	if len(deleteButtons) > 0 {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(deleteButtons...)
	}

	if _, err := h.bot.Send(msg); err != nil {
		log.Println("Ошибка отправки итогов за день:", err)
	}
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

func formatDifference(difference float64, callories string) string {
	if difference >= 0 {
		return fmt.Sprintf("%.1f %s", difference, callories)
	}

	return fmt.Sprintf("Переел на: %.1f %s", -difference, callories)
}
