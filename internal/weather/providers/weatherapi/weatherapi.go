package weatherapi

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"sync"
	"weatherbot/internal/weather"
	"weatherbot/internal/weather/handler"
)

type WeatherAPI struct {
	APIKey string
	Cache  *cache.Cache
	Logger *logrus.Logger
}

func (api *WeatherAPI) GetWeatherData(ctx context.Context, city string, ch chan<- *weather.WeatherData, wg *sync.WaitGroup) {
	handler.GetWeatherDataImpl(ctx, city, api, ch, wg)
}

func (api *WeatherAPI) GetCacheInstance() *cache.Cache {
	return api.Cache
}
