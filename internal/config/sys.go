package config

import (
	"os"
	"path/filepath"

	"github.com/dubbersthehoser/mayble/internal/status"
)

// HasUserConfigDir returns the user config directory string and true if found.
// If not returns empty string and false.
func HasUserConfigDir() (string, bool) {
	userConfig, err := os.UserConfigDir()
	if err != nil {
		return "", false
	}
	return userConfig, true
}

// HasUserHomeDir returns user home directory path and true if found. If not
// returns empty string and false.
func HasUserHomeDir() (string, bool) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}
	return userHome, true
}

// FindConfigDir with given application name returns the full path of exsisting
// folder if exists starting from user config then user home. If not found
// returns the approprite path to use if not returns empty string.
func FindConfigDir(appName string) (string, bool) {
	userConfig, isConfig := HasUserConfigDir()
	userHome, isHome := HasUserHomeDir()
	if isConfig {
		path := filepath.Join(userConfig, appName)
		_, err := os.Stat(path)
		if err == nil {
			return path, true
		}
	}
	if isHome {
		path := filepath.Join(userHome, "."+appName)
		_, err := os.Stat(path)
		if err == nil {
			return path, true
		}
	}

	if isConfig {
		path := filepath.Join(userConfig, appName)
		return path, false
	}
	if isHome {
		path := filepath.Join(userHome, "."+appName)
		return path, false
	}
	return "", false
}

// GetDefaultDir with given app name returns approprite path and creates
// directory if it dose not exsits.
func GetDefaultDir(appName string) (string, error) {
	const op status.Op = "config.GetDefaultDir"

	path, found := FindConfigDir(appName)
	if path == "" {
		return "", status.E(op, status.LevelError, "could not find config directory")
	}
	if !found {
		err := os.Mkdir(path, 0744)
		if err != nil {
			return "", status.E(op, status.LevelError, status.FailedToCreate, err)
		}
	}
	return path, nil
}
