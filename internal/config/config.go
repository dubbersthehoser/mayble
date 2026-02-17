package config

import (
	"os"
	"errors"
	"encoding/json"
	"path/filepath"
)

// Table config for table view.
type Table struct {
		Name     string
		Settings struct {
			ColsHidden []string `json:"hidden_columns"`
			ColWidths   []float32 `json:"column_width"`
		} `json:"settings"`
}

func (t *Table) SetColWidth(i int, s float32) error {
	w := &t.Settings.ColWidths
	if w == nil {
		*w = make([]float32, 0)
	}

	if i < 0 && i > len(*w){
		return errors.New("invalid index")
	}

	if len(*w) == i {
		*w = append(*w, s)
	} else {
		(*w)[i] = s
	} 
	return nil
}

// GetColWidth of column i. When is not in or out of range returns -1.0.
func (t *Table) GetColWidth(i int) float32 {
	w := &t.Settings.ColWidths
	if w == nil {
		return -1.0
	} 
	if i >= len(*w) || i < 0 {
		return -1.0
	}
	return (*w)[i]
}

type UI struct {
	Tables map[string]Table `json:"tables"`
}

// Config contains all configuration for the application.
//
type Config struct {
	Version    string `json:"version"` // NOTE check if this is empty then it's the 1.0.0 config.
	ConfigDir  string `json:"config_dir"`
	ConfigFile string `json:"config_file"`
	DBDriver   string `json:"db_driver"`
	DBFile     string `json:"db_file"`
	UI         UI
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

// UpdateUITable create or update table settings.
// Param name will set to t.Name.
// Returns error when t is nil or when c.save errors.
func (c *Config) UpdateUITable(name string, t *Table) error {
	if t == nil {
		return errors.New("passed nil value")
	}
	if c.UI.Tables == nil {
		c.UI.Tables = make(map[string]Table)
	}
	t.Name = name
	c.UI.Tables[name] = *t
	return c.save()
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
	}
	return &t
}

// save to config file.
func (c *Config) save() error {
	path := c.ConfigFile

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(c)
	return nil
}

func (c *Config) Open() (*Config, error) {
	file := filepath.Base(c.ConfigFile)
	return Load(c.ConfigDir, file)
}

// Load config with root directory and from file name. Sets .ConfigFile with 
// joined root and file. When the .ConfigFile is not found it will return new config.
func Load(root, file string) (*Config, error) {

	path := filepath.Join(root, file)

	fileIO, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Config{
			ConfigDir: root,
			ConfigFile: path,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	defer fileIO.Close()

	cfg := &Config{}

	decoder := json.NewDecoder(fileIO)
	
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

