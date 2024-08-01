package handler

import (
	"context"
	"fmt"
	"sync"
	"time"
	"weatherbot/internal/logger"
	"weatherbot/internal/weather"
	"weatherbot/utils"
)

const timeout = 45 * time.Second

// GetWeatherDataImpl implement of GetWeatherData by given weather provider
func GetWeatherDataImpl(city string, w weather.WeatherDataInterface) (result *weather.WeatherData) {
	var wg sync.WaitGroup

	ch1 := make(chan *weather.CurrentData)
	ch2 := make(chan *weather.ForecastData)
	errCh := make(chan error)

	defer func() {
		if r := recover(); r != nil {
			logger.Logger().Printf("Recovered from panic: %v", r)
			close(ch1)
			close(ch2)
			close(errCh)
		}
	}()

	// city main contains coordinates in form Yekaterinburg[51.456 60.560]
	// so take coordinates from name or make geolocation api-call
	cityInfo, err := utils.GetCityInfo(city, w)
	if err != nil {
		return &weather.WeatherData{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	wg.Add(2)
	go w.GetCurrentWeatherData(cityInfo, &wg, ch1, errCh)
	go w.GetWeatherDataForecast(cityInfo, &wg, ch2, errCh)

	go func() {
		wg.Wait()
		close(ch1)
		close(ch2)
		close(errCh)
	}()

	var combinedErr error
	result = &weather.WeatherData{}
	done := false
	for !done && (ch1 != nil || ch2 != nil) {
		select {
		case <-ctx.Done():
			logger.Logger().Printf("Context error: %v\n", ctx.Err())
			done = true
		case data, ok := <-ch1:
			if ok {
				result.CurrentData = data
			} else {
				ch1 = nil
			}
		case data, ok := <-ch2:
			if ok {
				result.ForecastData = data
			} else {
				ch2 = nil
			}
		case err, ok := <-errCh:
			if ok {
				if combinedErr == nil {
					combinedErr = err
				} else {
					combinedErr = fmt.Errorf("%v; %w", combinedErr, err)
				}
			} else {
				errCh = nil
			}
		}
	}

	if combinedErr != nil {
		logger.Logger().Errorf("Encountered errors: %v\n", combinedErr)
		return &weather.WeatherData{}
	}

	return
}
