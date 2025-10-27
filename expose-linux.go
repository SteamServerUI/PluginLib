//go:build linux
// +build linux

package PluginLib

import (
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

func ExposeAPI(wg *sync.WaitGroup) {

	if !configInitialized {
		log.Fatal("plugin configuration not initialized; call InitConfig first")
	}

	pluginSocketPath := "./SSUI/plugins/sockets/" + config.PluginName + ".sock"

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
	wg.Add(1)
	wg.Go(func() {
		defer wg.Done()
		Log("Unix socket server running at " + pluginSocketPath)
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			Log("Unix socket server error: " + err.Error())
		}
	})

}
