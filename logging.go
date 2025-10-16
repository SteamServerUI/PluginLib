package PluginLib

import (
	"fmt"
	"log"
)

type LogResponse struct {
	Status string `json:"status"`
}

type LogRequest struct {
	Level      string `json:"level"`
	PluginName string `json:"pluginname"`
	Message    string `json:"message"`
}

func Log(message string, level ...string) error {
	if !configInitialized {
		log.Fatal("plugin configuration not initialized; call InitConfig first")
	}

	fmt.Println("Sending a log line to the server...")

	usedLevel := config.DefaultLevel

	// if level is provided, use the first one
	if len(level) > 0 {
		usedLevel = level[0]
	}
	if usedLevel == "" {
		return fmt.Errorf("log level must not be empty")
	}

	// Create LogRequest struct internally
	payload := LogRequest{
		Level:      usedLevel,
		PluginName: config.PluginName,
		Message:    message,
	}

	var logResponse LogResponse
	err := Post("/api/v2/plugins/log", payload, &logResponse)
	if err != nil {
		return fmt.Errorf("failed to log a line: %v", err)
	}

	fmt.Printf("Log sent successfully: %s\n", logResponse.Status)
	return nil
}
