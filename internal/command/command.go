package command

import (
	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

type Command interface{
	Execute() error
	Undo() error
}





type DeleteBook struct {
	store repo.BookStore
	id    int64
}

func (d *DeleteBook) Execute() error {
	
}


type CreateBook struct {
	store repo.BookStore
	entry *repo.BookEntry
	storedID int64
}

func (c *CreateBook) Execute() error {
	id, err := c.store.CreateBook(c.entry)
	if err != nil {
		return err
	}
	c.storedID = id
	return nil
}

func (c *CreateBook) Undo() error {
	return c.store.DeleteBook(c.storedID)
}

func NewCreateBook(s repo.BookStore, b *repo.BookEntry) *CreateBook {
	return &CreateBook{
		store: s,
		entry: b,
	}
}
