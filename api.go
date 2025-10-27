package PluginLib

import (
	"log"
	"net/http"
)

var (
	pluginMux = http.NewServeMux()
)

type RegisterResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type RegisterRequest struct {
	PluginName string `json:"pluginname"`
}

func RegisterRoute(path string, handler http.HandlerFunc) {
	pluginMux.HandleFunc(path, handler)
	//Log("Registered route: "+path, "Debug")
}

func RegisterPluginAPI() {

	if !configInitialized {
		log.Fatal("plugin configuration not initialized; call InitConfig first")
	}

	// Send a POST to /api/v2/plugins/register
	var registerResponse RegisterResponse
	_, err := Post("/api/v2/plugins/register", RegisterRequest{PluginName: config.PluginName}, &registerResponse)
	if err != nil {
		Log("Failed to register pluginAPI: "+err.Error(), "Error")
		return
	}
	if registerResponse.Status != "success" {
		Log("Failed to register pluginAPI: "+registerResponse.Message, "Error")
		return
	}
}
