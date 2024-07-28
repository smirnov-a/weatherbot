package message

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"weatherbot/internal/app"
	"weatherbot/internal/weather"
)

const templatePath = "templates/weather.html"

func SendMessageToTelegram(ctx *app.AppContext, data *weather.WeatherData) {
	const method = "SendMessageToTelegram"
	defer func() {
		if r := recover(); r != nil {
			ctx.Logger.Printf("Recovered from panic in %s: %v", method, r)
		}
	}()

	if data.CurrentData == nil || data.ForecastData == nil {
		return
	}

	htmlContent, err := GenerateWeatherHtm(data, templatePath)
	if err != nil {
		ctx.Logger.Printf("%s. Failed to generate HTML: %v", method, err)
		return
	}

	tempFile, err := os.CreateTemp("", "weather_forecast_*.png")
	if err != nil {
		ctx.Logger.Printf("%s. Failed to create temporary file: %v", method, err)
		return
	}

	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	if err := RenderHTMLToImage(htmlContent, tempFile.Name()); err != nil {
		ctx.Logger.Printf("%s. Failed to render HTML to image: %v", method, err)
		return
	}

	photo := tgbotapi.NewPhoto(ctx.ChatID, tgbotapi.FilePath(tempFile.Name()))
	if _, err := ctx.TelegramBot.Bot.Send(photo); err != nil {
		ctx.Logger.Printf("%s. Telegram bot send error: %v", method, err)
	}
}
