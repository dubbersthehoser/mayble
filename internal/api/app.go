package app


import (
	"github.com/dubbersthehoser/mayble/internal/storage"
	"github.com/dubbersthehoser/mayble/internal/storage/memeory"
	"github.com/dubbersthehoser/mayble/internal/command"
	"github.com/dubbersthehoser/mayble/internal/data"
	"github.com/dubbersthehoser/mayble/internal/core"
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
	GetBookLoans(*data.BookLoan) ([]data.BookLoan, error)
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
	a.broker = b
	a.storeMgr = newManager(store)
	memStore := memory.NewStorage()
	a.memMgr = newManager(memStore)
	if err := a.load(); err != nil {
		return nil, err
	}
	return &c, nil
}

func (a *App) load() error {
	bookLoans, err := a.storeMgr.store.GetAllBookLoans()
	if err != nil {
		return fmt.Errorf("core: %w", err)
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

func (a *App) GetAllBookLoans() ([]data.BookLoan, error) {
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
	if a.memMgr.undos.length() == 0 {
		return false
	}
	return true
}

func (a *App) Redo() error {
	return a.memMgr.reExecute()
}
func (a *App) RedoIsEmpty() bool {
	if a.memMgr.redos.length() == 0 {
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


/* Create */

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


/* Delete */

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


/* Update */

type commandUpdateBookLoan struct {
	bookLoan *datc.BookLoan
	prevBookLoan *datc.BookLoan
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

func (m *manager) execute(cmd command) error {
	if err := cmd.do(m.store); err != nil {
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
