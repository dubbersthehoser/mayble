package app

import (
	"fmt"

	"github.com/dubbersthehoser/mayble/internal/command"
	"github.com/dubbersthehoser/mayble/internal/storage"
)

/**********************
        Commands
***********************/

/* Import Book Loans */

// NOTE need to re-think import system implementation.

type commandImportBookLoans struct {
	store     storage.BookLoanStore
	addedIDs  []int64   // ids map to bookLoans slice.
	bookLoans []BookLoan
}
func newCommandImportBookLoans(books []BookLoan) func(storage.BookLoanStore) *commandImportBookLoans {
	return func(s storage.BookLoanStore) *commandImportBookLoans {
		return &commandImportBookLoans{
			bookLoans: books,
			store: s,
		}
	}
}
func (c *commandImportBookLoans) Do() error {
	c.addedIDs = make([]int64, len(c.bookLoans))
	for i, BookLoan := range c.bookLoans {
		id, err := createBookLoan(c.store, &BookLoan)
		if err != nil {
			return fmt.Errorf("app: %w", err)
		}
		c.addedIDs[i] = id
	}
	return nil
}
func (c *commandImportBookLoans) Undo() error {
	for i, id := range c.addedIDs {
		book := c.bookLoans[i]
		book.ID = id
		err := deleteBookLoan(c.store, &book)
		if err != nil {
			return err
		}
	}
	return nil
}



/* Create Book Loan */

type commandCreateBookLoan struct {
	store    storage.BookLoanStore
	bookLoan *BookLoan
}
func newCommandCreateBookLoan(book *BookLoan) func(storage.BookLoanStore) *commandCreateBookLoan {
	return func(s storage.BookLoanStore) *commandCreateBookLoan {
		return &commandCreateBookLoan{
			bookLoan: book,
			store: s,
		}
	}
}

func (c *commandCreateBookLoan) Do() error {
	fmt.Println(c.bookLoan)
	id, err := createBookLoan(c.store, c.bookLoan)
	if err != nil {
		return err
	}
	c.bookLoan.ID = id
	return nil
}

func (c *commandCreateBookLoan) Undo() error {
	return deleteBookLoan(c.store, c.bookLoan)
}


/* Delete Book Loan */

type commandDeleteBookLoan struct {
	store    storage.BookLoanStore
	bookLoan *BookLoan
}
func newCommandDeleteBookLoan(book *BookLoan) func(storage.BookLoanStore) *commandDeleteBookLoan {
	return func(s storage.BookLoanStore) *commandDeleteBookLoan {
		return &commandDeleteBookLoan{
			bookLoan: book,
			store: s,
		}
	}
}

func (c *commandDeleteBookLoan) Do() error {
	return deleteBookLoan(c.store, c.bookLoan)
}

func (c *commandDeleteBookLoan) Undo() error {
	_, err := createBookLoan(c.store, c.bookLoan)
	return err
}


/* Update Book Loan */

type commandUpdateBookLoan struct {
	store        storage.BookLoanStore
	bookLoan     *BookLoan
	prevBookLoan *BookLoan
}

func newCommandUpdateBookLoan(book *BookLoan) func(storage.BookLoanStore) *commandUpdateBookLoan {
	return func(s storage.BookLoanStore) *commandUpdateBookLoan {
		return &commandUpdateBookLoan{
			bookLoan: book,
			store: s,
		}
	}
}

func (c *commandUpdateBookLoan) Do() error {
	bookLoan, err := getBookLoanByID(c.store, c.bookLoan.ID)
	if err != nil {
		return err
	}
	if c.prevBookLoan == nil {
		c.prevBookLoan = bookLoan
	} else {
		bookLoan = c.prevBookLoan
		c.prevBookLoan = c.bookLoan
		c.bookLoan = bookLoan
	}
	return updateBookLoan(c.store, c.bookLoan)
}

func (c *commandUpdateBookLoan) Undo() error {
	book := c.prevBookLoan
	c.prevBookLoan = c.bookLoan
	c.bookLoan = book
	return updateBookLoan(c.store, c.bookLoan)
}



/****************************************
        Command Storage Manager
*****************************************/

// manager of commands and invoking.
type manager struct {
	undos *command.Stack
	redos *command.Stack
	queue *command.Queue
}

func newManager() *manager{
	m := manager{
		undos: command.NewStack(),
		redos: command.NewStack(),
		queue: command.NewQueue(),
	}
	return &m
}

// execute command
func (m *manager) execute(cmd command.Command) error {
	if err := cmd.Do(); err != nil {
		return err
	}
	m.undos.Push(cmd)
	m.redos.Clear()
	return nil
}

// unExecute command
func (m *manager) unExecute() error {
	cmd := m.undos.Pop()
	if cmd == nil {
		return nil
	}
	err := cmd.Undo()
	if err != nil {
		return err
	}
	m.redos.Push(cmd)
	return nil
} 

// reExecute an undo'ed command
func (m *manager) reExecute() error {
	cmd := m.redos.Pop()
	if cmd == nil {
		return nil
	}
	err := cmd.Do()
	if err != nil {
		return err
	}
	m.undos.Push(cmd)
	return nil
}

// enqueue command into queue
func (m *manager) enqueue(cmd command.Command) {
	m.queue.Enqueue(cmd)
}

// dequeue command out of queue
func (m *manager) dequeue() command.Command {
	return m.queue.Dequeue()
}
