package openweathermap

// WeatherResponse structure of response from service
type WeatherResponse struct {
	Cod     string     `json:"cod"`
	Message int        `json:"message"`
	Cnt     int        `json:"cnt"`
	List    []ListItem `json:"list"`
	City    struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Country    string `json:"country"`
		Population int    `json:"population"`
		Timezone   int64  `json:"timezone"`
		Sunrise    int64  `json:"sunrise"`
		Sunset     int64  `json:"sunset"`
	} `json:"city"`
}

type ListItem struct {
	Dt   int64 `json:"dt"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  float64 `json:"pressure"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
		Humidity  int     `json:"humidity"`
		TempKf    float64 `json:"temp_kf"`
	} `json:"main"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	Visibility int     `json:"visibility"`
	Pop        float64 `json:"pop"`
	Rain       struct {
		ThreeH float64 `json:"3h"`
	} `json:"rain,omitempty"`
	Snow struct {
		ThreeH float64 `json:"3h"`
	} `json:"snow,omitempty"`
	Sys struct {
		Pod string `json:"pod"`
	} `json:"sys"`
	DtTxt string `json:"dt_txt"`
}
