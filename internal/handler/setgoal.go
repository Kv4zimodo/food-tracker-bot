package handler

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) SetGoal(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Выберите способ ввода цели",
	)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"По калориям",
				"goal_mode_calories",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"По БЖУ",
				"goal_mode_bju",
			),
		),
	)

	if _, err := h.bot.Send(msg); err != nil {
		return
	}
}

func (h *handler) GoalModeCallback(update tgbotapi.Update) {
	query := update.CallbackQuery

	telegramID := query.From.ID

	var text string

	switch query.Data {
	case "goal_mode_calories":
		states[telegramID] = userState{
			State: "waiting_goal_calories",
		}
		text = "Введи дневную цель в килокалориях (ккал):"

	case "goal_mode_bju":
		states[telegramID] = userState{
			State: "waiting_goal_protein",
		}
		text = "Введи дневную норму белков в граммах:"
	}

	callback := tgbotapi.NewCallback(query.ID, "")
	if _, err := h.bot.Request(callback); err != nil {
		return
	}

	msg := tgbotapi.NewMessage(query.Message.Chat.ID, text)
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
		return
	}
}

func (h *handler) GoalInput(ctx context.Context, update tgbotapi.Update) {
	telegramID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	state := states[telegramID]

	value, err := strconv.ParseFloat(update.Message.Text, 64)
	if err != nil {
		msg := tgbotapi.NewMessage(
			chatID,
			"Введи число, например: 120 или 120.5",
		)
		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		}
		return
	}

	switch state.State {
	case "waiting_goal_calories":
		user, err := h.serviceUser.GetUserByTelegramID(ctx, telegramID)
		if err != nil {
			msg := tgbotapi.NewMessage(
				chatID,
				"Не удалось получить пользователя.",
			)
			if _, err := h.bot.Send(msg); err != nil {
				log.Printf("Ошибка отправки сообщения: %v", err)
			}
			return
		}

		goal, err := h.serviceGoal.SetGoalByCalories(ctx, user.ID, value)
		if err != nil {
			msg := tgbotapi.NewMessage(
				chatID,
				"Не удалось сохранить цель: "+err.Error(),
			)
			if _, err := h.bot.Send(msg); err != nil {
				log.Printf("Ошибка отправки сообщения: %v", err)
			}
			return
		}

		delete(states, telegramID)

		text := fmt.Sprintf(
			"Цель сохранена!\n"+
				"Калории: %.0f ккал\n"+
				"Белки: %.1f г\n"+
				"Жиры: %.1f г\n"+
				"Углеводы: %.1f г",
			goal.Calories,
			goal.Protein,
			goal.Fat,
			goal.Carbs,
		)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyMarkup = MainMenu()

		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		}

	case "waiting_goal_protein":
		state.GoalProtein = value
		state.State = "waiting_goal_fat"
		states[telegramID] = state

		msg := tgbotapi.NewMessage(
			chatID,
			"Введи дневную норму жиров в граммах:",
		)
		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		}

	case "waiting_goal_fat":
		state.GoalFat = value
		state.State = "waiting_goal_carbs"
		states[telegramID] = state

		msg := tgbotapi.NewMessage(
			chatID,
			"Введи дневную норму углеводов в граммах:",
		)
		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		}

	case "waiting_goal_carbs":
		user, err := h.serviceUser.GetUserByTelegramID(ctx, telegramID)
		if err != nil {
			msg := tgbotapi.NewMessage(
				chatID,
				"Не удалось получить пользователя.",
			)
			if _, err := h.bot.Send(msg); err != nil {
				log.Printf("Ошибка отправки сообщения: %v", err)
			}
			return
		}

		goal, err := h.serviceGoal.SetGoalByMacros(
			ctx,
			user.ID,
			state.GoalProtein,
			state.GoalFat,
			value,
		)
		if err != nil {
			state.State = "waiting_goal_protein"
			state.GoalProtein = 0
			state.GoalFat = 0
			states[telegramID] = state

			msg := tgbotapi.NewMessage(
				chatID,
				"Не удалось сохранить цель: "+err.Error()+
					"\nВведи норму белков заново:",
			)
			if _, err := h.bot.Send(msg); err != nil {
				log.Printf("Ошибка отправки сообщения: %v", err)
			}
			return
		}

		delete(states, telegramID)

		text := fmt.Sprintf(
			"Цель сохранена!\n"+
				"Калории: %.0f ккал\n"+
				"Белки: %.1f г\n"+
				"Жиры: %.1f г\n"+
				"Углеводы: %.1f г",
			goal.Calories,
			goal.Protein,
			goal.Fat,
			goal.Carbs,
		)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyMarkup = MainMenu()

		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		}
	}
}
