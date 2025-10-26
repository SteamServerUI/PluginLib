package PluginLib

import (
	"log"
	"net"
	"net/http"
	"os"
	"sync"
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

func ExposeAPI(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	if !configInitialized {
		log.Fatal("plugin configuration not initialized; call InitConfig first")
	}

	pluginSocketPath := "./SSUI/plugins/sockets" + config.PluginName + ".sock"

	if err := os.RemoveAll(pluginSocketPath); err != nil {
		Log("Error removing existing socket: " + err.Error())
	}

	listener, err := net.Listen("unix", pluginSocketPath)
	if err != nil {
		Log("Error starting Unix socket server: " + err.Error())
		return
	}

	// Set socket permissions
	if err := os.Chmod(pluginSocketPath, 0600); err != nil {
		Log("Error setting socket permissions: " + err.Error())
	}

	// Create an HTTP server
	server := &http.Server{
		Handler: pluginMux,
	}

	// Start server in a goroutine
	wg.Go(func() {
		Log("Unix socket server running at " + pluginSocketPath)
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			Log("Unix socket server error: " + err.Error())
		}
	})

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
