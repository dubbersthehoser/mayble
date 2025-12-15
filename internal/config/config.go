package config

import (
	"os"
	"errors"
	"encoding/json"
	"path/filepath"
)

// Config contains all configuration of the application.
//
type Config struct {
	ConfigDir  string `json:"config_dir"` 
	ConfigFile string `json:"config_file"`
	DBDriver   string `json:"db_driver"`
	DBFile     string `json:"db_file"`
}

// SetDBFile sets storage driver and saves it to Confgi.ConfigFile.
func (c *Config) SetDBDriver(s string) error {
	c.DBDriver = s
	return c.save()
}

// SetDBFile sets database file path and saves it to Config.ConfigFile.
func (c *Config) SetDBFile(s string) error {
	c.DBFile = s
	return c.save()
}

func (c *Config) save() error {
	path := c.ConfigFile

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.Encode(c)
	return nil
}

// Load config with root directory and from file name. Sets .ConfigFile with 
// joined root and file
func Load(root, file string) (*Config, error) {

	path := filepath.Join(root, file)

	fileIO, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Config{ConfigDir: root, ConfigFile: path}, nil
	}
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	decoder := json.NewDecoder(fileIO)
	
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
