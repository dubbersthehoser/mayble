package config

import (
	"os"
	"fmt"
	"errors"
	"path/filepath"
)

const (
	filename string = "config.json"
)

func GetDefaultConfigFile(appName string) (string, error) {
	path, found := findConfigDir(appName)
	if path == "" {
		return "", errors.New("config: could not find config directory")
	}
	if !found {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return "", fmt.Errorf("config: %w", err)
		}
	}
	path = filepath.Join(path, filename)
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

