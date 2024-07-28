package app

import (
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"weatherbot/internal/telegram"
)

type AppContext struct {
	TelegramBot *telegram.TelegramBot
	Cache       *cache.Cache
	Crontab     string
	ChatID      int64
	Logger      *logrus.Logger
}
