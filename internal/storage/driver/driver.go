package driver

import (
	"fmt"
	"errors"

	"github.com/dubbersthehoser/mayble/internal/storage"
	"github.com/dubbersthehoser/mayble/internal/storage/memory"
	"github.com/dubbersthehoser/mayble/internal/storage/sqlite"
)

func Load(driver string, path string) (storage.BookLoanStore, error) {
	switch driver {
	case "memory":
		return memory.NewStorage(), nil
	case "sqlite"
		return sqlite.NewStorage()
	default:
		return nil, fmt.Errorf("storage: dirver not found '%s'", driver)
	}
}
