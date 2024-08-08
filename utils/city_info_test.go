package utils

import (
	"errors"
	"github.com/patrickmn/go-cache"
	"testing"
	"weatherbot/internal/weather"
)

type MockGeocoder struct{}

func (m *MockGeocoder) GetGeoCodeCityInfo(city string) (*weather.CityInfo, error) {
	var lat, lon float64
	var hasCoords bool
	var err error

	switch city {
	case "Yekaterinburg":
		lat = 50.001
		lon = 60.001
		hasCoords = true
	default:
		err = errors.New("Error reading response body")
	}
	return &weather.CityInfo{
		Name:      city,
		Latitude:  lat,
		Longitude: lon,
		HasCoords: hasCoords,
	}, err
}

func (m *MockGeocoder) GetCacheInstance() *cache.Cache {
	return cache.New(cache.NoExpiration, cache.NoExpiration)
}

func TestGetCityInfo(t *testing.T) {
	type test struct {
		city string
		want weather.CityInfo
		err  error
	}

	tests := []test{
		{
			city: "Moscow[30.9768 60.3456]",
			want: weather.CityInfo{
				Name:      "Moscow",
				Latitude:  30.9768,
				Longitude: 60.3456,
				HasCoords: true,
			},
			err: nil,
		},
		{
			city: "Yekaterinburg",
			want: weather.CityInfo{
				Name:      "Yekaterinburg",
				Latitude:  50.001,
				Longitude: 60.001,
				HasCoords: true,
			},
			err: nil,
		},
		{
			city: "nonexistent",
			want: weather.CityInfo{
				Name:      "nonexistent",
				Latitude:  0.0,
				Longitude: 0.0,
				HasCoords: false,
			},
			err: errors.New("Error reading response body"),
		},
	}
	for _, tt := range tests {
		m := &MockGeocoder{}
		got, err := GetCityInfo(tt.city, m)
		if got == nil || *got != tt.want || (err != nil && err.Error() != tt.err.Error()) {
			t.Errorf("GetCityInfo(%s) = %v, %v; want %v, %v", tt.city, got, err, tt.want, tt.err)
		}
	}
}
