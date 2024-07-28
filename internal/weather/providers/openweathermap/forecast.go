package openweathermap

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sync"
	"time"
	"weatherbot/internal/weather"
	"weatherbot/utils"
)

const forecastUrl = "https://api.openweathermap.org/data/2.5/forecast"

// limitOfResult count of items in forecast
const limitOfResult = "10"

func (owm *OpenWeatherMap) GetWeatherDataForecast(cityInfo *weather.CityInfo, wg *sync.WaitGroup, ch chan<- *weather.ForecastData, errCh chan<- error) { //(data weather.WeatherData, err error) {
	const method = "GetWeatherDataForecast"

	defer func() {
		if r := recover(); r != nil {
			errCh <- fmt.Errorf("panic in %s: %v", method, r)
		}
		wg.Done()
	}()

	additional := map[string]string{
		"cnt": limitOfResult,
	}
	url, _ := utils.GetUrl(forecastUrl, cityInfo, owm, &additional)
	client := utils.GetHttpClient()
	response, err := client.Get(url)
	if err != nil {
		errCh <- fmt.Errorf("%s. error fetching data: %w", method, err)
		return
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		errCh <- fmt.Errorf("%s. error read response: %w", method, err)
		return
	}

	var weatherResponse WeatherResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		errCh <- fmt.Errorf("%s. error Unmarshal result: %w", method, err)
		return
	}

	offset := weatherResponse.City.Timezone
	sunrise := weatherResponse.City.Sunrise
	sunset := weatherResponse.City.Sunset

	data := &weather.ForecastData{
		Days:    getDayDiff(weatherResponse),
		Offset:  offset,
		Sunrise: getLocalTime(sunrise, offset),
		Sunset:  getLocalTime(sunset, offset),
	}

	for _, item := range weatherResponse.List {
		row := weather.Row{
			Timestamp:     getLocalTime(item.Dt, offset),
			Temperature:   math.Round(item.Main.Temp),
			FeelsLike:     math.Round(item.Main.FeelsLike),
			Pressure:      convertPaToMmHg(item.Main.Pressure),
			Humidity:      item.Main.Humidity,
			Weather:       item.Weather[0].Description,
			Clouds:        item.Clouds.All,
			Visibility:    item.Visibility,
			Precipitation: getPrecipitation(item),
			Pop:           getPercentOfValue(item.Pop),
			Wind: weather.Wind{
				Speed: item.Wind.Speed,
				Deg:   item.Wind.Deg,
				Gust:  item.Wind.Gust,
			},
		}
		data.Rows = append(data.Rows, row)
	}

	ch <- data
}

// getPrecipitation try to get Rain first, Snow second
func getPrecipitation(item ListItem) float64 {
	if item.Rain.ThreeH > 0 {
		return item.Rain.ThreeH
	} else if item.Snow.ThreeH > 0 {
		return item.Snow.ThreeH
	}
	return 0
}

// convertPaToMmHg pressure hPa to mmHg
func convertPaToMmHg(pressurePa float64) float64 {
	const pascalToMmHg = 1.33322
	return math.Round(pressurePa / pascalToMmHg)
}

func getLocalTime(timestamp, offset int64) string {
	utcTime := time.Unix(timestamp, 0)
	return utcTime.Format(time.DateTime)
}

func getDayDiff(data WeatherResponse) (days float64) {
	if len(data.List) == 0 {
		return
	}

	currTime := time.Now()
	lastTimestamp := data.List[len(data.List)-1].Dt
	endTime := time.Unix(lastTimestamp, 0)
	duration := endTime.Sub(currTime)
	days = duration.Hours() / 24
	return math.Ceil(days)
}

func getPercentOfValue(value float64) string {
	return fmt.Sprintf("%.0f", value*100)
}
