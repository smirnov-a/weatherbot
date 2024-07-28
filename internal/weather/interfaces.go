package weather

import (
	"github.com/patrickmn/go-cache"
	"sync"
)

// WeatherDataInterface main interface
type WeatherDataInterface interface {
	GetWeatherData(city string) *WeatherData
	GetCurrentWeatherData(cityInfo *CityInfo, wg *sync.WaitGroup, ch chan<- *CurrentData, errCh chan<- error)
	GetWeatherDataForecast(cityInfo *CityInfo, wg *sync.WaitGroup, ch chan<- *ForecastData, errCh chan<- error)
	GeoCoderInterface
}

// GeoCoderInterface interface uses while working with geolocation api
type GeoCoderInterface interface {
	GetGeoCodeCityInfo(city string) (*CityInfo, error)
	CacheInterface
}

type CacheInterface interface {
	GetCacheInstance() *cache.Cache
}

// UrlParamsInterface interface for working with url parameters for api
type UrlParamsInterface interface {
	GetUrlParams(cityInfo *CityInfo) *map[string]string
	GetGeoCodingParams(city string) *map[string]string
}
