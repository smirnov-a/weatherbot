package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sync"
	"weatherbot/internal/weather"
	"weatherbot/utils"
)

const weatherUrl = "https://api.weatherapi.com/v1/current.json"

func (api *WeatherAPI) GetCurrentWeatherData(cityInfo *weather.CityInfo, wg *sync.WaitGroup, ch chan<- *weather.CurrentData, errCh chan<- error) {
	const method = "GetCurrentWeatherData"

	defer func() {
		if r := recover(); r != nil {
			errCh <- fmt.Errorf("panic in %s: %v", method, r)
		}
		wg.Done()
	}()

	additional := map[string]string{
		"current_fields": "temp_c,feelslike_c,wind_kph,wind_degree,wind_dir,pressure_mb,condition",
	}
	url, _ := utils.GetUrl(weatherUrl, cityInfo, api, &additional)
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

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		errCh <- fmt.Errorf("%s. error Unmarshal result: %w", method, err)
		return
	}

	data := &weather.CurrentData{
		City:    cityInfo.Name,
		Weather: math.Round(result["current"].(map[string]interface{})["temp_c"].(float64)),
	}

	ch <- data
}
