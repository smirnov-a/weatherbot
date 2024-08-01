package openweathermap

import (
	"encoding/json"
	"io"
	"net/http"
	"weatherbot/internal/weather"
	"weatherbot/utils"
)

const geoCodeUrl = "https://api.openweathermap.org/geo/1.0/direct"

func (owm *OpenWeatherMap) GetGeoCodeCityInfo(city string) (*weather.CityInfo, error) {
	const method = "GetGeoCodeCityInfo"
	var geoData weather.CityInfo

	params := &utils.RequestParams{
		Method:      http.MethodGet,
		Url:         geoCodeUrl,
		QueryParams: owm.GetGeoCodingParams(city),
	}
	req, err := utils.NewRequest(params)
	if err != nil {
		return nil, err
	}

	response, err := utils.DoRequestWithRetry(req, utils.Retries, utils.RetryTimeout)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		owm.Logger.Printf("%s. Error reading response body: %v", method, err)
		return nil, err
	}

	var result []map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		owm.Logger.Printf("%s. Error parse response body: %v", method, err)
		return nil, err
	}

	geoData = weather.CityInfo{
		Name:      city,
		Latitude:  result[0]["lat"].(float64),
		Longitude: result[0]["lon"].(float64),
		HasCoords: true,
	}
	return &geoData, err
}
