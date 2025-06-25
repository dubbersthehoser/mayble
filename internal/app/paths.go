package app

import (
	"os"
	"path/filepath"
)

type Paths struct {
	ConfigPath  string
	LogsPath    string
	DBPath      string
	StoragePath string
}

func NewPaths() (*Paths, error) {
	
	config := &Paths{}

	root, err := ConfigDir()
	if err != nil {
		return nil, err
	}
	config.StoragePath = root
	config.ConfigPath  = filepath.Join(root, ConfigFile)
	config.DBPath      = filepath.Join(root, DBFile)
	config.LogsPath    = filepath.Join(root, LogsFile)
	return config, nil
}

func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(home, "." + AppName), nil
	} 
	return "", err
}

func (c *Paths) InitStoreage() error {
	if _, err := os.Stat(c.StoragePath); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(c.StoragePath, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return nil
}
