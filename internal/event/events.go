package event

import (
	"github.com/dubbersthehoser/mayble/internal/models"
)

type SelectedFrom int
const (
	SelectedFromUser   SelectedFrom = iota
	SelectedFromSearch
)

type CellSelected struct {
	Has     bool
	Point   models.Cell
	Version int64
}

type TableSorted struct {
	Version  int64
	Sorted   []int64
	Column   string
	Asc      bool
}

type TableSearched struct {
	Version int64
	Pattern string
	Points  []models.Cell
	Scores  []int
}

type HiddenColumn struct {
	Hidden  []bool
}

func (hc *HiddenColumn) ShownColumns() []string {
	columns := make([]string, 0)
	for i, label := range models.BookEntryFields() {
		if !hc.Hidden[i] {
			columns = append(columns, label)
		}
	}
	return columns
}

func (hc *HiddenColumn) HiddenColumns() []string {
	columns := make([]string, 0)
	for i, label := range models.BookEntryFields() {
		if hc.Hidden[i] {
			columns = append(columns, label)
		}
	}
	return columns
}

type StoredSnapshot struct {
	Version int64
	Failed bool
	Message string
}

type CreatedBookEntry struct {
	Message string
	Failed bool
}

type DeletedBookEntry struct {
	Message string
	Failed bool
}

type UpdatedBookEntry struct {
	Message string
	Failed bool
}

type ImportedFile struct {
	Path    string
	Message string
	Failed  bool
}

type ExportedFile struct {
	Path    string
	Message string
	Failed  bool
}

type OpenedDatabase struct {
	Path    string
	Message string
	Failed  bool
	Err     error
}

type CreatedDatabase struct {
	Path string
	Message string
	Failed bool
	Err     error
}

