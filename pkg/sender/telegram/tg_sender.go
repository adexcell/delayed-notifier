package telegram

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramSender struct {
	botAPI *tgbotapi.BotAPI
}

func New() *TelegramSender {
	botAPI, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		log.Fatal("could not connect to telegram api: ", err)
	}

	botAPI.Debug = false

	return &TelegramSender{
		botAPI: botAPI,
	}
}

func (t *TelegramSender) SendToTelegram(telegramId int, text string) error {
	msg := tgbotapi.NewMessage(int64(telegramId), text)
	_, err := t.botAPI.Send(msg)
	if err != nil {
		return fmt.Errorf("could not send message to telegram user: %s", err.Error())
	}

	return nil
}
