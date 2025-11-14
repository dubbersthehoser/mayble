package config

import (
	"os"
	"io"
	"errors"
	"encoding/json"
)

type Properties map[string]string
type Sections map[string]Properties

type Config struct {
	Sections Sections
	Path string
}
func NewConfig() *Config {
	config := &Config{
		Sections: Sections{},
	}
	return config
}

// Load config file from file path. If file dose not exist returns a new config with given path.
func Load(path string) (*Config, error){
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		config := NewConfig()
		config.Path = path
		return NewConfig(), nil
	}else if err != nil {
		return nil, err
	}
	defer file.Close()
	s, err :=  Read(file)
	if err != nil {
		return nil, err
	}
	config := &Config{
		Sections: s,
	}
	return config, nil
}

// Save the config to path, if path == "" then uses config.Path to save to. 
func Save(config *Config, path string) error {
	if path == "" {
		path = config.Path
	}
	file, err := os.Create(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	return Write(config.Sections, file)

}

func Read(r io.Reader) (Sections, error) {
	s := Sections{}
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&s); err != nil {
		return nil, err
	}
	return s, nil
}

func Write(s Sections, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("  ", "")
	if err := encoder.Encode(s); err != nil {
		return err
	}
	return nil
}


