package openweathermap

import (
	"fmt"
	"weatherbot/internal/weather"
)

const limitParam = "10"
const country = "RU"
const state = ""

// GetUrlParams returns map with parameters for api call
func (owm *OpenWeatherMap) GetUrlParams(cityInfo *weather.CityInfo) *map[string]string {
	params := owm.getDefaultParams()
	if cityInfo != nil && cityInfo.Latitude != 0 && cityInfo.Longitude != 0 {
		params["lat"] = fmt.Sprintf("%f", cityInfo.Latitude)
		params["lon"] = fmt.Sprintf("%f", cityInfo.Longitude)
	}

	return &params
}

func (owm *OpenWeatherMap) getDefaultParams() map[string]string {
	return map[string]string{
		"appid": owm.APIKey,
		"units": "metric",
		"lang":  "ru",
	}
}

func (owm *OpenWeatherMap) GetGeoCodingParams(city string) *map[string]string {
	params := owm.getDefaultParams()
	params["q"] = fmt.Sprintf("%s,%s,%s", city, state, country)
	params["limit"] = limitParam
	return &params
}
