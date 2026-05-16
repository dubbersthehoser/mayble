package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"fmt"
)

const Version string ="2.0.0"

type OldConfig struct {
	ConfigDir  string `json:"config_dir"`
	ConfigFile string `json:"config_file"`
	DBDriver   string `json:"db_driver"`
	DBFile     string `json:"db_file"`
}

// Config contains all configuration for the application.
type Config struct {
	Version    string `json:"version"`
	ConfigFile string `json:"config_file"`
	DBFile     string `json:"db_file"`
	Table struct {
		Headers map[string]struct{
			IsHidden bool
			Width    float32
		}
	}
	UI struct {
		WindowBody int
	}
}

func NewConfigWithDefaults(appName string) (*Config, error) {
	configFile, err := GetDefaultConfigFile(appName)
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		Version: Version,
		ConfigFile: configFile,
	}
	return cfg, nil
}

// Save to config file.
func (c *Config) Save() error {
	path := c.ConfigFile

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("config.save %w", err)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(c)
	return nil
}

func Load(path, appName string) (*Config, error) {

	fileIO, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		cfg, err := NewConfigWithDefaults(appName)
		if err != nil {
			return nil, err
		}
		cfg.ConfigFile = path
		return cfg, nil
	}
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	defer fileIO.Close()

	raw, err := io.ReadAll(fileIO)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	if isOld(raw) {
		cfg, err := NewConfigWithDefaults(appName)
		if err != nil {
			return nil, err
		}
		cfg.ConfigFile = path
		err = backup(path, cfg)
		return cfg, err
	}

	cfg := &Config{}

	err = json.Unmarshal(raw, cfg)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return cfg, err
}

func isOld(jsonBytes []byte) bool {
	oldCfg := &OldConfig{}
	err := json.Unmarshal(jsonBytes, oldCfg)
	if err != nil {
		return false
	}
	return oldCfg.ConfigDir != "" || oldCfg.DBDriver != ""
}

// backup create a backup of the old config.
func backup(path string, cfg *Config) error {
	to, err := os.Create(path + ".bak")
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	from, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	_, err = from.WriteTo(to)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	return nil
}
