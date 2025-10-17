package PluginLib

import (
	"log"
	"net"
	"net/http"
	"os"
)

func ExposeAPI() {

	if !configInitialized {
		log.Fatal("plugin configuration not initialized; call InitConfig first")
	}

	pluginSocketPath := "/tmp/ssui/" + config.PluginName + ".sock"

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveIndex)

	listener, err := net.Listen("unix", pluginSocketPath)
	if err != nil {
		Log("Error starting Unix socket server: " + err.Error())
		return
	}

	// Set socket permissions
	if err := os.Chmod(pluginSocketPath, 0666); err != nil {
		Log("Error setting socket permissions: " + err.Error())
	}

	// Create an HTTP server
	server := &http.Server{
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		Log("Unix socket server running at " + pluginSocketPath)
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			Log("Unix socket server error: " + err.Error())
		}
	}()

}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
