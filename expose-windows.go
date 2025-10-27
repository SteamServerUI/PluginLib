//go:build windows
// +build windows

package PluginLib

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/Microsoft/go-winio"
)

var PluginPipePath string

func getPluginPipePath() string {
	if PluginPipePath != "" {
		return PluginPipePath
	}
	// check the ./SSUI/plugins/sockets/pipename.identifier file for the pipe name. It will only contain the name, nothing else.
	PluginPipePathFile, err := os.Open("./SSUI/plugins/sockets/pipename.identifier")
	if err != nil {
		fmt.Println("Error opening pipename.identifier file, I have to go...:", err)
		os.Exit(1)
	}
	defer PluginPipePathFile.Close()

	scanner := bufio.NewScanner(PluginPipePathFile)
	for scanner.Scan() {
		PluginPipePath = scanner.Text()
	}
	if PluginPipePath == "" {
		fmt.Println("Error reading pipename.identifier file, I have to go...:", err)
		os.Exit(1)
	}
	//returns something like \\.\pipe\ssui-1234567890\
	return PluginPipePath + config.PluginName
}

func ExposeAPI(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	pluginPipePath := getPluginPipePath()

	if !configInitialized {
		log.Fatal("plugin configuration not initialized; call InitConfig first")
	}

	// Remove existing pipe if it exists
	if err := os.RemoveAll(pluginPipePath); err != nil {
		Log("Error removing existing pipe: " + err.Error())
	}

	// Create a named pipe listener using go-winio
	listener, err := winio.ListenPipe(pluginPipePath, nil)
	if err != nil {
		Log("Error starting named pipe server: " + err.Error())
		return
	}

	// Named pipes on Windows have permissions set via Security Descriptor
	// Default nil config in ListenPipe provides reasonable security (creator-only access)

	// Create an HTTP server
	server := &http.Server{
		Handler: pluginMux,
	}

	// Start server in a goroutine
	go func() {
		Log("Named pipe server running at " + pluginPipePath)
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			Log("Named pipe server error: " + err.Error())
		}
	}()

	// Ensure listener is closed when the function exits
	defer listener.Close()
}
