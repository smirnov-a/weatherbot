package weather

// WeatherData structure with current weather and forecast
type WeatherData struct {
	CurrentData  *CurrentData
	ForecastData *ForecastData
}

type CurrentData struct {
	City    string
	Weather float64
}

type ForecastData struct {
	Days    float64
	Offset  int64
	Sunrise string
	Sunset  string
	Rows    []Row
}

type Row struct {
	Timestamp     string
	Temperature   float64
	FeelsLike     float64
	Pressure      float64
	Humidity      int
	Weather       string
	Clouds        int
	Visibility    int
	Pop           string
	Precipitation float64
	Wind          Wind
}

type Wind struct {
	Speed float64
	Deg   int
	Gust  float64
}

type LocalName struct {
	Locale string
	Name   string
}

// CityInfo structure with latitude/longitude
type CityInfo struct {
	Name      string
	Latitude  float64
	Longitude float64
	HasCoords bool
}
