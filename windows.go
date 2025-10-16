//go:build windows
// +build windows

package plugininterface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/Microsoft/go-winio"
)

const pipeName = `\\.\pipe\ssui`

// Get sends a GET request to the specified SSUI endpoint and unmarshals the JSON response into the provided response interface.
func Get(endpoint string, response any) error {
	timeout := 5 * time.Second
	return sendRequest("GET", endpoint, nil, response, &timeout)
}

// Post sends a POST request to the specified SSUI endpoint with the given payload and unmarshals the JSON response.
func Post(endpoint string, payload any, response any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	timeout := 5 * time.Second
	return sendRequest("POST", endpoint, body, response, &timeout)
}

// sendRequest sends an HTTP request to the named pipe and parses the response.
func sendRequest(method, endpoint string, body []byte, response any, timeout *time.Duration) error {
	// Connect to the named pipe
	pipe, err := winio.DialPipe(pipeName, timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to named pipe %s: %w", pipeName, err)
	}
	defer pipe.Close()

	// Construct HTTP request
	host := "localhost"
	request := fmt.Sprintf("%s %s HTTP/1.1\r\nHost: %s\r\n", method, endpoint, host)
	if body != nil {
		request += "Content-Type: application/json\r\n"
		request += fmt.Sprintf("Content-Length: %d\r\n", len(body))
	}
	request += "Connection: close\r\n\r\n"
	if body != nil {
		request += string(body)
	}

	// Send request
	_, err = pipe.Write([]byte(request))
	if err != nil {
		return fmt.Errorf("failed to write to named pipe: %w", err)
	}

	// Read response
	respBytes, err := io.ReadAll(pipe)
	if err != nil {
		return fmt.Errorf("failed to read response from named pipe: %w", err)
	}

	// Parse HTTP response (simplified, assumes valid HTTP response)
	respStr := string(respBytes)
	headerEnd := bytes.Index([]byte(respStr), []byte("\r\n\r\n"))
	if headerEnd == -1 {
		return fmt.Errorf("invalid HTTP response: no header-body separator")
	}
	bodyStart := headerEnd + 4
	respBody := respBytes[bodyStart:]

	// Check status code (simplified, assumes first line is status)
	lines := bytes.SplitN(respBytes, []byte("\r\n"), 2)
	if len(lines) < 1 {
		return fmt.Errorf("invalid HTTP response: no status line")
	}
	if !bytes.Contains(lines[0], []byte("200 OK")) {
		return fmt.Errorf("unexpected status: %s", lines[0])
	}

	if response != nil {
		if err := json.Unmarshal(respBody, response); err != nil {
			return fmt.Errorf("failed to unmarshal JSON response: %w", err)
		}
	}
	return nil
}
