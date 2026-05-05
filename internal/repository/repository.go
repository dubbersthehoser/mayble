package repository

import (
	"io"

	"github.com/dubbersthehoser/mayble/internal/models"
)


type CSVHandler interface {
	ImportFile(path string) error
	ExportFile(path string) error
}


type BookRetriever interface {
	GetAllBooks(Variant) ([]models.BookEntry, error)
	GetBookByID(id int64) (models.BookEntry, error)
}

type BookStore interface {
	BookCreator
	BookUpdator
	BookDeletor
}

type GenreRetriever interface {
	GetUniqueGenres() ([]string, error)
}

type BookCreator interface {
	CreateBook(*models.BookEntry) (int64, error)
}

type BookUpdator interface {
	UpdateBook(*models.BookEntry) error
}

type BookDeletor interface {
	DeleteBook(id int64) error
}

type BookImporter interface {
	BookImport(io.Reader) error
}

type BookExporter interface {
	BookExport(io.Writer, []BookEntry) error
}
