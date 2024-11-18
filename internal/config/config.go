package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

var (
	template_locations []string = []string{
		"Templates",
		"templates",
		"Documents/Templates",
		"Documents/templates",
	}

	config_files []string = []string{
		"blueprinter.toml",
		".blueprinter.toml",
		".blueprinterrc",
		".blueprinterrc.toml",
		"blueprinterrc.toml",
	}
)

type BlueprinterConfig struct {
	TemplateSources []string `toml:"templateSources"`
}

func ParseConfig() BlueprinterConfig {
	configPath, err := findConfigurationFile()
	if err != nil {
		fmt.Println("Configuration file not found.")
		return defaultConfig()
	}

	contents, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error reading contents of configuration file.")
		return defaultConfig()
	}

	var cfg BlueprinterConfig
	if err = toml.Unmarshal(contents, &cfg); err != nil {
		fmt.Println("Error reading configuration file.")
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
	configPath := "blueprinter.toml"
	if checkFileExists(configPath) {
		return configPath, nil
	}

	// check within git project next if config file not found.
	// if git root is not found, skip this step
	gr, err := findGitRoot()
	if err == nil {
		for _, f := range config_files {
			configPath := filepath.Join(gr, f)
			if checkFileExists(configPath) {
				return configPath, nil
			}
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
		// NOTE: error handling here may be necessary
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
