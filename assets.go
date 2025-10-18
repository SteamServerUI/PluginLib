package PluginLib

import (
	"embed"
	"io/fs"
)

// AssetManager provides access to embedded assets.
type AssetManager struct {
	fs fs.FS
}

// RegisterAssets registers an embedded filesystem for the plugin.
func RegisterAssets(embeddedFS *embed.FS) *AssetManager {
	return &AssetManager{fs: embeddedFS}
}

// GetAsset reads an embedded asset as a byte slice.
func (am *AssetManager) GetAsset(path string) ([]byte, error) {
	return fs.ReadFile(am.fs, path)
}

// GetAssetString reads an embedded asset as a string.
func (am *AssetManager) GetAssetString(path string) (string, error) {
	data, err := am.GetAsset(path)
	return string(data), err
}

// OpenAsset returns an fs.File for streaming or advanced use cases.
func (am *AssetManager) OpenAsset(path string) (fs.File, error) {
	return am.fs.Open(path)
}
