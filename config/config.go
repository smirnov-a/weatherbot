package config

import (
	"github.com/spf13/viper"
	"strings"
	"weatherbot/internal/logger"
)

func IniConfig() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		logger.Logger().Fatalf("Error reading config file, %s", err)
	}
	viper.AutomaticEnv()
}

func GetConfigValue(key string) string {
	return viper.GetString(key)
}

func GetApiKey() string {
	provider := GetConfigValue("WEATHER_PROVIDER")
	key := strings.ToUpper(provider) + "_API_KEY"
	return GetConfigValue(key)
}

func GetTelegramToken() string {
	return GetConfigValue("TELEGRAM_TOKEN")
}

func GetTelegramChatId() int64 {
	return viper.GetInt64("TELEGRAM_CHAT_ID")
}

func GetTelegramDebug() bool {
	return viper.GetBool("TELEGRAM_DEBUG")
}
