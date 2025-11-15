package core

import (
	"strings"
	"slices"
	"fmt"
	"errors"

	"github.com/dubbersthehoser/mayble/internal/storage"
	"github.com/dubbersthehoser/mayble/internal/memdb"
	"github.com/dubbersthehoser/mayble/internal/broker"
)

/*******************
	Core
********************/

type Core struct {
	storage  storage.Storage
	storeMgr *manager 
	memMgr   *manager 
	broker   *broker.Broker
}
func New(store storage.Storage, b *broker.Broker) (*Core, error) {
	var c Core
	c.broker = b
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
		_, err = c.memMgr.store.CreateBookLoan(&book)
		if errors.Is(err, storage.ErrEntryExists) {
			err = c.memMgr.store.UpdateBookLoan(&book)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func HandleBookLoanCreate(c *Core) broker.Handler {
	handler := func(d any) {
		bookLoan, ok := d.(data.BookLoan)
		if !ok {
			panic("core: create bookloan was given invalid data")
		}
		c.CreateBookLoan(&bookLoan)
	}
	return broker.Handler(handler)
}

func handleBookLoanDelete(c *Core) broker.Handler

func (c *Core) setUpListeners() {
	b.On(KeyBookLoanCreate, HandleBookCreate(c))

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
		err := c.storeMgr.execute(cmd)
		if err != nil {
			return err
		}
	}
	return c.load()
}

func (c *Core) GetAllBookLoans() ([]data.BookLoan, error) {
	bookLoans, err := c.memMgr.store.GetAllBookLoans()
	if err != nil {
		return nil, err
	}
	return bookLoans, nil
}

/* Listing and Ordering */

//type OrderBy string
//
//const (
//	ByTitle OrderBy = "Title"
//	ByAuthor        = "Author"
//	ByGenre         = "Genre"
//	ByRatting       = "Ratting"
//	ByBorrower      = "Borrower"
//	ByDate          = "Date"
//	ByID            = "ID"
//)
//
//type Order int
//
//const (
//	ASC Order = iota
//	DEC
//)
//
//
//func (c *Core) ListBookLoans(by OrderBy, order Order) ([]data.BookLoan, error) {
//	
//	bookLoans, err := c.memMgr.store.GetAllBookLoans()
//	if err != nil {
//		return nil, err
//	}
//
//	compare := func(x, y storage.BookLoan) int {
//		const (
//			GreaterX int = 1
//			Equal    int = 0
//			LesserX  int = -1
//		)
//		result := Equal
//		switch by {
//		case ByID:
//			switch {
//			case x.Ratting == y.Ratting:
//				result = Equal
//			case x.Ratting > y.Ratting:
//				result = GreaterX
//			case x.Ratting < y.Ratting:
//				result = LesserX
//			}
//		case ByTitle:
//			a := strings.ToLower(x.Title)
//			b := strings.ToLower(y.Title)
//			result = strings.Compare(a, b)
//		case ByAuthor:
//			a := strings.ToLower(x.Author)
//			b := strings.ToLower(y.Author)
//			result = strings.Compare(a, b)
//		case ByGenre:
//			a := strings.ToLower(x.Genre)
//			b := strings.ToLower(y.Genre)
//			result = strings.Compare(a, b)
//		case ByRatting:
//			switch {
//			case x.Ratting == y.Ratting:
//				result = Equal
//			case x.Ratting > y.Ratting:
//				result = GreaterX
//			case x.Ratting < y.Ratting:
//				result = LesserX
//			}
//		case ByBorrower, ByDate:
//			if x.Loan == nil && y.Loan == nil {
//				result = Equal
//			} else if x.Loan == nil {
//				result = LesserX
//			} else if y.Loan == nil {
//				result = GreaterX
//			} else if by == ByBorrower {
//				a := strings.ToLower(x.Loan.Name)
//				b := strings.ToLower(y.Loan.Name)
//				result = strings.Compare(a, b)
//			} else if by == ByDate {
//				result = x.Loan.Date.Compare(y.Loan.Date)
//			}
//		}
//		if order == DEC {
//			result = result * -1
//		}
//		return result
//	}
//	slices.SortFunc(bookLoans, compare)
//	return bookLoans, nil
//}


func (c *Core) ImportBookLoans(bookLoans []data.BookLoan) error {
	cmd := &commandImportBookLoans{
		bookLoans: bookLoans,
	}
	err := c.memMgr.execute(cmd)
	if err != nil {
		return err
	}
	c.storeMgr.enqueue(cmd)
	return nil
}


func (c *Core) CreateBookLoan(book *data.BookLoan) error {
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

func (c *Core) UpdateBookLoan(book *data.BookLoan) error {
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

func (c *Core) DeleteBookLoan(book *data.BookLoan) error {
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

func (c *Core) Undo() error {
	return c.memMgr.unExecute()
}
func (c *Core) IsUndo() bool {
	if c.memMgr.undos.length() == 0 {
		return false
	}
	return true
}

func (c *Core) Redo() error {
	return c.memMgr.reExecute()
}
func (c *Core) IsRedo() bool {
	if c.memMgr.redos.length() == 0 {
		return false
	}
	return true
	
}





/**********************
	Commands
***********************/

type command interface {
	do(storage.Storage)   error
	undo(storage.Storage) error
}

/* Import Book Loans */

type commandImportBookLoans struct {
	addedIDs []int64
	bookLoans []book.BookLoan
}
func (c *commandImportBookLoans) do(s storage.Storage) error {
	c.addedIDs = make([]int64, len(c.bookLoans))
	for i, BookLoan := range c.bookLoans {
		id, err := s.CreateBookLoan(&BookLoan)
		if err != nil {
			return fmt.Errorf("core: import: %w", err)
		}
		c.addedIDs[i] = id
	}
	return nil
	
}
func (c *commandImportBookLoans) undo(s storage.Storage) error {
	for _, id := range c.addedIDs {
		book, err :=s.GetBookLoanByID(id)
		if err != nil {
			return err
		}
		err = s.DeleteBookLoan(&book)
		if err != nil {
			return err
		}
	}
	return nil
}

/* Create Book Loan */

type commandCreateBookLoan struct {
	bookLoan *storage.BookLoan
}

func (c *commandCreateBookLoan) do(s storage.Storage) error {
	_, err := s.CreateBookLoan(c.bookLoan)
	return err
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
	_, err := s.CreateBookLoan(c.bookLoan)
	return err
}


/* Update Book Loan */

type commandUpdateBookLoan struct {
	bookLoan *data.BookLoan
	prevBookLoan *data.BookLoan
}

func (c *commandUpdateBookLoan) do(s storage.Storage) error {
	book, err := s.GetBookLoanByID(c.bookLoan.ID)
	if err != nil {
		return err
	}
	if c.prevBookLoan == nil {
		c.prevBookLoan = &book
	} else {
		book = *c.prevBookLoan
		c.prevBookLoan = c.bookLoan
		c.bookLoan = &book
	}
	return s.UpdateBookLoan(c.bookLoan)
}

func (c *commandUpdateBookLoan) undo(s storage.Storage) error {
	book := c.prevBookLoan
	c.prevBookLoan = c.bookLoan
	c.bookLoan = book
	return s.UpdateBookLoan(c.bookLoan)
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
	if cmd == nil {
		return nil
	}
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
	if length == 0 {
		return nil
	}
	cmd := cs.items[length-1]
	cs.items = cs.items[:length-1]
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



