package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"github.com/kv4zimodo/food-tracker-bot/internal/models"
	"github.com/kv4zimodo/food-tracker-bot/internal/service"
)

type userState struct {
	Meal        string
	MealItemID  int64
	State       string
	FoodName    string
	FoodID      int64
	Weight      int
	GoalProtein float64
	GoalFat     float64
}

var states = make(map[int64]userState)

func FoodMeals() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Завтрак", "meal_breakfast"),
			tgbotapi.NewInlineKeyboardButtonData("Обед", "meal_lunch"),
		), tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ужин", "meal_dinner"),
			tgbotapi.NewInlineKeyboardButtonData("Перекус", "meal_snack"),
		),
	)
}

func (h *handler) AddFood(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID

	msg := tgbotapi.NewMessage(chatID, "Выберите прием пищи")

	msg.ReplyMarkup = FoodMeals()

	if _, err := h.bot.Send(msg); err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
}

func (h *handler) MealCallBack(update tgbotapi.Update) {
	telegramID := update.CallbackQuery.From.ID
	chatID := update.CallbackQuery.Message.Chat.ID

	switch update.CallbackData() {
	case "meal_breakfast":
		states[telegramID] = userState{
			Meal:  "breakfast",
			State: "waiting_food",
		}
	case "meal_lunch":
		states[telegramID] = userState{
			Meal:  "lunch",
			State: "waiting_food",
		}
	case "meal_dinner":
		states[telegramID] = userState{
			Meal:  "dinner",
			State: "waiting_food",
		}
	case "meal_snack":
		states[telegramID] = userState{
			Meal:  "snack",
			State: "waiting_food",
		}
	default:
		return
	}

	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	if _, err := h.bot.Request(callback); err != nil {
		log.Println("Ошибка подтверждения callback:", err)
	}

	msg := tgbotapi.NewMessage(chatID, "Введите название продукта")

	if _, err := h.bot.Send(msg); err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
}

func (h *handler) FoodInput(ctx context.Context, update tgbotapi.Update) {
	telegramID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	foodName := update.Message.Text

	state := states[telegramID]

	if state.State != "waiting_food" {
		return
	}

	food, err := h.serviceFood.GetFoodByName(ctx, foodName)
	if err != nil {
		log.Println("Ошибка поиска продукта:", err)

		state.State = "waiting_food"
		states[telegramID] = state

		msg := tgbotapi.NewMessage(chatID, "Такого продукта нет в базе, введите название продукта снова")
		if _, err := h.bot.Send(msg); err != nil {
			log.Println("Ошибка отправки сообщения:", err)
		}

		return
	}

	state.FoodID = food.ID
	state.FoodName = food.Name
	state.State = "waiting_weight"
	states[telegramID] = state

	msg := tgbotapi.NewMessage(chatID, "Введите вес продукта в граммах")

	if _, err := h.bot.Send(msg); err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
}

func (h *handler) WeightInput(ctx context.Context, update tgbotapi.Update) {
	telegramID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	state := states[telegramID]

	if state.State != "waiting_weight" {
		return
	}

	weight, err := strconv.Atoi(update.Message.Text)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Введите вес числом, например: 250")
		if _, err := h.bot.Send(msg); err != nil {
			log.Println("Ошибка отправки сообщения:", err)
		}
		return
	}
	if weight <= 0 {
		msg := tgbotapi.NewMessage(chatID, "Вес должен быть больше нуля")
		if _, err := h.bot.Send(msg); err != nil {
			log.Println("Ошибка отправки сообщения:", err)
		}
		return
	}

	state.Weight = weight
	states[telegramID] = state

	user, err := h.serviceUser.GetUserByTelegramID(
		ctx,
		telegramID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("Пользователь не найден")

			msg := tgbotapi.NewMessage(
				chatID,
				"Пользователь не найден. Используйте /start",
			)

			if _, err := h.bot.Send(msg); err != nil {
				log.Println("Ошибка отправки сообщения:", err)
			}

			return
		}

		log.Println("Ошибка получения пользователя:", err)

		msg := tgbotapi.NewMessage(
			chatID,
			"Произошла ошибка при получении пользователя",
		)

		if _, err := h.bot.Send(msg); err != nil {
			log.Println("Ошибка отправки сообщения:", err)
		}

		return
	}

	category := models.MealCategory(state.Meal)

	now := time.Now()

	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location())

	meal, err := h.serviceMeal.GetMealByUserAndCategoryAndDate(
		ctx,
		user.ID,
		category,
		today)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			meal, err = h.serviceMeal.CreateMeal(
				ctx,
				models.Meal{
					UserID:   user.ID,
					Category: category,
					Date:     today,
				},
			)

			if err != nil {
				log.Println("Ошибка создания приема пищи:", err)

				msg := tgbotapi.NewMessage(
					chatID,
					"Не удалось сохранить прием пищи",
				)

				if _, err := h.bot.Send(msg); err != nil {
					log.Println("Ошибка отправки сообщения:", err)
				}

				return
			}
		} else {
			log.Println("Ошибка поиска приема пищи:", err)

			msg := tgbotapi.NewMessage(
				chatID,
				"Не удалось найти прием пищи",
			)

			if _, err := h.bot.Send(msg); err != nil {
				log.Println("Ошибка отправки сообщения:", err)
			}

			return
		}
	}

	food, err := h.serviceFood.GetFoodByID(ctx, state.FoodID)
	if err != nil {
		log.Println("Ошибка получения продукта:", err)

		msg := tgbotapi.NewMessage(chatID, "Не удалось получить продукт из базы, попробуй добавить его заново.")
		if _, err := h.bot.Send(msg); err != nil {
			log.Println("Ошибка отправки сообщения:", err)
		}

		return
	}

	_, err = h.serviceMealItem.CreateMealItem(
		ctx,
		models.MealItem{
			MealID: meal.ID,
			FoodID: food.ID,
			Weight: float64(weight),
		},
	)
	if err != nil {
		log.Println("Ошибка создания приема пищи:", err)

		msg := tgbotapi.NewMessage(
			chatID,
			"Не удалось сохранить продукт",
		)

		if _, err := h.bot.Send(msg); err != nil {
			log.Println("Ошибка отправки сообщения:", err)
		}

		return
	}

	nutrition := service.CalculateNutrition(*food, float64(weight))

	text := fmt.Sprintf(
		"Добавлено в прием пищи\n\n"+
			"%s — %d г\n\n"+
			"Калории: %.1f ккал\n"+
			"Белки: %.1f г\n"+
			"Жиры: %.1f г\n"+
			"Углеводы: %.1f г",
		food.Name,
		weight,
		nutrition.Calories,
		nutrition.Protein,
		nutrition.Fat,
		nutrition.Carbs,
	)

	delete(states, telegramID)

	msg := tgbotapi.NewMessage(chatID, text)

	if _, err := h.bot.Send(msg); err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
}
