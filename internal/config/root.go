package config

import (
	"os"
)

type ConfigRoot struct {
	root *os.Root
	filename string
}

func OpenConfigRoot(dir string) (*ConfigRoot, error) {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, err
	}
	configRoot := &ConfigRoot{
		root: root,
		filename: FileName,

	}
	return configRoot, nil
}

func (c *ConfigRoot) SetFileName(name string) {
	c.filename = name
}

func (c *ConfigRoot) Open() (*Config, error) {
	return OpenWithRoot(c.root)
}

func (c *ConfigRoot) Save(config *Config) error {
	return SaveWithRoot(c.root, config)
}

