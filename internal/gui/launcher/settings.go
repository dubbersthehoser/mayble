package launcher

import (
	"os"
	"fmt"
	"errors"
	"path/filepath"

)

type Settings struct {
	configDir  string
	configPath string
	dbPath     string
}

type Option func(s *Settings) 

func WithDBPath(path string) Option {
	return func(s *Settings) {
		s.configDir = path
	}
}

func WithConfigPath(path string) Option {
	return func(s *Settings) {
		s.configPath = path
	}
}


func defaultDBPath(s *Settings) {
	if s.configDir == "" {
		panic("configuration path was not set before setting default dbPath.")
	}
	s.dbPath = filepath.Join(s.configDir, "db.sqlite")
}

func defaultConfigPath(s *Settings) {
	if s.configDir == "" {
		panic("configuration path was not set before setting default configPath.")
	}
	s.configPath = filepath.Join(s.configDir, "config.json")
}

func defaultConfigDir(s *Settings) {
	dir, err := determinConfigDir("mayble")
	if err == nil {
		s.configDir = dir
	} else {
		panic("could not determin configuration path for application")
	}
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
