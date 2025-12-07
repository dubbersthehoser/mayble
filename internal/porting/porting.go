package porting

import (
	"io"
	"errors"
	"path/filepath"

	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/porting/csv"
)

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


func GetBookLoanPorterByFilePath(filePath string) (BookLoanPorter, error) {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".csv":
		return &csv.BookLoanPorter{}, nil
	default:
		return nil, errors.New("porting: driver not found for file path.")
	}


}
