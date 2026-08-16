package config

import (
	"errors"
	"os"
)

var (
	ErrRootNotExist error = errors.New("config: could not find root directory")
)

const (
	Filename string = "config.json"
)

func GetRootDir() (string, error) {
	var (
		hasHomeDir   bool
		hasConfigDir bool
	)
	userConfig, err := os.UserConfigDir()
	hasConfigDir = err == nil

	userHome, err := os.UserHomeDir()
	hasHomeDir = err == nil

	if hasConfigDir {
		return userConfig, nil
	}
	if hasHomeDir {
		return userHome, nil
	}
	return "", ErrRootNotExist
}
