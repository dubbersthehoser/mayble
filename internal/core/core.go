package core

import (
	"strings"
	"slices"
	"fmt"
	"errors"

	"github.com/dubbersthehoser/mayble/internal/storage"
	"github.com/dubbersthehoser/mayble/internal/memdb"
)

/*******************
	Core
********************/

type Core struct {
	storage  storage.Storage
	storeMgr *manager 
	memMgr   *manager 
}
func New(store storage.Storage) (*Core, error) {
	var c Core
	c.storeMgr = newManager(store)
	memStore := memdb.NewMemStorage()
	c.memMgr = newManager(memStore)
	if err := c.load(); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Core) load() error {
	bookLoans, err := c.storeMgr.store.GetAllBookLoans()
	if err != nil {
		return fmt.Errorf("core: %w", err)
	}
	for _, book := range bookLoans {
		err = c.memMgr.store.CreateBookLoan(&book)
		if errors.Is(err, storage.ErrEntryExists) {
			err = c.memMgr.store.UpdateBookLoan(&book)
		}
		if err != nil {
			return err
		}
	}
	return nil
}




/***********************
	Core API
************************/

func (c *Core) Save() error {
	for {
		cmd := c.storeMgr.dequeue()
		if cmd == nil {
			break
		}
		fmt.Printf("%#v\n", cmd)
		err := c.storeMgr.execute(cmd)
		if err != nil {
			return err
		}
	}
	return c.load()
}


/* Listing and Ordering */

type OrderBy string

const (
	ByTitle OrderBy = "Title"
	ByAuthor        = "Author"
	ByGenre         = "Genre"
	ByRatting       = "Ratting"
	ByBorrower      = "Borrower"
	ByDate          = "Date"
)

type Order int

const (
	ASC Order = iota
	DEC
)

func (c *Core) ListBookLoans(by OrderBy, order Order) ([]storage.BookLoan, error) {
	bookLoans, err := c.memMgr.store.GetAllBookLoans()
	if err != nil {
		return nil, err
	}
	compare := func(x, y storage.BookLoan) int {
		const (
			GreaterX int = 1
			Equal    int = 0
			LesserX  int = -1
		)
		result := Equal
		switch by {
		case ByTitle:
			result = strings.Compare(x.Title, y.Title)
		case ByAuthor:
			result = strings.Compare(x.Author, y.Author)
		case ByGenre:
			result = strings.Compare(x.Genre, y.Genre)
		case ByRatting:
			switch {
			case x.Ratting == y.Ratting:
				result = Equal
			case x.Ratting > y.Ratting:
				result = GreaterX
			case x.Ratting < y.Ratting:
				result = LesserX
			}
		case ByBorrower:
			result = strings.Compare(x.Loan.Name, y.Loan.Name)
		case ByDate:
			result = x.Loan.Date.Compare(y.Loan.Date)
		}
		if order == DEC {
			result = result * -1
		}
		return result
	}
	slices.SortFunc(bookLoans, compare)
	return bookLoans, nil
}

func (c *Core) CreateBookLoan(book *storage.BookLoan) error {
	cmd := &commandCreateBookLoan{
		bookLoan: book,
	}
	err := c.memMgr.execute(cmd)
	if err != nil {
		return err
	}
	c.storeMgr.enqueue(cmd)
	return nil
}

func (c *Core) UpdateBookLoan(book *storage.BookLoan) error {
	cmd := &commandUpdateBookLoan{
		bookLoan: book,
	}
	err := c.memMgr.execute(cmd)
	if err != nil {
		return err
	}
	c.storeMgr.enqueue(cmd)
	return nil
}

func (c *Core) DeleteBookLoan(book *storage.BookLoan) error {
	cmd := &commandDeleteBookLoan{
		bookLoan: book,
	}
	err := c.memMgr.execute(cmd)
	if err != nil {
		return err
	}
	c.storeMgr.enqueue(cmd)
	return nil
}




/********************************
	Command Stack
*********************************/

type commandStack struct {
	items []command
}
func newCommandStack() *commandStack {
	c := commandStack{
		items: make([]command, 0),
	}
	return &c
}

func (cs *commandStack) pop() command {
	length := len(cs.items)
	cmd := cs.items[length-1]
	cs.items = cs.items[:length]
	return cmd
}

func (cs *commandStack) push(cmd command) {
	cs.items = append(cs.items, cmd)
}

func (cs *commandStack) length() int {
	return len(cs.items)
}

func (cs *commandStack) clear() {
	cs.items = make([]command, 0)
}





/*******************************
	Command Manager
********************************/

type manager struct {
	store storage.Storage
	undos *commandStack
	redos *commandStack
	queue []command
}

func newManager(store storage.Storage) *manager{
	m := manager{
		undos: newCommandStack(),
		redos: newCommandStack(),
		store: store,
	}
	return &m
}

func (m *manager) execute(cmd command) error {
	if err := cmd.do(m.store); err != nil {
		return err
	}
	m.undos.push(cmd)
	m.redos.clear()
	return nil
}

func (m *manager) unExecute() error {
	cmd := m.undos.pop()
	err := cmd.undo(m.store)
	if err != nil {
		return err
	}
	m.redos.push(cmd)
	return nil
} 

func (m *manager) reExecute() error {
	cmd := m.redos.pop()
	if cmd == nil {
		return nil
	}
	err := cmd.do(m.store)
	if err != nil {
		return err
	}
	m.undos.push(cmd)
	return nil
}

func (m *manager) enqueue(cmd command) {
	m.queue = append(m.queue, cmd)
}

func (m *manager) dequeue() command {
	if len(m.queue) == 0 {
		return nil
	}
	cmd := m.queue[0]
	m.queue = m.queue[1:]
	return cmd
}






/**********************
	Commands
***********************/

type command interface {
	do(storage.Storage)   error
	undo(storage.Storage) error
}

/* Create Book Loan */

type commandCreateBookLoan struct {
	bookLoan *storage.BookLoan
}

func (c *commandCreateBookLoan) do(s storage.Storage) error {
	return s.CreateBookLoan(c.bookLoan)
}

func (c *commandCreateBookLoan) undo(s storage.Storage) error {
	return s.DeleteBookLoan(c.bookLoan)
}


/* Delete Book Loan */

type commandDeleteBookLoan struct {
	bookLoan *storage.BookLoan
}

func (c *commandDeleteBookLoan) do(s storage.Storage) error {
	return s.DeleteBookLoan(c.bookLoan)
}

func (c *commandDeleteBookLoan) undo(s storage.Storage) error {
	return s.CreateBookLoan(c.bookLoan)
}


/* Update Book Loan */

type commandUpdateBookLoan struct {
	bookLoan *storage.BookLoan
	prevBookLoan *storage.BookLoan
}

func (c *commandUpdateBookLoan) do(s storage.Storage) error {
	book, err := s.GetBookLoanByID(c.bookLoan.ID)
	if err != nil {
		return err
	}
	c.prevBookLoan = &book
	return s.UpdateBookLoan(c.bookLoan)
}

func (c *commandUpdateBookLoan) undo(s storage.Storage) error {
	book := c.prevBookLoan
	c.prevBookLoan = c.bookLoan
	c.bookLoan = book
	return s.UpdateBookLoan(book)
}
