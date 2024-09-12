package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	template_locations []string = []string{
		"Templates",
		"templates",
		"Documents/Templates",
		"Documents/templates",
	}

	config_files []string = []string{
		"blueprinter.json",
		".blueprinter.json",
		".blueprinterrc",
		".blueprinterrc.json",
		"blueprinterrc.json",
	}
)

type BlueprinterConfig struct {
	TemplateSources []string `json:"template-sources"`
}

func ParseConfig() BlueprinterConfig {
	configPath, err := findConfigurationFile()
	if err != nil {
		return defaultConfig()
	}

	contents, err := os.ReadFile(configPath)
	if err != nil {
		return defaultConfig()
	}

	var cfg BlueprinterConfig
	if err = json.Unmarshal(contents, &cfg); err != nil {
		return defaultConfig()
	}

	if len(cfg.TemplateSources) == 0 {
		return defaultConfig()
	}

	for i, templateSource := range cfg.TemplateSources {
		if strings.HasPrefix(templateSource, "~") {
			home, err := os.UserHomeDir()
			if err != nil {
				panic(fmt.Sprintf("Error getting user home directory: %s\n", err))
			}

			cfg.TemplateSources[i] = strings.Replace(templateSource, "~", home, 1)
		}
	}

	return cfg
}

func findConfigurationFile() (string, error) {
	configPath := "blueprinter.json"
	if checkFileExists(configPath) {
		return configPath, nil
	}

	gr, err := findGitRoot()
	if err != nil {
		return "", err
	}
	for _, f := range config_files {
		configPath := filepath.Join(gr, f)
		if checkFileExists(configPath) {
			return configPath, nil
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return "", err
	}
	for _, f := range config_files {
		configPath := filepath.Join(home, f)
		if checkFileExists(configPath) {
			return configPath, nil
		}
	}

	return "", errors.New("No configuration file found.")
}

func findGitRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")

	out, err := cmd.Output()
	if err != nil {
		// not in git repository
		return "", err
	}

	gitRoot := strings.TrimSpace(string(out))

	return gitRoot, nil
}

func checkFileExists(configPath string) bool {
	if info, err := os.Stat(configPath); err == nil {
		return !info.IsDir()
	}
	return false
}

func defaultConfig() BlueprinterConfig {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Error getting user home directory: %s\n", err))
	}

	for _, template_location := range template_locations {
		dir := filepath.Join(home, template_location)
		if _, err := os.Stat(dir); err != nil {
			continue
		}

		return BlueprinterConfig{
			TemplateSources: []string{dir},
		}
	}

	panic("Error: No Bluprinter config directory specified and no fallback was found.")
}
