package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/patrickmn/go-cache"
	"os"
	"weatherbot/config"
	"weatherbot/i18n"
	"weatherbot/internal/app"
	"weatherbot/internal/logger"
	"weatherbot/internal/scheduler"
	"weatherbot/internal/telegram"
)

const defaultLang = "en"

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {

	crontabFile := flag.String("crontab", "crontab", "Path to crontab file")
	help := flag.Bool("help", false, "Show help")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if err := checkCronTabFile(*crontabFile); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	logger.InitLogger()
	log := logger.Logger()

	config.IniConfig()
	initLocale()

	bot, err := telegram.NewTelegramBot(config.GetTelegramToken())
	if err != nil {
		log.Fatalf("Failed to create telegram bot: %v", err)
	}

	app := &app.AppContext{
		TelegramBot: bot,
		Cache:       cache.New(cache.NoExpiration, cache.NoExpiration),
		Crontab:     *crontabFile,
		ChatID:      config.GetTelegramChatId(),
		Logger:      log,
		Context:     context.Background(),
	}

	scheduler.Start(app)
}

// initLocale initialize locale
func initLocale() {
	lang := config.GetConfigValue("LANGUAGE")
	if lang != "" {
		i18n.Initialize(defaultLang, lang)
		i18n.SetLocale(lang)
	}
}

func checkCronTabFile(f string) error {
	_, err := os.ReadFile(f)
	return err
}
