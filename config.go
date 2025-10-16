package PluginLib

import "fmt"

// Global configuration (must be initialized by plugins).
var config PluginConfig
var configInitialized bool

type PluginConfig struct {
	PluginName   string
	DefaultLevel string
}

func InitConfig(pluginName, defaultLevel string) error {
	if pluginName == "" || defaultLevel == "" {
		return fmt.Errorf("pluginName and defaultLevel must not be empty")
	}
	config = PluginConfig{
		PluginName:   pluginName,
		DefaultLevel: defaultLevel,
	}
	configInitialized = true
	return nil
}
