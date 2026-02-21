package config

import (
	"os"
	"errors"
	"encoding/json"
	"path/filepath"

	"github.com/dubbersthehoser/mayble/internal/status"
)

// Table config for table view.
type Table struct {
		Name     string
		Settings struct {
			ColsHidden  []string           `json:"hidden_columns"`
			ColWidths   map[string]float32 `json:"column_width"`
		} `json:"settings"`
}


func (t *Table) SetColWidth(label string, s float32) error {
	w := t.Settings.ColWidths
	if w == nil {
		w = make(map[string]float32)
	}
	w[label] = s
	return nil
}

// GetColWidth of column i. When is not in or out of range returns -1.0.
func (t *Table) GetColWidth(label string) float32 {
	w := t.Settings.ColWidths
	if w == nil {
		return -1.0
	} 
	v, ok := w[label]
	if !ok {
		return 0.0
	}
	return v
}

type UI struct {
	Tables map[string]Table `json:"tables"`
}

// Config contains all configuration for the application.
type Config struct {
	Version    string `json:"version"` // NOTE check if this is empty then it's the 1.0.0 config.
	ConfigDir  string `json:"config_dir"`
	ConfigFile string `json:"config_file"`
	DBDriver   string `json:"db_driver"`
	DBFile     string `json:"db_file"`
	UI         UI
}

// UpdateUITable create or update table settings.
// Param name will set to t.Name.
// Returns error when t is nil or when c.save errors.
func (c *Config) UpdateUITable(name string, t *Table) error {
	const op status.Op = "config.UpdateUITable"
	if t == nil {
		return status.E(op, status.LevelDebug, "table was nil") 
	}
	if c.UI.Tables == nil {
		c.UI.Tables = make(map[string]Table)
	}
	t.Name = name
	c.UI.Tables[name] = *t
	return nil
}

// GetUITable grab table by name if not found returns an new table.
func (c *Config) GetUITable(name string) *Table {
	if c.UI.Tables == nil {
		c.UI.Tables = make(map[string]Table)
	}
	t, ok := c.UI.Tables[name]
	if !ok {
		t = Table{
			Name: name,
		}
		t.Settings.ColsHidden = make([]string, 0)
		t.Settings.ColWidths = make(map[string]float32)
	}
	return &t
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
	return Load(c.ConfigFile)
}

// Load config from root directory. 
// The config file will be named 'config.json' and if file is not found it will
// return new Config type.
func Load(root string) (*Config, error) {
	const op status.Op = "config.load"

	file := "config.json"
	path := filepath.Join(root, file)

	fileIO, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		cfg := &Config{
			ConfigDir: root,
			ConfigFile: path,
		}
		return cfg, status.E(op, status.FileNotFound, status.LevelInfo, err)
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
	if !(cfg.Version == "") {
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
