package porting

import (
	"io"
	"fmt"
	"errors"
	"path/filepath"

	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/porting/csv"
)

type NamedWriteCloser interface {
	io.WriteCloser
	Name() string
}

type NamedReadCloser interface {
	io.ReadCloser
	Name() string  
}

type BookLoanImporter interface {
	ImportBookLoans(io.Reader) ([]app.BookLoan, error)
}

type BookLoanExporter interface {
	ExportBookLoans(io.Writer, []app.BookLoan) error
}

type BookLoanPorter interface {
	BookLoanImporter
	BookLoanExporter
}


func GetBookLoanPorterByName(name string) (BookLoanPorter, error) {
	ext := filepath.Ext(name)
	switch ext {
	case ".csv":
		return &csv.BookLoanPorter{}, nil
	case "":
		return nil, errors.New("porting: unspecified file extention")
	default:
		return nil, fmt.Errorf("porting: unsupported file extention: '%s'", ext)
	}


}
