package storage

import (
	"io"

	"github.com/dubbersthehoser/mayble/internal/data"
)

type Importer interface {
	ImportBooks(io.Reader) ([]data.BookLoan, error)
}

type Exporter interface {
	ExportBooks(io.Writer, []data.BookLoan) error
}
