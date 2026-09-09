package command

import (
	"github.com/dubbersthehoser/mayble/internal/models"
)

type CellSelect struct {
	Version int64
	Point   models.Cell
	Has     bool
}

type TableSort struct {
	Column string
	Asc    bool
}

type TableSearch struct {
	Column  string
	Pattern string
}

type CreateBookEntry struct {
	Book models.BookEntry
}

type UpdateBookEntry struct {
	Book models.BookEntry
}

type DeleteBookEntry struct {
	BookID int64
}

type ImportFile struct {
	Path string
}

type ExportFile struct {
	Path string
}

type CreateDatabase struct {
	Path string
}

type OpenDatabase struct {
	Path string
}

type TakeSnapshot struct {}


