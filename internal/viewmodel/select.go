package viewmodel

import (
	"fyne.io/fyne/v2/data/binding"

	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

type EntrySelect struct {
	retriever  repo.BookRetriever
	isSelected bool
	selected   int64

	cellPos struct {
		col int
		row int
	}

	l *listener
}

func newEntrySelect(r repo.BookRetriever) *EntrySelect {
	e := &EntrySelect{
		retriever: r,
		l:         &listener{},
	}
	return e
}

func (e *EntrySelect) getBook() (*repo.BookEntry, error) {
	b, err := e.retriever.GetBookByID(e.selected)
	return &b, err
}

func (e *EntrySelect) selectID(id int64, notify bool) {
	e.selected = id
	e.isSelected = true
	if notify {
		e.l.notify()
	}
}
func (e *EntrySelect) selectCell(row, col int, notify bool) {
	e.cellPos.row = row
	e.cellPos.col = col
	e.isSelected = true
	if notify {
		e.l.notify()
	}
}

func (e *EntrySelect) SelectedCell() (row, col int) {
	return e.cellPos.row, e.cellPos.col
}

func (e *EntrySelect) unselect(notify bool) {
	e.isSelected = false
	if notify {
		e.l.notify()
	}
}

func (e *EntrySelect) HasSelected() bool {
	return e.isSelected
}

func (e *EntrySelect) AddListener(l binding.DataListener) {
	e.l.AddListener(l)
}
