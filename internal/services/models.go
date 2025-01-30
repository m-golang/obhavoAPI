package services

// Weather holds the location and current weather data.
type Weather struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}

// Location holds the essential location details.
type Location struct {
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"` // Using float64 for better precision
	Lon     float64 `json:"lon"` // Using float64 for better precision
}

// Current holds the essential weather details.
type Current struct {
	TempC   float64 `json:"temp_c"`   // Temperature in Celsius
	WindKph float64 `json:"wind_kph"` // Wind speed in kilometers per hour
	Cloud   int     `json:"cloud"`    // Cloud cover percentage
}

type FormattedWeatherData struct {
	Name       string  `json:"name"`
	Country    string  `json:"country"`
	Lat        float64 `json:"lat"`    // Using float64 for better precision
	Lon        float64 `json:"lon"`    // Using float64 for better precision
	TempC      float64 `json:"temp_c"` // Temperature in Celsius
	TempColor  string  `json:"temp_color"`
	WindKph    float64 `json:"wind_kph"` // Wind speed in kilometers per hour
	WindColor  string  `json:"wind_color"`
	Cloud      int     `json:"cloud"` // Cloud cover percentage
	CloudColor string  `json:"cloud_color"`
}
