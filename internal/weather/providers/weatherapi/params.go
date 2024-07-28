package weatherapi

import (
	"fmt"
	"weatherbot/internal/weather"
)

// GetUrlParams returns map with parameters for api call
func (api *WeatherAPI) GetUrlParams(cityInfo *weather.CityInfo) *map[string]string {
	params := map[string]string{
		"key":  api.APIKey,
		"lang": "ru",
		"aqi":  "no",
	}
	if cityInfo != nil && cityInfo.Latitude != 0 && cityInfo.Longitude != 0 {
		params["q"] = fmt.Sprintf("%f,%f", cityInfo.Latitude, cityInfo.Longitude)
	}

	return &params
}
