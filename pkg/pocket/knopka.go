package pocket

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var userStates = make(map[int64]string)

func main() {
	// Создаем бота
	bot, err := tgbotapi.NewBotAPI("7543227307:AAGYpAkSZofJDLv5SIKIm2nETeLO0cxwTzw")
	if err != nil {
		log.Fatalf("Ошибка инициализации бота: %v", err)
	}

	bot.Debug = true
	log.Printf("Авторизован под: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// Обработка нажатий кнопок
		if update.CallbackQuery != nil {
			handleButtonClick(bot, update.CallbackQuery)
		}

		// Обработка текстовых сообщений
		if update.Message != nil {
			handleMessage(bot, update.Message)
		}
	}
}

func handleButtonClick(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	// Сохраняем выбранную кнопку в состоянии пользователя
	//  userID := callback.From.ID
	chatID := callback.Message.Chat.ID
	userStates[chatID] = callback.Data

	// Уведомляем Telegram, что кнопка обработана
	callbackResponse := tgbotapi.NewCallback(callback.ID, "Кнопка выбрана")
	bot.Request(callbackResponse)

	// Просим пользователя ввести PIN
	msg := tgbotapi.NewMessage(chatID, "Введите ваш PIN:")
	bot.Send(msg)
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	//userID := message.From.ID

	// Проверяем, находится ли пользователь в состоянии ожидания ввода PIN
	if selectedButton, exists := userStates[chatID]; exists {
		// Обрабатываем введенный PIN
		pin := message.Text

		// Отправляем ответ с кнопкой и PIN
		response := "Вы выбрали: " + selectedButton + "\nВаш PIN: " + pin
		msg := tgbotapi.NewMessage(chatID, response)
		bot.Send(msg)

		// Убираем состояние пользователя
		delete(userStates, chatID)
		return
	}

	// Если пользователь не нажал кнопку, выводим сообщение по умолчанию
	msg := tgbotapi.NewMessage(chatID, "Сначала выберите кнопку!")
	bot.Send(msg)
}
