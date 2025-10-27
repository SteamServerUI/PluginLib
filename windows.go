//go:build windows
// +build windows

package PluginLib

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Microsoft/go-winio"
)

var sSUIPipePath string

func getSSUIPipePath() string {
	if sSUIPipePath != "" {
		return sSUIPipePath
	}
	// check the ./SSUI/plugins/sockets/pipename.identifier file for the pipe name. It will only contain the name, nothing else.
	sSUIPipePathFile, err := os.Open("./SSUI/plugins/sockets/pipename.identifier")
	if err != nil {
		fmt.Println("Error opening pipename.identifier file, I have to go...:", err)
		os.Exit(1)
	}
	defer sSUIPipePathFile.Close()

	scanner := bufio.NewScanner(sSUIPipePathFile)
	for scanner.Scan() {
		sSUIPipePath = scanner.Text()
	}
	if sSUIPipePath == "" {
		fmt.Println("Error reading pipename.identifier file, I have to go...:", err)
		os.Exit(1)
	}
	sSUIPipePath = sSUIPipePath + "ssui"
	fmt.Println("SSUIPipePath:", sSUIPipePath)
	return sSUIPipePath
}

// Get sends a GET request to the specified SSUI endpoint and unmarshals the JSON response into the provided response interface.
func Get(endpoint string, response any) (any, error) {
	pipeName := getSSUIPipePath()
	transport := &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return winio.DialPipe(pipeName, nil)
		},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	resp, err := client.Get("http://localhost" + endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request to %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if response != nil {
		if err := json.Unmarshal(body, response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
		}
	}
	return response, nil
}

// Post sends a POST request to the specified SSUI endpoint with the given payload and unmarshals the JSON response.
func Post(endpoint string, payload any, response any) (any, error) {
	pipeName := getSSUIPipePath()
	transport := &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return winio.DialPipe(pipeName, nil)
		},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := client.Post("http://localhost"+endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request to %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if response != nil {
		if err := json.Unmarshal(respBody, response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
		}
	}
	return response, nil
}
