package tools

import "fmt"

// Example function implementation
func GetWeather(location, format string) string {
	return fmt.Sprintf("Weather in %s is 22° %s", location, format)
}
