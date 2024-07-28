package utils

import (
	"net/url"
	"weatherbot/internal/weather"
)

func GetUrl(urlBase string, cityInfo *weather.CityInfo, params weather.UrlParamsInterface, additional *map[string]string) (string, error) {
	u, _ := url.Parse(urlBase)
	query := u.Query()
	for key, value := range *params.GetUrlParams(cityInfo) {
		query.Add(key, value)
	}
	if additional != nil {
		for key, value := range *additional {
			query.Add(key, value)
		}
	}
	u.RawQuery = query.Encode()
	return u.String(), nil
}
