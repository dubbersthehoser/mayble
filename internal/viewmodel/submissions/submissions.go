package submissions

import (
	"errors"
	"slices"

	"github.com/dubbersthehoser/mayble/internal/models"
)

type List struct {
	items []models.BookEntry
	m     int
}

func NewList(max int) *List {
	return &List{
		m: max,
		items: make([]models.BookEntry, 0),
	}
}

func (l *List) Cap() int {
	return l.m
}

func (l *List) Length() int {
	return len(l.items)
}

func (l *List) Get(idx int) (*models.BookEntry, error) {
	if len(l.items) <= idx || idx < 0 {
		return nil, errors.New("index out of range")
	}
	return &l.items[idx], nil
}

func (l *List) Append(b models.BookEntry) error {
	if l.m == l.Length() {
		return errors.New("out of space")
	}
	l.items = append(l.items, b)
	return nil
}

func (l *List) Clear() {
	l.items = l.items[:0]
}

func (l *List) Pop() *models.BookEntry {
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

