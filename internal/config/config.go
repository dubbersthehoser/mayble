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

type Header struct {
	IsHidden bool    `json:"is_hidden"`
	Width    float32 `json:"width"`
}

// Config contains all configuration for the application.
type Config struct{
	Version    string `json:"version"`
	ConfigFile string `json:"config_file"`
	DBFile     string `json:"db_file"`
	UI struct{
		OpenBody int              `json:"open_body"`
		Headers map[string]Header `json:"headers"`
		TableSortBy   string      `json:"table_sort_by"`
		TableAscending bool       `json:"table_ascending"`
	} `json:"ui"`
}

func NewConfigWithDefaults(appName string) (*Config, error) {
	configFile, err := GetDefaultConfigFile(appName)
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		Version: Version,
		ConfigFile: configFile,
		UI: struct{
			OpenBody int              `json:"open_body"`
			Headers map[string]Header `json:"headers"`
			TableSortBy string        `json:"table_sort_by"`
			TableAscending bool       `json:"table_ascending"`
		}{
			Headers: make(map[string]Header),
		},
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
		err = backup(path)
		return cfg, err
	}

	cfg := &Config{}

	err = json.Unmarshal(raw, cfg)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return cfg, err
}

// isOld check wether json string is in the format of the old config.
func isOld(jsonBytes []byte) bool {
	oldCfg := &OldConfig{}
	err := json.Unmarshal(jsonBytes, oldCfg)
	if err != nil {
		return false
	}
	// DBDriver was removed for v2 config.
	return oldCfg.DBDriver != ""
}

// backup create a backup of the old config.
func backup(path string) error {
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
