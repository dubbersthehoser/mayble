package submissions

import (
	"errors"
	"slices"

	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

type List struct {
	items []repo.BookEntry
	m     int
}

func NewList(max int) *List {
	return &List{
		m: max,
		items: make([]repo.BookEntry, 0),
	}
}

func (l *List) Length() int {
	return len(l.items)
}

func (l *List) Get(idx int) (*repo.BookEntry, error) {
	if len(l.items) <= idx || idx < 0 {
		return nil, errors.New("index out of range")
	}
	return &l.items[idx], nil
}

func (l *List) Append(b repo.BookEntry) {
	l.items = append(l.items, b)
}

func (l *List) Clear() {
	l.items = l.items[:0]
}

func (l *List) Pop() *repo.BookEntry {
	if l.Length() == 0 {
		return nil
	}
	idx := len(l.items)-1
	item := l.items[idx]
	l.items = l.items[:idx]
	return &item
}

func (l *List) Remove(idx int) error {
	if len(l.items) <= idx || idx < 0 {
		return errors.New("index out of range")
	}
	l.items = slices.Delete(l.items, idx, idx+1)
	return nil
}

