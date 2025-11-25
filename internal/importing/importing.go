package importing

import (
	"io"

	"github.com/dubbersthehoser/mayble/internal/app"
)

type Importer interface {
	ImportBooks(io.Reader) ([]app.BookLoan, error)
}

type Exporter interface {
	ExportBooks(io.Writer, []app.BookLoan) error
}
