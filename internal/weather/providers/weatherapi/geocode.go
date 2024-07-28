package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"weatherbot/internal/weather"
	"weatherbot/utils"
)

const geoCodeUrl = "https://api.weatherapi.com/v1/search.json"

func (api *WeatherAPI) GetGeoCodeCityInfo(city string) (*weather.CityInfo, error) {
	const method = "GetGeoCodeCityInfo"
	var geoData weather.CityInfo

	additionalParams := api.GetGeoCodingParams(city)
	url, _ := utils.GetUrl(geoCodeUrl, nil, api, additionalParams)
	client := utils.GetHttpClient()
	response, err := client.Get(url)
	if err != nil {
		api.Logger.Printf("%s. Error making request: %v", method, err)
		return nil, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		api.Logger.Printf("%s. Error reading response body: %v", method, err)
		return nil, err
	}

	var result []map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		api.Logger.Printf("%s. Error parse response body: %v", method, err)
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

func (api *WeatherAPI) GetGeoCodingParams(city string) *map[string]string {
	state := ""
	country := "RU"
	return &map[string]string{
		"q":     fmt.Sprintf("%s,%s,%s", city, state, country),
		"limit": "10",
	}
}
