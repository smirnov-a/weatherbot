package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"weatherbot/internal/logger"
)

type TelegramBot struct {
	Bot *tgbotapi.BotAPI
}

func NewTelegramBot(token string) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	return &TelegramBot{
		Bot: bot,
	}, nil
}

func (t *TelegramBot) SendMessage(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := t.Bot.Send(msg)
	if err != nil {
		logger.Logger().Printf("Failed to send message: %v", err)
	}
	return err
}

// HandleUpdates - update mode
func (t *TelegramBot) HandleUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for update := range t.Bot.GetUpdatesChan(u) {
		if update.Message != nil {
			logger.Logger().Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			// Example: echo the message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			t.Bot.Send(msg)
		}
	}
}

// HandleWebhook - webhook mode
func (t *TelegramBot) HandleWebhook(webhookURL string) {
	webhookConfig, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		logger.Logger().Panic(err)
	}

	_, err = t.Bot.Request(webhookConfig)
	if err != nil {
		logger.Logger().Panic(err)
	}

	webhookInfo, err := t.Bot.GetWebhookInfo()
	if err != nil {
		logger.Logger().Panic(err)
	}

	if webhookInfo.LastErrorDate != 0 {
		logger.Logger().Printf("Telegram callback failed: %s", webhookInfo.LastErrorMessage)
	}

	updates := t.Bot.ListenForWebhook("/")

	go http.ListenAndServe(":8080", nil)

	for update := range updates {
		if update.Message != nil {
			logger.Logger().Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			// Example: echo the message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			t.Bot.Send(msg)
		}
	}
}
