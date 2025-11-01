package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	apiKey string
	chatID int64
}

func New(apiKey string, chatID int64) *TelegramService {
	return &TelegramService{
		apiKey: apiKey,
		chatID: chatID,
	}
}

func (svc *TelegramService) SendNotification(title string, description string) {
	bot, err := tgbotapi.NewBotAPI(svc.apiKey)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal("There was an error when parsing the telegram chat id")
	}

	messageText := title + "\n" + description
	msg := tgbotapi.NewMessage(svc.chatID, messageText)

	_, err = bot.Send(msg)
	if err != nil {
		log.Fatal("Failed to send message:", err)
	}
}
//
// func SendNotification(title string, description string, apiKey string, telegramChatID int64) {
// 	bot, err := tgbotapi.NewBotAPI(apiKey)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	if err != nil {
// 		log.Fatal("There was an error when parsing the telegram chat id")
// 	}
//
// 	messageText := title + "\n" + description
// 	msg := tgbotapi.NewMessage(telegramChatID, messageText)
//
// 	_, err = bot.Send(msg)
// 	if err != nil {
// 		log.Fatal("Failed to send message:", err)
// 	}
// }
