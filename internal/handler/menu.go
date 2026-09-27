package handler

func MainMenu() tgbotapi.ReplyKeyboardMarkup{
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить еду"),
			tgbotapi.NewKeyboardButton("Сегодня"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Моя цель"),
			tgbotapi.NewKeyboardButton("Настройки"),
		),
	)
}