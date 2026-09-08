package snapshot

import (
	"sync/atomic"
	"slices"
	"fmt"

	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/viewmodel/display"
)

var version atomic.Int64

var Current atomic.Pointer[Snapshot]

type Point struct {
	Col  string
	ID   int64
}


// Snapshot is a inmutable table of BookEntry data.
type Snapshot struct {
	data    []models.BookEntry
	version int64

	rowToID map[int]int64
	idToRow map[int64]int
}

func NewSnapshot(data []models.BookEntry) *Snapshot {
	s := &Snapshot{
		data: data,
		version: version.Load(),
		rowToID: make(map[int]int64),
		idToRow: make(map[int64]int),
	}
	version.Add(1)
	return s
}

func (ss *Snapshot) Version() int64 {
	return ss.version
}

func (ss *Snapshot) Get(p Point) (string, error) {
	row, err := ss.IDToRow(p.ID)
	if err != nil {
		return "", fmt.Errorf("get %d: %w", row, err)
	}
	fields := display.EntryValues(&ss.data[row])
	idx := slices.Index(fields, p.Col)
	if idx == -1 {
		return "", fmt.Errorf("get %s: invalid column", p.Col)
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

