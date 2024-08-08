package utils

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"regexp"
	"strconv"
	"strings"
	"weatherbot/internal/logger"
	"weatherbot/internal/weather"
)

// GetCityInfo - returns city information like latitude/longitude
// city in config may be like "Moscow[30.9768 60.3456]" (geolocation in brackets)
// so it tries to parse coordinates. if no coordinates then get it via api
func GetCityInfo(city string, geoCoder weather.GeoCoderInterface) (cityInfo *weather.CityInfo, err error) {
	re := regexp.MustCompile(`^(.*?)(?:\[(\d+\.\d+)\s+(\d+\.\d+)\])?$`)
	matches := re.FindStringSubmatch(city)
	if matches == nil {
		return cityInfo, fmt.Errorf("wrong city format: %s", city)
	}

	cityName := strings.TrimSpace(matches[1])
	cityInfo = &weather.CityInfo{
		Name: cityName,
	}
	if len(matches) > 2 && matches[2] != "" && matches[3] != "" {
		lat, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			return cityInfo, fmt.Errorf("wrong latitude: %s", matches[2])
		}
		lon, err := strconv.ParseFloat(matches[3], 64)
		if err != nil {
			return cityInfo, fmt.Errorf("wrong longitude: %s", matches[3])
		}
		cityInfo.Latitude = lat
		cityInfo.Longitude = lon
		cityInfo.HasCoords = true
	}
	if !cityInfo.HasCoords {
		geoData, err := GetGeoCoderData(cityName, geoCoder)
		if err != nil {
			return cityInfo, err
		}
		cityInfo.Latitude = geoData.Latitude
		cityInfo.Longitude = geoData.Longitude
		cityInfo.HasCoords = true
	}

	return cityInfo, nil
}

// GetGeoCoderData - get city geolocation by api
// and save it to local cache
func GetGeoCoderData(city string, geoCoder weather.GeoCoderInterface) (cityInfo *weather.CityInfo, err error) {
	cacheKey := fmt.Sprintf("geocode_%s", city)
	if cacheData, found := geoCoder.GetCacheInstance().Get(cacheKey); found {
		cityInfo = cacheData.(*weather.CityInfo)
	} else {
		cityInfo, err = geoCoder.GetGeoCodeCityInfo(city)
		if err != nil {
			logger.Logger().Print("err:", err)
			return nil, err
		}
		geoCoder.GetCacheInstance().Set(cacheKey, cityInfo, cache.NoExpiration)
	}
	return
}
