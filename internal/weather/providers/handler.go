package providers

import (
	"context"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	wg := &sync.WaitGroup{}
	// chanData weather data
	chanData := make(chan *weather.WeatherData)
	// chanMessage channel for sending message to telegram
	chanMessage := make(chan *weather.WeatherData)
	once := &sync.Once{}
	closeDataChan := func(ch chan *weather.WeatherData) {
		once.Do(func() {
			close(ch)
		})
	}
	defer func() {
		closeDataChan(chanData)
		close(chanMessage)
	}()

	sendMessageFunc := func(data *weather.WeatherData) {
		message.SendMessageToTelegram(app, data)
	}
	go worker(sendMessageFunc, chanMessage)

	provider := getProvider(app.Cache)
	for _, city := range cities {
		wg.Add(1)
		go provider.GetWeatherData(ctx, city, chanData, wg)
	}

	go func() {
		wg.Wait()
		// close data channel. when closed it will stop cycle below
		closeDataChan(chanData)
	}()

	done := false
	for !done {
		select {
		case <-ctx.Done():
			closeDataChan(chanData)
			done = true
		case data, ok := <-chanData:
			if !ok {
				done = true
				break
			}
			// send data to channel. it will be sent to telegram
			chanMessage <- data
			res = append(res, data)
		}
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

// worker read data from data channel and execute given function with data
func worker(f func(data *weather.WeatherData), ch <-chan *weather.WeatherData) {
	// do while channel is not closed
	for data := range ch {
		f(data)
	}
}
