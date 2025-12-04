package launcher

import (
	"fmt"
	"os"
	"errors"
	"path/filepath"

	"github.com/dubbersthehoser/mayble/internal/settings"
)

/* Lancher Options */

type Option func(s *settings.Settings) 

func WithDBPath(path string) Option {
	return func(s *settings.Settings) {
		s.DBPath = path
	}
}
func WithDBDriver(driver string) Option {
	return func(s *settings.Settings) {
		s.DBDriver = driver
	}
}

func WithConfigPath(path string) Option {
	return func(s *settings.Settings) {
		s.ConfigPath = path
	}
}

func WithConfigDir(dir string) Option {
	return func(s *settings.Settings) {
		s.ConfigDir = dir
	}
}


/* Default Settings */

// defaultConfigDir function must be ran first before other default functions.
func defaultConfigDir(s *settings.Settings) { 
	if s.ConfigDir != "" {
		return
	}
	dir, err := determinConfigDir("mayble")
	if err == nil {
		s.ConfigDir = dir
	} else {
		panic("could not determin configuration path for application")
	}
}

func defaultDBPath(s *settings.Settings) {
	if s.ConfigDir == "" {
		panic("configuration path was not set before setting default dbPath.")
	}
	if s.DBPath != "" {
		return
	}
	s.DBPath = filepath.Join(s.ConfigDir, "db.sqlite")
}


func defaultConfigPath(s *settings.Settings) {
	if s.ConfigDir == "" {
		panic("configuration path was not set before setting default configPath.")
	}
	if s.ConfigPath != "" {
		return
	}
	s.ConfigPath = filepath.Join(s.ConfigDir, "config.json")
}

func defaultDBDriver(s *settings.Settings) {
	s.DBDriver = "memory"
}

func determinConfigDir(appName string) (string, error) {

	// NOTE: This Function could be better, but it works. Has cased unwanted panics
	//       from defaultConfigDir() in previous verions (or is that a bug with defaultConfigDir()?).

	userConfig, err1 := os.UserConfigDir()
	userHome, err2 := os.UserHomeDir()
	if err1 != nil && err2 != nil {
		return "", fmt.Errorf("failed to determin config directory: UserHome: %w, UserConfig: %w", err2, err1)
	}

	var (
		ThePath    string
		fullConfig string
		fullHome   string
	)

	if err1 == nil {
		fullConfig = filepath.Join(userConfig, appName)
	}
	if err2 == nil {
		fullHome = filepath.Join(userHome, "." + appName)
	}

	if fullConfig == "" {
		ThePath = fullHome
	}
	if fullHome == "" {
		ThePath = fullConfig
	}

	if ThePath == "" {
		// return one that exsits
		paths := []string{fullConfig, fullHome}
		for _, path := range paths {
			_, err := os.Stat(path)
			if err != nil {
				continue
			} else {
				ThePath = path
			}
		}
	}

	// if all checks above failed to determin ThePath
	if ThePath == "" {
		ThePath = fullConfig
	}

	err := os.Mkdir(ThePath, 0744)
	var (
		Exist   bool = errors.Is(err, os.ErrExist)
		NoError bool = err != nil
	)

	if !Exist && NoError {
		return ThePath, fmt.Errorf("config: %w", err)
	}

	return ThePath, nil
}

