package PluginLib

import (
	"fmt"
)

var serverStatusResponse ServerStatusResponse

func GetServerStatus() (ServerStatusResponse, error) {
	if err := Get("/api/v2/server/status", &serverStatusResponse); err != nil {
		fmt.Println("Error:", err)
		return serverStatusResponse, err
	}
	return serverStatusResponse, nil
}
