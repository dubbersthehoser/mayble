package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/dubbersthehoser/mayble/internal/status"
)

// Table config for table view.
type Table struct {
	ColumnsHidden []string           `json:"hidden_columns"`
	ColumnWidths  map[string]float32 `json:"column_width"`
	cfg           *Config
}

// SetColumnWidth for column lable for size s.
func (t *Table) SetColumnWidth(label string, s float32) {
	w := t.ColumnWidths
	if w == nil {
		w = make(map[string]float32)
	}
	w[label] = s
}

// GetColumnWidth from named label.
func (t *Table) GetColumnWidth(label string) float32 {
	w := t.ColumnWidths
	if w == nil {
		return 0.0
	}
	v, ok := w[label]
	if !ok {
		return 0.0
	}
	return v
}

// UI contains ui settings.
type UI struct {
	Table Table `json:"table"`
}

// Config contains all configuration for the application.
type Config struct {
	Version    string `json:"version"`
	ConfigDir  string `json:"config_dir"`
	ConfigFile string `json:"config_file"`
	DBDriver   string `json:"db_driver"` // NOTE deprecated.
	DBFile     string `json:"db_file"`
	UI         UI
}

// GetUITable grab table by name if not found returns an new table.
func (c *Config) GetUITable() *Table {
	if c.UI.Table.ColumnsHidden == nil {
		c.UI.Table.ColumnsHidden = make([]string, 0)
	}
	if c.UI.Table.ColumnWidths == nil {
		c.UI.Table.ColumnWidths = make(map[string]float32)
	}
	t := &c.UI.Table
	return t
}

// save to config file.
func (c *Config) Save() error {
	const op status.Op = "config.save"
	path := c.ConfigFile

	file, err := os.Create(path)
	if err != nil {
		return status.E(op, status.Unexpected, status.LevelError, err)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(c)
	return nil
}

func (c *Config) Open() (*Config, error) {
	return Load(c.ConfigDir)
}

// Load config from root directory. The file name will 'config.json'
// and if file is not found it will return new Config type to be saved later.
func Load(root string) (*Config, error) {
	const op status.Op = "config.load"

	file := "config.json"
	path := filepath.Join(root, file)

	fileIO, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		cfg := &Config{
			ConfigDir:  root,
			ConfigFile: path,
		}
		return cfg, nil
	}
	if err != nil {
		return nil, status.E(op, status.Unexpected, status.LevelError, err)
	}
	defer fileIO.Close()

	cfg := &Config{}

	decoder := json.NewDecoder(fileIO)

	err = decoder.Decode(cfg)
	if err != nil {
		return nil, status.E(op, status.FailedToDecode, status.LevelError, err)
	}
	err = backupV1(path, cfg)
	return cfg, err
}

// backupV1 create a backup of the old config if it is.
func backupV1(path string, cfg *Config) error {
	const op status.Op = "config.backupV1"
	if cfg.Version != "" || cfg.DBFile == "" {
		return nil
	}

	cfg.Version = "2.0.0"

	to, err := os.Create(path + ".bak")
	if err != nil {
		return status.E(op, status.FailedToCreate, status.LevelError, err)
	}
	from, err := os.Open(path)
	if err != nil {
		return status.E(op, status.FailedToOpen, status.LevelError, err)
	}

	_, err = from.WriteTo(to)
	if err != nil {
		return status.E(op, status.Unexpected, status.LevelError, err)
	}
	return nil
}
