package PluginLib

import (
	"fmt"
)

var serverStatusResponse ServerStatusResponse
var getArgResponse ArgResponse

func GetServerStatus() (ServerStatusResponse, error) {
	if _, err := Get("/api/v2/server/status", &serverStatusResponse); err != nil {
		fmt.Println("Error:", err)
		return serverStatusResponse, err
	}
	return serverStatusResponse, nil
}

func GetSingleArgFromRunfile(flag string) (string, error) {
	payload := ArgsRequest{
		Flag: flag,
	}

	_, err := Post("/api/v2/runfile/args/getarg", payload, &getArgResponse)
	if err != nil {
		return "", fmt.Errorf("failed to get arg from runfile api: %v", err)
	}
	if getArgResponse.Status != "success" {
		return "", fmt.Errorf("failed to get arg from runfile api: %s", getArgResponse.Error)
	}
	return getArgResponse.Value, nil

}
