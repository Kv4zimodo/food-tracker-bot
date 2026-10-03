package handler

import (
	"context"
	"errors"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"github.com/kv4zimodo/food-tracker-bot/internal/models"
)

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

	msg.ReplyMarkup = MainMenu()

	if _, err := h.bot.Send(msg); err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
}
