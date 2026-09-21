package table

import (
	"cmp"
	"context"
	"slices"

	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/search"
	"github.com/dubbersthehoser/mayble/internal/snapshot"
)

type SearchResult struct {
	Point search.Point
	Score int
}

func searchSearcherWithContext(ctx context.Context, srch *search.Searcher) []SearchResult {

	results := make([]SearchResult, 0)

	for srch.Next() {
		if ctx.Err() != nil {
			return []SearchResult{}
		}
		point := srch.Point()
		score := srch.Score()
		if score == -1 {
			continue
		}
		r := SearchResult{
			Score: score,
			Point: point,
		}
		results = append(results, r)
	}

	if len(results) == 0 {
		return []SearchResult{}
	}

	slices.SortFunc(results, func(a, b SearchResult) int {
		r := cmp.Compare(a.Score, b.Score)
		if r == 0 {
			return cmp.Compare(a.Point.Row, b.Point.Row)
		}
		return r * -1
	})
	return results
}

type tableTraverse struct {
	snapshot *snapshot.Snapshot
	row, col int
	isDone   bool
	setDone  func()
}

func newTableTraverse(ss *snapshot.Snapshot) *tableTraverse {
	tt := &tableTraverse{
		snapshot: ss,
		row:      0,
		col:      -1,
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
	tt.col += 1
	if cols <= tt.col {
		tt.row += 1
		tt.col = 0
	}
	if rows <= tt.row {
		return tt.retDone()
	}
	label := models.BookEntryFields()[tt.col]
	id, _ := tt.snapshot.RowToID(tt.row)
	v, err := tt.snapshot.Get(models.Cell{ID: id, Column: label})
	if err != nil {
		return tt.retDone()
	}
	return v, true
}

func (tt *tableTraverse) IsDone() bool {
	return tt.isDone
}

func (tt *tableTraverse) Point() search.Point {
	return search.Point{Row: tt.row, Col: tt.col}
}

func (tt *tableTraverse) retDone() (string, bool) {
	tt.isDone = true
	return "", false
}

type columnTraverse struct {
	snapshot *snapshot.Snapshot
	row, col int
	setDone  func()
	isDone   bool
}

func newColumnTraverse(ss *snapshot.Snapshot, col int) *columnTraverse {
	ct := &columnTraverse{
		snapshot: ss,
		row:      -1,
		col:      col,
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

	label := models.BookEntryFields()[ct.col]
	id, _ := ct.snapshot.RowToID(ct.row)
	v, err := ct.snapshot.Get(models.Cell{ID: id, Column: label})
	if err != nil {
		return ct.retDone()
	}
	return v, true
}

func (ct *columnTraverse) Point() search.Point {
	return search.Point{Row: ct.row, Col: ct.col}
}

func (ct *columnTraverse) IsDone() bool {
	return ct.isDone
}

func (ct *columnTraverse) retDone() (string, bool) {
	ct.isDone = true
	return "", false
}
