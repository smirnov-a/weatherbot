package openweathermap

import (
	"fmt"
	"weatherbot/internal/weather"
)

// GetUrlParams returns map with parameters for api call
func (owm *OpenWeatherMap) GetUrlParams(cityInfo *weather.CityInfo) *map[string]string {
	params := map[string]string{
		"appid": owm.APIKey,
		"units": "metric",
		"lang":  "ru",
	}
	if cityInfo != nil && cityInfo.Latitude != 0 && cityInfo.Longitude != 0 {
		params["lat"] = fmt.Sprintf("%f", cityInfo.Latitude)
		params["lon"] = fmt.Sprintf("%f", cityInfo.Longitude)
	}

	return &params
}
