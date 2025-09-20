package core

import (
	"fmt"
	"errors"

	"github.com/dubbersthehoser/mayble/internal/storage"
)

type condition int
const (
	Untouched  condition = iota
	Deleted 
	Updated 
	Created
)
func (c condition) String() string {
	switch c {
	case Untouched:
		return "Untouched"
	case Deleted:
		return "Deleted"
	case Updated:
		return "Update"
	case Created:
		return "Created"
	}
	panic("core: condition to string was not found!")
}
type Core struct {
	storage       storage.Storage
	bookLoanTable map[int64]storage.BookLoan
	conditions    map[int64]condition
}
func New(store storage.Storage) (*Core, error) {
	var c Core
	c.storage       = store
	c.bookLoanTable = make(map[int64]storage.BookLoan)
	c.conditions    = make(map[int64]condition)
	if err := c.loadTable(); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Core) getNewID() int64 {
	return int64(len(c.bookLoanTable)) + 1
}

func (c *Core) loadTable() error {
	bookLoans, err := c.storage.GetAllBookLoans()
	if err != nil {
		return fmt.Errorf("core: %w", err)
	}
	for _, book := range bookLoans {
		id := c.getNewID()
		c.bookLoanTable[id] = book
	}
	return nil
}

func (c *Core) saveTable() error {
	for id, book := range c.bookLoanTable {
		condition := c.conditions[id]
		switch condition {
		case Created:
			if err := c.storage.CreateBookLoan(&book); err != nil {
				return err
			}
		case Deleted:
			if err := c.storage.DeleteBookLoan(&book); err != nil {
				return err
			}
		case Updated:
			if err := c.storage.UpdateBookLoan(&book); err != nil {
				return err
			}
		case Untouched:
			continue
		}
	}
	c.bookLoanTable = make(map[int64]storage.BookLoan)
	return c.loadTable()
}

func (c *Core) Save() error {
	return c.saveTable()
}

func (c *Core) ListBookLoanIDs() []int64 {
	ids := make([]int64, len(c.bookLoanTable))
	i := 0
	for id := range c.bookLoanTable {
		ids[i] = id
		i++
	}
	return ids
}
func (c *Core) GetBookLoanByID(id int64) (storage.BookLoan, error) {
	book, ok := c.bookLoanTable[id]
	if !ok {
		return book, errors.New("core: couldn't get book: id not found")
	}
	return book, nil
}
func (c *Core) CreateBookLoan(book storage.BookLoan) (int64, error) {
	id := c.getNewID()
	c.bookLoanTable[id] = book
	c.conditions[id] = Created
	return id, nil
}
func (c *Core) UpdateBookLoan(id int64, book storage.BookLoan) error {
	_, ok := c.bookLoanTable[id]
	if !ok {
		return errors.New("core: couldn't update book id not found")
	}
	c.bookLoanTable[id] = book
	condition := c.conditions[id]
	if condition == Deleted {
		return errors.New("core: couldn't update book is deleted")
	}
	if condition != Created {
		c.conditions[id] = Updated
	}
	return nil
}
func (c *Core) DeleteBookLoan(id int64) error {
	_, ok := c.bookLoanTable[id]
	if !ok {
		return errors.New("core: couldn't delete book id not found")
	}
	c.conditions[id] = Deleted
	return nil
}



