package openweathermap

import (
	"encoding/json"
	"fmt"
	"io"
	"weatherbot/internal/weather"
	"weatherbot/utils"
)

const geoCodeUrl = "https://api.openweathermap.org/geo/1.0/direct"
const limitParam = "10"
const country = "RU"
const state = ""

func (owm *OpenWeatherMap) GetGeoCodeCityInfo(city string) (*weather.CityInfo, error) {
	const method = "GetGeoCodeCityInfo"
	var geoData weather.CityInfo

	additionalParams := owm.GetGeoCodingParams(city)
	url, _ := utils.GetUrl(geoCodeUrl, nil, owm, additionalParams)
	client := utils.GetHttpClient()
	resp, err := client.Get(url)
	if err != nil {
		owm.Logger.Printf("%s. Error making request: %v", method, err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
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

func (owm *OpenWeatherMap) GetGeoCodingParams(city string) *map[string]string {
	return &map[string]string{
		"q":     fmt.Sprintf("%s,%s,%s", city, state, country),
		"limit": limitParam,
	}
}
