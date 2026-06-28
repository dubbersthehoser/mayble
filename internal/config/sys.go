package config

import (
	"os"
	"errors"
	"path/filepath"
)

const (
	filename string = "config.json"
)

func GetDefaultConfigFile(appName string) (string, error) {
	dir, ok := findConfigDir(appName)
	if !ok {
		return "", errors.New("config: could not find default config directory")

	}
	path := filepath.Join(dir, filename)
	return path, nil
}

func findConfigDir(appName string) (string, bool) {
	var (
		hasHomeDir bool
		hasConfigDir bool
	)
	userConfig, err := os.UserConfigDir()
	hasConfigDir = err == nil

	userHome, err := os.UserHomeDir()
	hasHomeDir = err == nil

	if hasConfigDir {
		path := filepath.Join(userConfig, appName)
		_, err := os.Stat(path)
		if err == nil {
			return path, true
		}
	}

	if hasHomeDir {
		path := filepath.Join(userHome, "."+appName)
		_, err := os.Stat(path)
		if err == nil {
			return path, true
		}
	}
	return "", false
}

