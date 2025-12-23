package config

import (
	"os"
	"testing"
	"path/filepath"
)

func TestConfig(t *testing.T) {
	
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	err = file.Close()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	os.Remove(file.Name())

	dir := filepath.Dir(file.Name())
	fileName := filepath.Base(file.Name())

	cfg, err := Load(dir, fileName)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	

	err = cfg.SetDBFile("test.db")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	err = cfg.SetDBDriver("memory")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}
