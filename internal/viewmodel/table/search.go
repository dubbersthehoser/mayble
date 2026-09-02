package table

import (
	"github.com/dubbersthehoser/mayble/internal/search"
	"github.com/dubbersthehoser/mayble/internal/snapshot"
)

type tableTraverse struct {
	snapshot    *snapshot.Snapshot
	row, col int
	isDone   bool
	setDone  func()
}

func newTableTraverse(ss *snapshot.Snapshot) *tableTraverse {
	tt := &tableTraverse{
		snapshot: ss,
		row: 0,
		col: -1,
	}
	tt.setDone = func() {
		tt.isDone = true
	}

	return tt
}

func (tt *tableTraverse) Next() (string, bool) {
	rows, cols := tt.snapshot.Size()
	if rows == 0 {
		return tt.retDone()
	}
	if cols <= tt.col {
		tt.row += 1
		tt.col = 0
	}
	if rows <= tt.row {
		return tt.retDone()
	}
	v, err := tt.snapshot.Get(snapshot.Point{Row: tt.row, Col: tt.col})
	if err != nil {
		return tt.retDone()
	}
	return v, false
}

func (tt *tableTraverse) IsDone() bool {
	return tt.isDone
}

func (tt *tableTraverse) Point() search.Point {
	return search.Point{Row: tt.row, Col: tt.col}
}

func (tt *tableTraverse) retDone() (string, bool) {
	tt.isDone = true
	return "", true
}

type columnTraverse struct {
	snapshot *snapshot.Snapshot
	row, col int
	setDone func()
	isDone  bool
}

func newColumnTraverse(ss *snapshot.Snapshot, col int) *columnTraverse {
	ct := &columnTraverse{
		snapshot: ss,
		row: -1,
		col: col,

	}
	return ct
}

func (ct *columnTraverse) Next() (string, bool) {
	rows, cols := ct.snapshot.Size()
	ct.row += 1
	if rows <= ct.row {
		return ct.retDone()
	}
	if cols <= ct.col {
		return ct.retDone()
	}

	v, err := ct.snapshot.Get(snapshot.Point{Row: ct.row, Col: ct.col})
	if err != nil {
		return ct.retDone()
	}
	return v, false
}

func (ct *columnTraverse) Point() search.Point {
	return search.Point{Row: ct.row, Col: ct.col}
}

func (ct *columnTraverse) IsDone() bool {
	return ct.isDone
}

func (ct *columnTraverse) retDone() (string, bool) {
	ct.isDone = true
	return "", true
}
