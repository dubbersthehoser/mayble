package config

import (
	"os"
	"io"
	"io/fs"
	"errors"
	"encoding/json"
)

const FileName string = "config.json"

type Setting struct {
	fileName string
	root     os.Root
}

type Option func(*Setting)

func WithFileName(name string) Option {
	return func(s *Setting) {
		s.fileName = name
	}
}

func WithDirectory(dir string) (Option, error) {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, err
	}
	return WithRoot(root)
}

func WithRoot(root os.Root) Option {
	return func(s *Setting) {
		s.root = root
	}
}

type ConfigStore struct {
	root     os.Root
	fileName string
}

func NewConfigStore(op []Options) {
	for _, o := range op {
		o(s)
	}
}

type Config struct {
	DBPath     string `json:"db_path"`
	ExportPath string `json:"export_path"`
}

func OpenWithRoot(root *os.Root) (*Config, error) {
	file, err := root.Open(FileName)
	if errors.Is(err, fs.ErrNotExist) {
		return &Config{}, nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()
	return Read(file)
}

func SaveWithRoot(root *os.Root, config *Config) error {
	file, err := root.Open(FileName)
	if errors.Is(err, fs.ErrNotExist) {
		file, err = root.Create(FileName)
		if err != nil{
			return err
		}
	} else if err != nil {
		return err
	}
	defer file.Close()
	return Write(config, file)
}

func Read(r io.Reader) (*Config, error) {
	config := &Config{}
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

func Write(config *Config, w io.Writer) error {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(config); err != nil {
		return err
	}
	return nil
}

