package handler

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

func (h *handler) DeleteMealItemCallback(ctx context.Context, update tgbotapi.Update) {
	callback := update.CallbackQuery

	if callback == nil {
		return
	}

	answer := func(text string) {
		_, err := h.bot.Request(
			tgbotapi.NewCallback(callback.ID, text),
		)
		if err != nil {
			log.Println("Ошибка ответа на кнопку:", err)
		}
	}

	const delete = "delete_meal_item_"

	if !strings.HasPrefix(callback.Data, delete) {
		answer("Некорректная кнопка")
		return
	}

	itemID, err := strconv.ParseInt(strings.TrimPrefix(callback.Data, delete), 10, 64)
	if err != nil || itemID <= 0 {
		answer("Некорректный ID продукта")
		return
	}

	item, err := h.serviceMealItem.GetMealItemByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			answer("Продукт не найден")
			return
		}

		log.Println("Ошибка получения продукта:", err)
		answer("Не удалось получить продукт")
		return
	}

	meal, err := h.serviceMeal.GetMealByID(ctx, item.MealID)
	if err != nil {
		log.Println("Ошибка получения приёма пищи:", err)
		answer("Не удалось проверить владельца записи")
		return
	}

	user, err := h.serviceUser.GetUserByTelegramID(ctx, callback.From.ID)
	if err != nil {
		log.Println("Ошибка получения пользователя:", err)
		answer("Не удалось определить пользователя")
		return
	}

	if meal.UserID != user.ID {
		answer("Ты че!?Нельзя удалять чужие записи")
		return
	}

	err = h.serviceMealItem.DeleteMealItem(ctx, itemID)
	if err != nil {
		log.Println("Ошибка удаления продукта:", err)
		answer("Не удалось удалить продукт")
		return
	}

	answer("Продукт удалён")

	if callback.Message == nil {
		return
	}

	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	_, err = h.bot.Request(
		tgbotapi.NewDeleteMessage(chatID, messageID),
	)
	if err != nil {
		log.Println("Ошибка удаления старого сообщения:", err)
	}

	message := *callback.Message
	message.From = callback.From

	refreshUpdate := tgbotapi.Update{
		Message: &message,
	}

	h.Today(ctx, refreshUpdate)
}
