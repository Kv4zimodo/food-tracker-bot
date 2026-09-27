package handler

func HandlerHelp(update tgbotapi.Update) {
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