package repository

import (
	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/models"
)

type BookStoreNotify struct {
	store BookStore
	bus   *bus.Bus
}

func NewBookStoreNotify(s BookStore, b *bus.Bus) *BookStoreNotify {
	return &BookStoreNotify{
		store: s,
		bus: b,
	}
}

func (bn *BookStoreNotify) CreateBook(b *models.BookEntry) (int64, error) {
	id, err := bn.store.CreateBook(b)
	// notify
	return id, err
}

func (bn *BookStoreNotify) UpdateBook(b *models.BookEntry) (error) {
	err := bn.store.UpdateBook(b)
	// notify
	return err
}

func (bn *BookStoreNotify) DeleteBook(id int64) (error) {
	err := bn.store.DeleteBook(id)
	// notify
	return err
}
