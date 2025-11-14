package api

import (
	"testing"
)

func TestSQLiteFS(t *testing.T) {
	testDirs := []string{
		"sqlite/schemas",
		"sqlite/queries",
	}
	f := SQLiteFS

	for _, path := range testDirs {
		_, err := f.ReadDir(path)
		if err != nil {
			t.Fatal(err)
		}
	}
}
