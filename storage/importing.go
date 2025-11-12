package storage

import (
	"io"
)

type Importer interface {
	ImportBooks(io.Reader) ([]BookLoan, error)
}

type Exporter interface {
	ExportBooks( io.Writer, []BookLoan) error
}
