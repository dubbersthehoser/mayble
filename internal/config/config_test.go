package config

import (
	"testing"
	"os"
	"io"
)

func Test_backup(t *testing.T) {
	

	dataInFile := []byte(`This is data to backup`)

	tmpFile, err := os.CreateTemp("", "backup-testing_*")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	tmpFile.Write(dataInFile)
	tmpFile.Close()

	t.Logf("created file '%s'", tmpFile.Name())

	err = backup(tmpFile.Name())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

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
