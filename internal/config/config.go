package config

import (
	"os"
	"io"
	"encoding/json"
	"path/filepath"
)

type Config struct {
	ConfigDir  string `json:"config_dir"` 
	ConfigFile string `json:"config_file"`
	DBDriver   string `json:"db_driver"`
	DBFile     string `json:"db_file"`
}

func (c *Config) SetDBDirver(s string) error {
	c.DBDriver = s
	return c.save()
}

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
}

func Load(dir, file, string) (*Config, error) {

	path := filepath.Join(dir, file)

	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExsit) {
		return &Config{ConfigDir: dir, ConfigFile: path}, nil
	}
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	decoder := json.NewDecoder(file)
	
	err := decoder.Decode(cfg)
	if err != nil {
		return err
	}

	return cfg, nil
}
