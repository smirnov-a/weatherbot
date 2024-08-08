package weather

import (
	"context"
	"github.com/patrickmn/go-cache"
	"sync"
)

// WeatherDataInterface main interface
type WeatherDataInterface interface {
	GetWeatherData(context.Context, string, chan<- *WeatherData, *sync.WaitGroup)
	GetCurrentWeatherData(*CityInfo, *sync.WaitGroup, chan<- *CurrentData, chan<- error)
	GetWeatherDataForecast(*CityInfo, *sync.WaitGroup, chan<- *ForecastData, chan<- error)
	GeoCoderInterface
}

// GeoCoderInterface interface uses while working with geolocation api
type GeoCoderInterface interface {
	GetGeoCodeCityInfo(string) (*CityInfo, error)
	CacheInterface
}

// CacheInterface cache instance
type CacheInterface interface {
	GetCacheInstance() *cache.Cache
}

// UrlParamsInterface interface for working with url parameters for api
type UrlParamsInterface interface {
	GetUrlParams(*CityInfo) *map[string]string
	GetGeoCodingParams(string) *map[string]string
}
