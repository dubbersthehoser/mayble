package config

import (
	"testing"
	"os"
	"io"
)

func Test_backup(t *testing.T) {

	dataInFile := []byte(`This is data to backup`)

	// create test file and add test data to it.
	tmpFile, err := os.CreateTemp("", "backup-testing_*")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	tmpFile.Write(dataInFile)
	tmpFile.Close()

	t.Logf("created file '%s'", tmpFile.Name())

	// backup.
	err = backup(tmpFile.Name())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// open back up read data then compare.
	backupName := tmpFile.Name() + ".bak"
	bakFile, err := os.Open(backupName)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	dataInBak, err := io.ReadAll(bakFile)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if string(dataInBak) != string(dataInFile) {
		t.Fatalf("expect '%s', got '%s'", dataInFile, dataInBak)
	}
}

func Test_isOld(t *testing.T) {
	oldConfig := []byte(`{
		"config_dir": "/home/somthing/.config/mayble",
		"config_file": "/home/somthing/.config/mayble/config.json",
		"db_driver": "memory",
		"db_file": "/home/somthing/Documnets/mayble.db"
	}`)
	if !isOld(oldConfig) {
		t.Fatalf("expect %t, got %t", true, false)
	}

	newConfig := []byte(`
		{
			"version": "2.0.0",
			"config_file": "/home/something/.config/mayble/config.json",
			"db_file": "/home/something/Documents/mayble.db",
			"ui": {
				"open_body": 0,
				"headers": {},
				"table_sort_by": "Table",
				"table_ascending": false
			}
		}
	`)

	if isOld(newConfig) {
		t.Fatalf("expect %t, got %t", false, true)
	}
}




