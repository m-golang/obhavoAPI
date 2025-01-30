package services

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func formatWeatherData(weatherData Weather) FormattedWeatherData {
	var formattedData FormattedWeatherData
	formattedData.Name = weatherData.Location.Name
	formattedData.Country = weatherData.Location.Country
	formattedData.Lat = weatherData.Location.Lat
	formattedData.Lon = weatherData.Location.Lon

	// Set Temperature Color
	formattedData.TempC = weatherData.Current.TempC
	formattedData.TempColor = getTempColor(formattedData.TempC)

	// Set Wind Color
	formattedData.WindKph = weatherData.Current.WindKph
	formattedData.WindColor = getWindColor(formattedData.WindKph)

	// Set Cloud Color
	formattedData.Cloud = weatherData.Current.Cloud
	formattedData.CloudColor = getCloudColor(formattedData.Cloud)

	return formattedData
}

func getTempColor(tempC float64) string {
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

func getWindColor(windKph float64) string {
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

func capitalizeFirstLetter(s string) string {
	caser := cases.Title(language.Und)
	return caser.String(s)
}
