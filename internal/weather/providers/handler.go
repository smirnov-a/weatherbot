package providers

import (
	"github.com/patrickmn/go-cache"
	"weatherbot/config"
	"weatherbot/internal/app"
	"weatherbot/internal/logger"
	"weatherbot/internal/telegram/message"
	"weatherbot/internal/weather"
	"weatherbot/internal/weather/providers/openweathermap"
	"weatherbot/internal/weather/providers/weatherapi"
)

const providerOpenweathermap = "openweathermap"
const providerWeatherapi = "weatherapi"

// GetWeather get current and forecast weather for given cities
// and send int to telegram chat
func GetWeather(app *app.AppContext, cities []string) (res []*weather.WeatherData) {
	provider := getProvider(app.Cache)
	for _, city := range cities {
		data := provider.GetWeatherData(city)
		go message.SendMessageToTelegram(app, data)
		res = append(res, data)
	}
	return
}

// getProvider depends on config setting
func getProvider(cache *cache.Cache) (provider weather.WeatherDataInterface) {
	log := logger.Logger()
	prov := config.GetConfigValue("WEATHER_PROVIDER")
	switch prov {
	case providerOpenweathermap:
		provider = &openweathermap.OpenWeatherMap{
			APIKey: config.GetApiKey(),
			Cache:  cache,
			Logger: log,
		}
	case providerWeatherapi:
		provider = &weatherapi.WeatherAPI{
			APIKey: config.GetApiKey(),
			Cache:  cache,
			Logger: log,
		}
	default:
		log.Println("Unknown weather provider:", prov)
	}
	return
}
