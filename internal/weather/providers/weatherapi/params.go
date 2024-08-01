package weatherapi

import (
	"fmt"
	"weatherbot/internal/weather"
)

const limitParam = "10"
const country = "RU"
const state = ""

// GetUrlParams returns map with parameters for api call
func (api *WeatherAPI) GetUrlParams(cityInfo *weather.CityInfo) *map[string]string {
	params := api.getDefaultParams()
	if cityInfo != nil && cityInfo.Latitude != 0 && cityInfo.Longitude != 0 {
		params["q"] = fmt.Sprintf("%f,%f", cityInfo.Latitude, cityInfo.Longitude)
	}

	return &params
}

func (api *WeatherAPI) getDefaultParams() map[string]string {
	return map[string]string{
		"key":  api.APIKey,
		"lang": "ru",
		"aqi":  "no",
	}
}

func (api *WeatherAPI) GetGeoCodingParams(city string) *map[string]string {
	params := api.getDefaultParams()
	params["q"] = fmt.Sprintf("%s,%s,%s", city, state, country)
	params["limit"] = limitParam
	return &params
}
