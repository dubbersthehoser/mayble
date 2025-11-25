package driver

import (
	"errors"

	"github.com/dubbersthehoser/mayble/internal/storage"
	"github.com/dubbersthehoser/mayble/internal/storage/memory"
)

func Load(driver string, path string) (storage.BookLoanStore, error) {
	switch driver {
	case "memory":
		return memory.NewStorage(), nil
	default:
		return nil, errors.New("storage driver not found")
	}
}
