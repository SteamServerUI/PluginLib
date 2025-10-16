package PluginLib

type ServerStatusResponse struct {
	Status bool   `json:"isRunning"`
	UUID   string `json:"uuid"`
}

type SettingsResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}
