package utils

import "os"

func CheckFileExists(configPath string) bool {
	if info, err := os.Stat(configPath); err == nil {
		return !info.IsDir()
	}
	return false
}
