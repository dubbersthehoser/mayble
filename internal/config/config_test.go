package config

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func unexpectError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestConfig(t *testing.T) {

	dir, err := os.MkdirTemp("", "")
	unexpectError(t, err)
	t.Cleanup(func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Log(err)
		}
	})

	cfg, err := Load(dir)
	unexpectError(t, err)

	if cfg.ConfigDir != dir {
		t.Fatalf("expect '%s', got '%s'", dir, cfg.ConfigDir)
	}
	if cfg.ConfigFile != filepath.Join(dir, "config.json") {
		t.Fatalf("expect '%s', got '%s'", dir, cfg.ConfigDir)
	}

	cfg.DBFile = "test.db"

	table := cfg.GetUITable()
	if table == nil {
		t.Fatalf("unexpected nil value returned")
	}

	if table.ColumnsHidden == nil {
		t.Fatalf("table settings column hidden is nil")
	}

	hidden := []string{
		"Title",
		"Author",
	}
	table.ColumnsHidden = hidden

	err = cfg.Save()
	unexpectError(t, err)

	cfg, err = cfg.Open()
	unexpectError(t, err)

	table = cfg.GetUITable()

	if r := slices.Compare(table.ColumnsHidden, hidden); r != 0 {
		t.Fatalf("expect \n\t%#v\ngot\n\t%#v\n", hidden, table.ColumnsHidden)
	}
}
