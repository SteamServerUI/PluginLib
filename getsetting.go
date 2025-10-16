package PluginLib

import (
	"fmt"
	"log"
)

// GetSetting retrieves the value of a specific setting by name from the /api/v2/settings endpoint.
func GetSetting(name string) (any, error) {
	if name == "" {
		return nil, fmt.Errorf("setting name must not be empty")
	}

	var settings map[string][]map[string]any
	err := Get("/api/v2/settings", &settings)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %v", err)
	}
	data, ok := settings["data"]
	if !ok {
		return nil, fmt.Errorf("no 'data' key in settings response")
	}
	for _, setting := range data {
		if settingName, ok := setting["name"].(string); ok && settingName == name {
			value, exists := setting["value"]
			if !exists {
				return nil, fmt.Errorf("setting '%s' has no 'value' field", name)
			}
			return value, nil
		}
	}

	return nil, fmt.Errorf("setting '%s' not found", name)
}

func GetAllSettings() (settings map[string][]map[string]any) {
	err := Get("/api/v2/settings", &settings)
	if err != nil {
		log.Fatalf("Failed to get settings: %v", err)
	}
	return settings
}
