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
func GetWeatherDataImpl[T weather.WeatherDataInterface](ctx context.Context, city string, w T, ch chan<- *weather.WeatherData, wg *sync.WaitGroup) {
	defer wg.Done()

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
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	wg1 := &sync.WaitGroup{}
	wg1.Add(2)
	go w.GetCurrentWeatherData(cityInfo, wg1, ch1, errCh)
	go w.GetWeatherDataForecast(cityInfo, wg1, ch2, errCh)

	go func() {
		wg1.Wait()
		close(ch1)
		close(ch2)
		close(errCh)
	}()

	// build result weather: current and forecast
	var combinedErr error
	result := &weather.WeatherData{}
	done := false
	for !done && (ch1 != nil || ch2 != nil) {
		select {
		case <-ctx.Done():
			combinedErr = ctx.Err()
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
		return
	}

	// output result to data weather channel
	ch <- result
}
