package config

import (
	"os"
	"testing"
)

func TestGetDefaultDir(t *testing.T) {
	appName := "Testing-Application-Name"

	path, err := GetDefaultDir(appName)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	t.Logf("created: '%s'", path)
	defer func(){
		_ = os.Remove(path)
		t.Logf("removed: '%s'", path)
	}()

	_, err = os.Stat(path)
	if err != nil {
		t.Fatalf("path was not created '%s'", path)
	}


	pathTwo, err := GetDefaultDir(appName)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if pathTwo != path {
		t.Fatalf("expect '%s', got '%s'", path, pathTwo)
	}

}
