//go:build linux
// +build linux

package plugininterface

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

const socketPath = "/tmp/ssui/ssui.sock"

// Get sends a GET request to the specified SSUI endpoint and unmarshals the JSON response into the provided response interface.
func Get(endpoint string, response any) error {
	transport := &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", socketPath)
		},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	resp, err := client.Get("http://localhost" + endpoint)
	if err != nil {
		return fmt.Errorf("failed to send GET request to %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if response != nil {
		if err := json.Unmarshal(body, response); err != nil {
			return fmt.Errorf("failed to unmarshal JSON response: %w", err)
		}
	}
	return nil
}

// Post sends a POST request to the specified SSUI endpoint with the given payload and unmarshals the JSON response.
func Post(endpoint string, payload any, response any) error {
	transport := &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", socketPath)
		},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := client.Post("http://localhost"+endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to send POST request to %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if response != nil {
		if err := json.Unmarshal(respBody, response); err != nil {
			return fmt.Errorf("failed to unmarshal JSON response: %w", err)
		}
	}
	return nil
}
