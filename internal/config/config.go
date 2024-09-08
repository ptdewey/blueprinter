package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type BlueprinterConfig struct {
	TemplateSources []string `json:"template-sources"`
}

func ParseConfig(filepath string) (BlueprinterConfig, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return BlueprinterConfig{}, err
	}

	var cfg BlueprinterConfig
	if err = json.Unmarshal(contents, &cfg); err != nil {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting user home directory: ", err)
			return BlueprinterConfig{}, err
		}

		// TODO: maybe check if default directories exist
		return BlueprinterConfig{
			TemplateSources: []string{home + "/Templates"},
		}, nil
	}

	// TODO: check if '~/' was used in template-sources and replace with '$HOME'

	return cfg, nil
}
