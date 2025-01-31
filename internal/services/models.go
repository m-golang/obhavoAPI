package services

// Weather holds the location and current weather data.
// It represents the full weather report for a specific location.
type Weather struct {
	Location Location `json:"location"` // Location contains geographical details like name, country, and coordinates.
	Current  Current  `json:"current"`  // Current contains the current weather data like temperature, wind speed, and cloud cover.

}

// Location holds the essential location details such as name, country, and coordinates.
// It is used to represent the geographical information for the weather data.
type Location struct {
	Name    string  `json:"name"`    // Name represents the name of the location (e.g., city, town, etc.).
	Country string  `json:"country"` // Country represents the country of the location.
	Lat     float64 `json:"lat"`     // Using float64 for better precision.
	Lon     float64 `json:"lon"`     // Using float64 for better precision.
}

// Current holds the essential weather details for the current conditions.
// It represents data such as temperature, wind speed, and cloud coverage.
type Current struct {
	TempC   float64 `json:"temp_c"`   // Temperature in Celsius.
	WindKph float64 `json:"wind_kph"` // Wind speed in kilometers per hour.
	Cloud   int     `json:"cloud"`    // Cloud cover percentage.
}

// FormattedWeatherData holds the weather data after it has been processed and formatted,
// including additional properties such as color codes for visual representation.
type FormattedWeatherData struct {
	Name       string  `json:"name"`        // Name represents the name of the location (e.g., city, town, etc.).
	Country    string  `json:"country"`     // Country represents the country of the location.
	Lat        float64 `json:"lat"`         // Using float64 for better precision.
	Lon        float64 `json:"lon"`         // Using float64 for better precision.
	TempC      float64 `json:"temp_c"`      // Temperature in Celsius.
	TempColor  string  `json:"temp_color"`  // TempColor represents the color code associated with the current temperature.
	WindKph    float64 `json:"wind_kph"`    // Wind speed in kilometers per hour.
	WindColor  string  `json:"wind_color"`  // WindColor represents the color code associated with the wind speed.
	Cloud      int     `json:"cloud"`       // Cloud cover percentage.
	CloudColor string  `json:"cloud_color"` // This can be used for visual representation of different cloud cover levels.
}
