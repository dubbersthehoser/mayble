package config

import (
	"os"
	"testing"
	"path/filepath"
	"slices"
)

func unexpectError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestConfig(t *testing.T) {
	
	file, err := os.CreateTemp("", "")
	unexpectError(t, err)

	t.Cleanup(func(){
		os.Remove(file.Name())
	})


	err = file.Close()
	unexpectError(t, err)

	err = os.Remove(file.Name())
	unexpectError(t, err)

	dir := filepath.Dir(file.Name())
	fileName := filepath.Base(file.Name())

	cfg, err := Load(dir, fileName)
	unexpectError(t, err)

	err = cfg.SetDBFile("test.db")
	unexpectError(t, err)

	err = cfg.SetDBDriver("memory")
	unexpectError(t, err)

	table := cfg.GetTable("Main")

	if table == nil {
		t.Fatalf("unexpected nil value return")
	}

	if table.Name != "Main" {
		t.Fatalf("table name was not set")
	}

	if table.Settings.ColsHidden == nil {
		t.Fatalf("table settings column hidden is nil")
	}

	hidden := []string{
		"Title",
		"Author",
	}
	table.Settings.ColsHidden = hidden
	err = cfg.UpdateTable("Main", table)
	unexpectError(t, err)

	cfg, err = cfg.Open()
	unexpectError(t, err)

	table = cfg.GetTable("Main")

	if r := slices.Compare(table.Settings.ColsHidden, hidden); r != 0 {
		t.Fatalf("expect \n\t%#v\ngot\n\t%#v\n", table.Settings.ColsHidden, hidden)
	}
}
