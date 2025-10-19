package PluginLib

type ServerStatusResponse struct {
	Status bool   `json:"isRunning"`
	UUID   string `json:"uuid"`
}

type SettingsResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type ArgsRequest struct {
	Flag string `json:"flag"`
}

type ArgResponse struct {
	Value  string `json:"value,omitempty"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}
