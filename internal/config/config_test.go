package config

import (
	"os"
	"testing"
)

func TestConfigRoot(t *testing.T) {
	
	dir := os.TempDir()
	filename := "test.json"
	dbPath := "/test/path.db"

	tmpDir, err := os.MkdirTemp(dir, "test_*")
	t.Log("created: ", tmpDir)
	defer func() {
		os.RemoveAll(dir)
		t.Log("removed: ", tmpDir)
	}()

	configRoot, err := OpenConfigRoot(tmpDir)
	configRoot.SetFileName(filename)
	if err != nil {
		t.Error(err)
	}

	t.Log("opening config")
	config, err := configRoot.Open()
	if err != nil {
		t.Error(err)
	}

	config.DBPath = dbPath

	t.Log("saving config")
	if err := configRoot.Save(config); err != nil {
		t.Error(err)
	}

	t.Log("copying and opening config, then comparing")
	configCopy := *config
	config, err = configRoot.Open()
	if err != nil {
		t.Error(err)
	}

	if config.DBPath != configCopy.DBPath {
		t.Errorf("expect config.DBPath='%s', got config.DBPath='%s'", configCopy.DBPath, config.DBPath)
	}

}
