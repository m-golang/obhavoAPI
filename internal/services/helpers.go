package services

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// formatWeatherData formats the raw weather data into a user-friendly structure
// with additional properties like color codes for temperature, wind, and cloud conditions.
func formatWeatherData(weatherData Weather) FormattedWeatherData {
	// Initialize the formatted weather data structure.
	var formattedData FormattedWeatherData

	// Extract location details from the weather data.
	formattedData.Name = weatherData.Location.Name
	formattedData.Country = weatherData.Location.Country
	formattedData.Lat = weatherData.Location.Lat
	formattedData.Lon = weatherData.Location.Lon

	// Set temperature and corresponding color code based on the temperature.
	formattedData.TempC = weatherData.Current.TempC
	formattedData.TempColor = getTempColor(formattedData.TempC)

	// Set wind speed and corresponding color code based on the wind speed.
	formattedData.WindKph = weatherData.Current.WindKph
	formattedData.WindColor = getWindColor(formattedData.WindKph)

	// Set cloud coverage percentage and corresponding color code based on the cloud coverage.
	formattedData.Cloud = weatherData.Current.Cloud
	formattedData.CloudColor = getCloudColor(formattedData.Cloud)

	// Return the fully formatted weather data.
	return formattedData
}

// getTempColor determines the color associated with the temperature.
// The color changes based on the temperature value to visually represent different temperature ranges.
func getTempColor(tempC float64) string {
	// Define color ranges for different temperature values (in Celsius).
	if tempC < -20 {
		return "#003366" // Deep Blue
	} else if tempC >= -20 && tempC < -10 {
		return "#4A90E2" // Ice Blue
	} else if tempC >= -10 && tempC < 0 {
		return "#B3DFFD" // Light Blue
	} else if tempC >= 0 && tempC < 10 {
		return "#E6F7FF" // Pale Grayish Blue
	} else if tempC >= 10 && tempC < 20 {
		return "#D1F2D3" // Light Green
	} else if tempC >= 20 && tempC < 30 {
		return "#FFFACD" // Soft Yellow
	} else if tempC >= 30 && tempC < 40 {
		return "#FFCC80" // Light Orange
	} else if tempC >= 40 && tempC < 50 {
		return "#FF7043" // Deep Orange
	} else if tempC >= 50 {
		return "#D32F2F" // Bright Red
	}

	return "#FFFFFF" // Default color if no condition matches
}

// getWindColor determines the color associated with wind speed.
// The color changes based on the wind speed to visually represent different wind intensities.
func getWindColor(windKph float64) string {
	// Define color ranges for different wind speeds (in kilometers per hour).
	if windKph >= 0 && windKph < 10 {
		return "#E0F7FA" // Light Cyan
	} else if windKph >= 10 && windKph < 20 {
		return "#B2EBF2" // Pale Blue
	} else if windKph >= 20 && windKph < 40 {
		return "#4DD0E1" // Soft Teal
	} else if windKph >= 40 && windKph < 60 {
		return "#0288D1" // Bright Blue
	} else if windKph >= 60 {
		return "#01579B" // Deep Navy Blue
	}
	return "#FFFFFF" // Default color if no condition matches
}

// getCloudColor determines the color associated with cloud coverage.
// The color changes based on the cloud coverage percentage to visually represent different cloud conditions.
func getCloudColor(cloud int) string {
	if cloud >= 0 && cloud < 10 {
		return "#FFF9C4" // Light Yellow
	} else if cloud >= 10 && cloud < 30 {
		return "#FFF176" // Soft Yellow
	} else if cloud >= 30 && cloud < 60 {
		return "#E0E0E0" // Light Gray
	} else if cloud >= 60 && cloud < 90 {
		return "#9E9E9E" // Gray
	} else if cloud >= 90 && cloud <= 100 {
		return "#616161" // Dark Gray
	}
	return "#FFFFFF" // Default color if no condition matches
}

// capitalizeFirstLetter capitalizes the first letter of a string.
// It is useful for formatting location names or other textual data that should follow proper casing.
func capitalizeFirstLetter(s string) string {
	// Use the Title casing rules for capitalization.
	caser := cases.Title(language.Und)
	return caser.String(s)
}
