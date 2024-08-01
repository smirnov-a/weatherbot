package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sync"
	"weatherbot/internal/weather"
	"weatherbot/utils"
)

const weatherUrl = "https://api.weatherapi.com/v1/current.json"

// GetCurrentWeatherData get current weather from data provider
func (api *WeatherAPI) GetCurrentWeatherData(cityInfo *weather.CityInfo, wg *sync.WaitGroup, ch chan<- *weather.CurrentData, errCh chan<- error) {
	const method = "GetCurrentWeatherData"

	defer func() {
		if r := recover(); r != nil {
			errCh <- fmt.Errorf("panic in %s: %v", method, r)
		}
		wg.Done()
	}()

	params := &utils.RequestParams{
		Method:      http.MethodGet,
		Url:         weatherUrl,
		QueryParams: api.GetUrlParams(cityInfo),
	}

	req, err := utils.NewRequest(params)
	if err != nil {
		errCh <- fmt.Errorf("%s. error creating request: %w", method, err)
		return
	}

	response, err := utils.DoRequestWithRetry(req, utils.Retries, utils.RetryTimeout)
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

	wData := 0.0
	if resCurrent, found := result["current"]; found {
		wData = math.Round(resCurrent.(map[string]interface{})["temp_c"].(float64))
	}
	data := &weather.CurrentData{
		City:    cityInfo.Name,
		Weather: wData,
	}

	ch <- data
}
