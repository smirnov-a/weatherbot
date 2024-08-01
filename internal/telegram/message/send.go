package message

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"weatherbot/internal/app"
	"weatherbot/internal/weather"
)

const templatePath = "templates/weather.html"

// SendMessageToTelegram send message to telegram with weather data
func SendMessageToTelegram(app *app.AppContext, data *weather.WeatherData) {
	const method = "SendMessageToTelegram"
	defer func() {
		if r := recover(); r != nil {
			app.Logger.Printf("Recovered from panic in %s: %v", method, r)
		}
	}()

	if data.CurrentData == nil || data.ForecastData == nil {
		return
	}

	htmlContent, err := GenerateWeatherHtm(data, templatePath)
	if err != nil {
		app.Logger.Printf("%s. Failed to generate HTML: %v", method, err)
		return
	}

	tempFile, err := os.CreateTemp("", "weather_forecast_*.png")
	if err != nil {
		app.Logger.Printf("%s. Failed to create temporary file: %v", method, err)
		return
	}

	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	if err := RenderHTMLToImage(htmlContent, tempFile.Name()); err != nil {
		app.Logger.Printf("%s. Failed to render HTML to image: %v", method, err)
		return
	}

	photo := tgbotapi.NewPhoto(app.ChatID, tgbotapi.FilePath(tempFile.Name()))
	if _, err := app.TelegramBot.Bot.Send(photo); err != nil {
		app.Logger.Printf("%s. Telegram bot send error: %v", method, err)
	}
}
