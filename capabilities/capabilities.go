package capabilities

import (
	"encoding/json"
	"fmt"

	t "github.com/adam-fraga/eris/capabilities/tools"
)

// Capability represents a function the LLM can call
type Capability struct {
	Name        string
	Description string
	Func        interface{}
}

// Map of all capabilities
var Capabilities = map[string]Capability{
	"get_current_weather": {
		Name:        "get_current_weather",
		Description: "Get the current weather for a location",
		Func:        t.GetWeather,
	},
	"create_file": {
		Name:        "create_file",
		Description: "Create a file with specified content",
		Func:        t.CreateFile,
	},
}

func CallCapability(name string, argsJson []byte) (string, error) {
	cap, ok := Capabilities[name]
	if !ok {
		return "", fmt.Errorf("capability not found: %s", name)
	}

	switch name {
	case "get_current_weather":
		var args struct {
			Location string `json:"location"`
			Format   string `json:"format"`
		}
		if err := json.Unmarshal(argsJson, &args); err != nil {
			return "", err
		}
		f := cap.Func.(func(string, string) string)
		return f(args.Location, args.Format), nil

	case "create_file":
		var args struct {
			Filename string `json:"filename"`
			Content  string `json:"content"`
		}
		if err := json.Unmarshal(argsJson, &args); err != nil {
			return "", err
		}
		f := cap.Func.(func(string, string) string)
		return f(args.Filename, args.Content), nil

	default:
		return "", fmt.Errorf("unsupported function signature for %s", name)
	}
}

// JSON structure for the LLM
type ChatCapability struct {
	Type     string        `json:"type"` // always "function"
	Function FunctionDescr `json:"function"`
}

type FunctionDescr struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Parameters  map[string]CapabilityProp `json:"parameters"`
}

type CapabilityProp struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

// Convert capabilities to []interface{} for ChatRequest
func ToChatCapabilitiesInterface() []interface{} {
	caps := []interface{}{}
	for _, cap := range Capabilities {
		var params map[string]CapabilityProp

		// Dynamically define params based on the capability name
		switch cap.Name {
		case "get_current_weather":
			params = map[string]CapabilityProp{
				"location": {Type: "string", Description: "Location to get the weather for"},
				"format":   {Type: "string", Description: "celsius or fahrenheit", Enum: []string{"celsius", "fahrenheit"}},
			}
		case "create_file":
			params = map[string]CapabilityProp{
				"filename": {Type: "string", Description: "Name of the file to create"},
				"content":  {Type: "string", Description: "Content to write into the file"},
			}
		default:
			params = map[string]CapabilityProp{}
		}

		caps = append(caps, ChatCapability{
			Type: "function",
			Function: FunctionDescr{
				Name:        cap.Name,
				Description: cap.Description,
				Parameters:  params,
			},
		})
	}
	return caps
}
