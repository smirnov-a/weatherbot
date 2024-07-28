package openweathermap

import (
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"weatherbot/internal/weather"
	"weatherbot/internal/weather/handler"
)

type OpenWeatherMap struct {
	APIKey string
	Cache  *cache.Cache
	Logger *logrus.Logger
}

func (owm *OpenWeatherMap) GetWeatherData(city string) (result *weather.WeatherData) {
	return handler.GetWeatherDataImpl(city, owm)
}

func (owm *OpenWeatherMap) GetCacheInstance() *cache.Cache {
	return owm.Cache
}
