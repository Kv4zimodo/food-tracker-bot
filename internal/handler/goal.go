package handler

import (
	"context"
	"errors"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

func (h *handler) Goal(ctx context.Context, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	telegramID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	user, err := h.serviceUser.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		log.Println("Ошибка получения пользователя:", err)
		h.bot.Send(tgbotapi.NewMessage(
			chatID,
			"Не удалось получить данные пользователя",
		))
		return
	}

	goal, err := h.serviceGoal.GetGoalByUserID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.bot.Send(tgbotapi.NewMessage(
				chatID,
				"У тебя пока не задана дневная цель",
			))
			return
		}

		log.Println("Ошибка получения цели:", err)
		h.bot.Send(tgbotapi.NewMessage(
			chatID,
			"Не удалось получить дневную цель",
		))
		return
	}

	text := fmt.Sprintf(
		"Твоя дневная цель\n\n"+
			"Калории: %.0f ккал\n"+
			"Белки: %.1f г\n"+
			"Жиры: %.1f г\n"+
			"Углеводы: %.1f г",
		goal.Calories,
		goal.Protein,
		goal.Fat,
		goal.Carbs,
	)

	h.bot.Send(tgbotapi.NewMessage(chatID, text))
}
