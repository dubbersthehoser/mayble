package driver

import (
	"fmt"

	"github.com/dubbersthehoser/mayble/internal/storage"
	"github.com/dubbersthehoser/mayble/internal/storage/memory"
	"github.com/dubbersthehoser/mayble/internal/sqlite"
)

func Load(driver string, path string) (storage.BookLoanStore, error) {
	switch driver {
	case "memory":
		return memory.NewStorage(), nil
	case "sqlite":
		return sqlite.NewStorage(path)
	default:
		return nil, fmt.Errorf("storage: dirver not found '%s'", driver)
	}
}
