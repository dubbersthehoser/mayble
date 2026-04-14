package config

import (
	"encoding/json"
	"errors"
	"os"
	"fmt"
	"path/filepath"
)

// Table config for table view.
type Table struct {
	ColumnsHidden []string           `json:"hidden_columns"`
	ColumnWidths  map[string]float32 `json:"column_width"`
}

// UI contains ui settings.
type UI struct {
	Table Table `json:"table"`
}

// SetColumnWidth for column lable for size s.
func (u *UI) SetColumnWidth(label string, s float32) {
	w := u.Table.ColumnWidths
	if w == nil {
		w = make(map[string]float32)
	}
	w[label] = s
}

// GetColumnWidth from named label.
func (u *UI) GetColumnWidth(label string) float32 {
	w := u.Table.ColumnWidths
	if w == nil {
		return 0.0
	}
	v, ok := w[label]
	if !ok {
		return 0.0
	}
	return v
}

func (u *UI) SetHiddenColumns(headers []string) {
	u.Table.ColumnsHidden = headers
}

func (u *UI) GetHiddenColumns() []string {
	return u.Table.ColumnsHidden
}


// Config contains all configuration for the application.
type Config struct {
	Version    string `json:"version"`
	ConfigDir  string `json:"config_dir"`
	ConfigFile string `json:"config_file"`
	DBDriver   string `json:"db_driver"` // NOTE deprecated.
	DBFile     string `json:"db_file"`
	UI
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

func (c *Config) Open() (*Config, error) {
	return Load(c.ConfigDir)
}

// Load config from root directory. The file name will 'config.json'
// and if file is not found it will return new Config type to be saved later.
func Load(root string) (*Config, error) {

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
		return nil, fmt.Errorf("config.load %w", err)
	}
	defer fileIO.Close()

	cfg := &Config{}

	decoder := json.NewDecoder(fileIO)

	err = decoder.Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("config.load %w", err)
	}
	err = backupV1(path, cfg)
	return cfg, err
}

// backupV1 create a backup of the old config if it is.
func backupV1(path string, cfg *Config) error {
	if cfg.Version != "" || cfg.DBFile == "" {
		return nil
	}

	cfg.Version = "2.0.0"

	to, err := os.Create(path + ".bak")
	if err != nil {
		return fmt.Errorf("config.backup %w", err)
	}
	from, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("config.backup %w", err)
	}

	_, err = from.WriteTo(to)
	if err != nil {
		return fmt.Errorf("config.backup %w", err)
	}
	return nil
}
