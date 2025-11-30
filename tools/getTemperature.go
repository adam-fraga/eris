package tools

func GetTemperatureTool() Tool {
	return Tool{
		Type: "function",
		Function: struct {
			Name        string         `json:"name"`
			Description string         `json:"description"`
			Parameters  map[string]any `json:"parameters"`
		}{
			Name:        "get_temperature",
			Description: "Return the current temperature of a city",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"city": map[string]any{
						"type":        "string",
						"description": "The city name",
					},
				},
				"required": []string{"city"},
			},
		},
	}
}

func GetTemperature(city string) string {
	if city == "New York" {
		return "22°C"
	}
	if city == "Toronto" {
		return "-25°C"
	}
	return "unknown"
}
