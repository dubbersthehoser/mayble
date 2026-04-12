package command

import (
	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

type Command interface{
	Execute() error
	Undo() error
}

//
// Book Storage Commands
//

func NewDeleteBook(id int64, s repo.BookStore, r repo.BookRetriever) *DeleteBook {
	return &DeleteBook{
		store: s,
		retriever: r,
		id: id,
	}
}

func NewUpdateBook(b *repo.BookEntry, s repo.BookStore, r repo.BookRetriever) *UpdateBook {
	return &UpdateBook{
		store: s,
		retriever: r,
		book: b,
	}
}

func NewCreateBook(s repo.BookStore, b *repo.BookEntry) *CreateBook {
	return &CreateBook{
		store: s,
		book: b,
	}
}


type UpdateBook struct {
	store     repo.BookStore
	retriever repo.BookRetriever

	book      *repo.BookEntry
	original  *repo.BookEntry
}


func (u *UpdateBook) Execute() error {
	if u.original == nil {
		old, err := u.retriever.GetBookByID(u.book.ID)
		if err != nil {
			return err
		}
		u.original = &old
	}

	err := u.store.UpdateBook(u.book)
	if err != nil {
		return err
	}
	return nil
}

func (u *UpdateBook) Undo() error {
	if u.original == nil {
		return nil
	}

	err := u.store.UpdateBook(u.original)
	if err != nil {
		return err
	}

	return nil
}


type DeleteBook struct {
	store repo.BookStore
	retriever repo.BookRetriever
	id    int64
	book  *repo.BookEntry
}

func (d *DeleteBook) Execute() error {
	book, err := d.retriever.GetBookByID(d.id)
	if err != nil {
		return err
	}
	d.book = &book
	return d.store.DeleteBook(d.id)
}

func (d *DeleteBook) Undo() error {
	id, err := d.store.CreateBook(d.book)
	if err != nil {
		return err
	}
	d.book.ID = id
	d.id = id
	return nil
}


type CreateBook struct {
	store repo.BookStore
	book *repo.BookEntry
	storedID int64
}

func (c *CreateBook) Execute() error {
	id, err := c.store.CreateBook(c.book)
	if err != nil {
		return err
	}
	c.storedID = id
	return nil
}

func (c *CreateBook) Undo() error {
	return c.store.DeleteBook(c.storedID)
}

