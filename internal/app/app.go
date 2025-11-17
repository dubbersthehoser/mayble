package app

import (
	"fmt"
	"errors"
	"github.com/dubbersthehoser/mayble/internal/storage"
	"github.com/dubbersthehoser/mayble/internal/storage/memory"
	"github.com/dubbersthehoser/mayble/internal/command"
	"github.com/dubbersthehoser/mayble/internal/data"
)

type Redoable interface {
	Redo() error
	RedoIsEmpty()  bool
}
type Undoable interface {
	Undo() error
	UndoIsEmpty() bool
}
type Savable interface {
	Save() error
}

type BookLoaning interface {
	CreateBookLoan(*data.BookLoan) error
	UpdateBookLoan(*data.BookLoan) error
	DeleteBookLoan(*data.BookLoan) error
	GetBookLoans() ([]data.BookLoan, error)
	ImportBookLoans([]data.BookLoan) error
}

type Mayble interface {
	BookLoaning
	Redoable
	Undoable
	Savable
}

type App struct {
	storage  storage.Storage
	storeMgr *manager 
	memMgr   *manager 
}
func New(store storage.Storage) (*App, error) {
	var a App
	a.storeMgr = newManager(store)
	memStore := memory.NewStorage()
	a.memMgr = newManager(memStore)
	if err := a.load(); err != nil {
		return nil, err
	}
	return &a, nil
}

func (a *App) load() error {
	bookLoans, err := a.storeMgr.store.GetAllBookLoans()
	if err != nil {
		return fmt.Errorf("app: %w", err)
	}
	for _, book := range bookLoans {
		_, err = a.memMgr.store.CreateBookLoan(&book)
		if errors.Is(err, storage.ErrEntryExists) {
			err = a.memMgr.store.UpdateBookLoan(&book)
		}
		if err != nil {
			return err
		}
	}
	return nil
}


func (a *App) Save() error {
	for {
		cmd := a.storeMgr.dequeue()
		if cmd == nil {
			break
		}
		err := a.storeMgr.execute(cmd)
		if err != nil {
			return err
		}
	}
	return a.load()
}

func (a *App) GetBookLoans() ([]data.BookLoan, error) {
	bookLoans, err := a.memMgr.store.GetAllBookLoans()
	if err != nil {
		return nil, err
	}
	return bookLoans, nil
}

func (a *App) ImportBookLoans(bookLoans []data.BookLoan) error {
	cmd := &commandImportBookLoans{
		bookLoans: bookLoans,
	}
	err := a.memMgr.execute(cmd)
	if err != nil {
		return err
	}
	a.storeMgr.enqueue(cmd)
	return nil
}


func (a *App) CreateBookLoan(book *data.BookLoan) error {
	cmd := &commandCreateBookLoan{
		bookLoan: book,
	}
	err := a.memMgr.execute(cmd)
	if err != nil {
		return err
	}
	a.storeMgr.enqueue(cmd)
	return nil
}

func (a *App) UpdateBookLoan(book *data.BookLoan) error {
	cmd := &commandUpdateBookLoan{
		bookLoan: book,
	}
	err := a.memMgr.execute(cmd)
	if err != nil {
		return err
	}
	a.storeMgr.enqueue(cmd)
	return nil
}

func (a *App) DeleteBookLoan(book *data.BookLoan) error {
	cmd := &commandDeleteBookLoan{
		bookLoan: book,
	}
	err := a.memMgr.execute(cmd)
	if err != nil {
		return err
	}
	a.storeMgr.enqueue(cmd)
	return nil
}

func (a *App) Undo() error {
	return a.memMgr.unExecute()
}
func (a *App) UndoIsEmpty() bool {
	if a.memMgr.undos.Length() == 0 {
		return false
	}
	return true
}

func (a *App) Redo() error {
	return a.memMgr.reExecute()
}
func (a *App) RedoIsEmpty() bool {
	if a.memMgr.redos.Length() == 0 {
		return false
	}
	return true
	
}





/**********************
	Commands
***********************/

/* Import */

type commandImportBookLoans struct {
	addedIDs []int64
	bookLoans []data.BookLoan
}
func (c *commandImportBookLoans) Do(s storage.Storage) error {
	c.addedIDs = make([]int64, len(c.bookLoans))
	for i, BookLoan := range c.bookLoans {
		id, err := s.CreateBookLoan(&BookLoan)
		if err != nil {
			return fmt.Errorf("app: import: %w", err)
		}
		c.addedIDs[i] = id
	}
	return nil
	
}
func (c *commandImportBookLoans) Undo(s storage.Storage) error {
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


/* Create */

type commandCreateBookLoan struct {
	bookLoan *data.BookLoan
}

func (c *commandCreateBookLoan) Do(s storage.Storage) error {
	_, err := s.CreateBookLoan(c.bookLoan)
	return err
}

func (c *commandCreateBookLoan) Undo(s storage.Storage) error {
	return s.DeleteBookLoan(c.bookLoan)
}


/* Delete */

type commandDeleteBookLoan struct {
	bookLoan *data.BookLoan
}

func (c *commandDeleteBookLoan) Do(s storage.Storage) error {
	return s.DeleteBookLoan(c.bookLoan)
}

func (c *commandDeleteBookLoan) Undo(s storage.Storage) error {
	_, err := s.CreateBookLoan(c.bookLoan)
	return err
}


/* Update */

type commandUpdateBookLoan struct {
	bookLoan *data.BookLoan
	prevBookLoan *data.BookLoan
}

func (c *commandUpdateBookLoan) Do(s storage.Storage) error {
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

func (c *commandUpdateBookLoan) Undo(s storage.Storage) error {
	book := c.prevBookLoan
	c.prevBookLoan = c.bookLoan
	c.bookLoan = book
	return s.UpdateBookLoan(c.bookLoan)
}





/****************************************
        Command Storage Manager
*****************************************/

type manager struct {
	store storage.Storage
	undos *command.Stack
	redos *command.Stack
	queue []command.Command
}

func newManager(store storage.Storage) *manager{
	m := manager{
		undos: command.NewStack(),
		redos: command.NewStack(),
		store: store,
	}
	return &m
}

func (m *manager) execute(cmd command.Command) error {
	if err := cmd.Do(m.store); err != nil {
		return err
	}
	m.undos.Push(cmd)
	m.redos.Clear()
	return nil
}

func (m *manager) unExecute() error {
	cmd := m.undos.Pop()
	if cmd == nil {
		return nil
	}
	err := cmd.Undo(m.store)
	if err != nil {
		return err
	}
	m.redos.Push(cmd)
	return nil
} 

func (m *manager) reExecute() error {
	cmd := m.redos.Pop()
	if cmd == nil {
		return nil
	}
	err := cmd.Do(m.store)
	if err != nil {
		return err
	}
	m.undos.Push(cmd)
	return nil
}

func (m *manager) enqueue(cmd command.Command) {
	m.queue = append(m.queue, cmd)
}

func (m *manager) dequeue() command.Command {
	if len(m.queue) == 0 {
		return nil
	}
	cmd := m.queue[0]
	m.queue = m.queue[1:]
	return cmd
}
