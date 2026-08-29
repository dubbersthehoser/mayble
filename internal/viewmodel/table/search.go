package table

import (
	"sync"
	"github.com/dubbersthehoser/mayble/internal/search"
)

type tableTraverse struct {
	sheet    *Sheet
	row, col int
	isDone   bool
	setDone  func()
}

func newTableTraverse(sh *Sheet) *tableTraverse {
	tt := &tableTraverse{
		sheet: sh,
		row: 0,
		col: -1,
	}

	tt.sheet.data.mu.RLock()

	tt.setDone = sync.OnceFunc(func() {
		tt.isDone = true
		tt.sheet.data.mu.RUnlock()
	})

	return tt
}

func (tt *tableTraverse) Next() (string, bool) {
	rows, cols := tt.sheet.Size()
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
	v, err := tt.sheet.Get(Point{Row: tt.row, Col: tt.col})
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
	tt.setDone()
	return "", true
}

type columnTraverse struct {
	sheet *Sheet
	row, col int
	setDone func()
	isDone  bool
}

func newColumnTraverse(s *Sheet, col int) *columnTraverse {
	ct := &columnTraverse{
		sheet: s,
		row: -1,
		col: col,

	}
	ct.sheet.data.mu.RLock()
	ct.setDone = sync.OnceFunc(func() {
		ct.isDone = true
		ct.sheet.data.mu.RUnlock()
	})
	return ct
}

func (ct *columnTraverse) Next() (string, bool) {
	rows, cols := ct.sheet.Size()
	ct.row += 1
	if rows <= ct.row {
		return ct.retDone()
	}
	if cols <= ct.col {
		return ct.retDone()
	}

	v, err := ct.sheet.Get(Point{Row: ct.row, Col: ct.col})
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
	ct.setDone()
	return "", true
}
