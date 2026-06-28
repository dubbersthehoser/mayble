package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"fmt"

	"github.com/dubbersthehoser/mayble/internal/models"
)

// Current config version.
const Version string = "2.0.0"

// default values for config.
const (
	windowHeight         float32 = 600
	windowWidth          float32 = 900
	windowFullScreen     bool = false
	windowCenterOnScreen bool = true
	minColumnSize        float32 = 100
	tableSortByIdx       int = models.IdxTitle               
)

var (
	ErrIsOldConfig error = errors.New("old config")
)

type OldConfig struct {
	ConfigDir  string `json:"config_dir"`
	ConfigFile string `json:"config_file"`
	DBDriver   string `json:"db_driver"`
	DBFile     string `json:"db_file"`
}

type Header struct {
	Name     string  `json:"name"`
	IsHidden bool    `json:"is_hidden"`
	Width    float32 `json:"width"`
}

// Config contains all configuration for the application.
type Config struct{
	Version    string `json:"version"`
	ConfigFile string `json:"config_file"`
	DBFile     string `json:"db_file"`
	UI         UI     `json:"ui"`
}

type UI struct {
	Headers              map[int]Header `json:"headers"`
	TableSortBy          int            `json:"table_sort_by"`
	TableAscending       bool           `json:"table_ascending"`
	TableMinWidth        float32        `json:"table_min_width"`
	WindowHeight         float32        `json:"window_height"`
	WindowWidth          float32        `json:"window_width"`
	WindowFullScreen     bool           `json:"window_fullscreen"`
	WindowCenterOnScreen bool           `json:"window_center_on_screen"`
}

// NewConfigWithDefaults returns a fresh configuration for application.
func NewConfigWithDefaults(configFile string) *Config {
	cfg := &Config{
		Version: Version,
		ConfigFile: configFile,
		UI: UI{
			Headers: make(map[int]Header),
			TableMinWidth:  minColumnSize,
			TableSortBy:  tableSortByIdx,
			WindowHeight: windowHeight,
			WindowWidth: windowWidth,
			WindowFullScreen: windowFullScreen,
			WindowCenterOnScreen: windowCenterOnScreen,
		},
	}

	headerLables := models.BookEntryFields()
	for i, label := range headerLables {
		header := Header{
			Name: label,
		}
		if label == headerLables[models.IdxID] {
			header.IsHidden = true
		}
		cfg.UI.Headers[i] = header
	}

	return cfg
}

// Save to config file.
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

// Load config file form file path. 
func Load(path string) (*Config, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	defer file.Close()

	raw, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	if isOld(raw) {
		return nil, ErrIsOldConfig
	}

	cfg := &Config{}

	err = json.Unmarshal(raw, cfg)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return cfg, err
}

// Migrate old config to current config.
func Migrate(path string) (*Config, error) {
	cfg := NewConfigWithDefaults(path)
	err := backup(path)
	return cfg, err
}

// isOld check wether json string is in the format of the old config.
func isOld(jsonBytes []byte) bool {
	oldCfg := &OldConfig{}
	err := json.Unmarshal(jsonBytes, oldCfg)
	if err != nil {
		return false
	}
	// DBDriver was removed for v2 config.
	return oldCfg.DBDriver != ""
}

// backup create a backup of the old config.
func backup(path string) error {
	to, err := os.Create(path + ".bak")
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	from, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	_, err = from.WriteTo(to)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	return nil
}
