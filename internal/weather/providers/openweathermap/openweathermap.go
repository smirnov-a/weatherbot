package openweathermap

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"sync"
	"weatherbot/internal/weather"
	"weatherbot/internal/weather/handler"
)

type OpenWeatherMap struct {
	APIKey string
	Cache  *cache.Cache
	Logger *logrus.Logger
}

func (owm *OpenWeatherMap) GetWeatherData(ctx context.Context, city string, ch chan<- *weather.WeatherData, wg *sync.WaitGroup) {
	handler.GetWeatherDataImpl(ctx, city, owm, ch, wg)
}

func (owm *OpenWeatherMap) GetCacheInstance() *cache.Cache {
	return owm.Cache
}
