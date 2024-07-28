package weatherapi

import (
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"weatherbot/internal/weather"
	"weatherbot/internal/weather/handler"
)

type WeatherAPI struct {
	APIKey string
	Cache  *cache.Cache
	Logger *logrus.Logger
}

func (api *WeatherAPI) GetWeatherData(city string) (result *weather.WeatherData) {
	return handler.GetWeatherDataImpl(city, api)
}

func (api *WeatherAPI) GetCacheInstance() *cache.Cache {
	return api.Cache
}
