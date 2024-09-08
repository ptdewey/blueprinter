package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type BlueprinterConfig struct {
	TemplateSources []string `json:"template-sources"`
}

func ParseConfig(configPath string) BlueprinterConfig {
	contents, err := os.ReadFile(configPath)
	if err != nil {
		return defaultConfig()
	}

	var cfg BlueprinterConfig
	if err = json.Unmarshal(contents, &cfg); err != nil {
		return defaultConfig()
	}

	return cfg
}

func defaultConfig() BlueprinterConfig {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Error getting user home directory: %s\n", err))
	}

	// TODO: maybe check if default directory exists
	// - other templates could be "~/templates", "~/Documents/Templates", or "~/Documents/templates"
	return BlueprinterConfig{
		TemplateSources: []string{
			filepath.Join(home, "Templates"),
		},
	}
}
