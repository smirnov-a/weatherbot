package weatherapi

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

const forecastUrl = "https://api.weatherapi.com/v1/forecast.json"
const cntDays = "2"
const cntRows = 18

func (api *WeatherAPI) GetWeatherDataForecast(cityInfo *weather.CityInfo, wg *sync.WaitGroup, ch chan<- *weather.ForecastData, errCh chan<- error) {
	const method = "GetWeatherDataForecast"
	defer wg.Done()

	additional := map[string]string{
		"days":        cntDays,
		"hour_fields": "time,temp_c,feelslike_c,pressure_mb,humidity,wind_kph,condition,cloud,vis_km,precip_mm",
	}
	url, _ := utils.GetUrl(forecastUrl, cityInfo, api, &additional)
	client := utils.GetHttpClient()
	response, err := client.Get(url)
	if err != nil {
		errCh <- fmt.Errorf("%s. error fetching data: %w", method, err)
		return
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		errCh <- fmt.Errorf("GetWeatherDataForecast. error read response: %w", err)
		return
	}

	var weatherResponse WeatherResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		errCh <- fmt.Errorf("%s. error Unmarshal result: %w", method, err)
		return
	}

	currentTime := time.Now()
	data := &weather.ForecastData{
		Days:    getDayDiff(weatherResponse, currentTime),
		Sunrise: weatherResponse.Forecast.Forecastday[0].Astro.Sunrise,
		Sunset:  weatherResponse.Forecast.Forecastday[0].Astro.Sunset,
	}

	for _, day := range weatherResponse.Forecast.Forecastday {
		for _, item := range day.Hour {
			itemTime := time.Unix(item.TimeEpoch, 0)
			if itemTime.Before(currentTime) || len(data.Rows) > cntRows {
				continue
			}
			row := weather.Row{
				Timestamp:     getLocalTime(item.TimeEpoch),
				Temperature:   math.Round(item.TempC),
				FeelsLike:     math.Round(item.FeelslikeC),
				Pressure:      convertPaToMmHg(item.PressureMb),
				Humidity:      item.Humidity,
				Weather:       item.Condition.Text,
				Clouds:        item.Cloud,
				Visibility:    int(item.VisKm),
				Precipitation: item.PrecipMm,
				Pop:           getPercentOfValue(item.WillItRain, item.WillItSnow),
				Wind: weather.Wind{
					Speed: kmhToMs(item.WindKph),
					Deg:   item.WindDegree,
				},
			}
			data.Rows = append(data.Rows, row)
		}
	}

	ch <- data
}

// convertPaToMmHg pressure hPa to mmHg
func convertPaToMmHg(pressurePa float64) float64 {
	const pascalToMmHg = 1.33322
	return math.Round(pressurePa / pascalToMmHg)
}

func kmhToMs(kmh float64) float64 {
	return math.Round(kmh / 3.6)
}

func getLocalTime(timestamp int64) string {
	utcTime := time.Unix(timestamp, 0)
	return utcTime.Format(time.DateTime)
}

// getDayDiff count diff days between current date and date from last array element
func getDayDiff(data WeatherResponse, currentTime time.Time) (days float64) {
	if len(data.Forecast.Forecastday) == 0 {
		return
	}

	lastForecastDay := data.Forecast.Forecastday[len(data.Forecast.Forecastday)-1]
	lastHour := lastForecastDay.Hour[len(lastForecastDay.Hour)-1]
	lastTimestamp := lastHour.TimeEpoch
	endTime := time.Unix(lastTimestamp, 0)
	duration := endTime.Sub(currentTime)
	days = duration.Hours() / 24
	return math.Ceil(days)
}

func getPercentOfValue(willRain, willSnow int) string {
	if willRain > 0 {
		return fmt.Sprintf("%.0d", willRain*100)
	} else if willSnow > 0 {
		return fmt.Sprintf("%.0d", willSnow*100)
	}
	return ""
}
