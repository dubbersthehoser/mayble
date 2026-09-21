package snapshot

import (
	"fmt"
	"slices"
	"sync/atomic"

	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/viewmodel/display"
)

var version atomic.Int64

var Current atomic.Pointer[Snapshot]

func init() {
	ss := NewSnapshot([]models.BookEntry{})
	Current.Store(ss)
}

// Snapshot is a inmutable table of BookEntry data.
type Snapshot struct {
	data         []models.BookEntry
	uniqueGenres []string
	version      int64

	rowToID map[int]int64
	idToRow map[int64]int
}

func NewSnapshot(data []models.BookEntry) *Snapshot {
	ss := &Snapshot{
		data:         data,
		version:      version.Load(),
		uniqueGenres: make([]string, 0),
		rowToID:      make(map[int]int64),
		idToRow:      make(map[int64]int),
	}

	for row, book := range data {
		ss.uniqueGenres = append(ss.uniqueGenres, book.Genre)
		ss.rowToID[row] = book.ID
		ss.idToRow[book.ID] = row
	}

	version.Add(1)
	return ss
}

func (ss *Snapshot) Version() int64 {
	return ss.version
}

func (ss *Snapshot) Get(p models.Cell) (string, error) {
	row, err := ss.IDToRow(p.ID)
	if err != nil {
		return "", fmt.Errorf("get %d: %w", row, err)
	}
	fields := display.EntryValues(&ss.data[row])
	idx := slices.Index(models.BookEntryFields(), p.Column)
	if idx == -1 {
		return "", fmt.Errorf("get %s: invalid column", p.Column)
	}
	return fields[idx], nil
}

func (ss *Snapshot) GetBookEntryByRow(row int) (*models.BookEntry, error) {
	if row >= len(ss.data) || row < 0 {
		return nil, fmt.Errorf("getbookentry %d: row out of bounds", row)
	}
	return &ss.data[row], nil
}

func (ss *Snapshot) GetBookEntryByID(id int64) (*models.BookEntry, error) {
	row, err := ss.IDToRow(id)
	if err != nil {
		return nil, err
	}
	return ss.GetBookEntryByRow(row)
}

func (ss *Snapshot) RowToID(row int) (int64, error) {
	id, ok := ss.rowToID[row]
	if !ok {
		return 0, fmt.Errorf("row_to_id %d: row not found", row)
	}
	return id, nil
}

func (ss *Snapshot) IDToRow(id int64) (int, error) {
	row, ok := ss.idToRow[id]
	if !ok {
		return 0, fmt.Errorf("id_to_row %d: id not found", id)
	}
	return row, nil
}

func (ss *Snapshot) Size() (rows, cols int) {
	rows = len(ss.data)
	if rows == 0 {
		return 0, 0
	}
	cols = len(models.BookEntryFields())
	return
}

func (ss *Snapshot) Length() int {
	return len(ss.data)
}

func (ss *Snapshot) UniqueGenres() []string {
	return ss.uniqueGenres
}
