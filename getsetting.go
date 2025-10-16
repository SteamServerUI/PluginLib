package PluginLib

import (
	"log"
)

func GetSetting() (setting map[string][]map[string]any) {
	var settings map[string][]map[string]any
	err := Get("/api/v2/settings", &settings)
	if err != nil {
		log.Fatalf("Failed to get settings: %v", err)
	}
	return settings
}
