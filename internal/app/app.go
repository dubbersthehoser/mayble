package app

import (
	"fmt"
	"errors"
	"github.com/dubbersthehoser/mayble/internal/storage"
	"github.com/dubbersthehoser/mayble/internal/storage/memory"
)


type App struct {
	storage  storage.BookLoanStore
	storeMgr *manager 

	memory   *memory.Storage
	memMgr   *manager 
}

func New(store storage.BookLoanStore) (*App, error) {
	var a App
	a.storage = store
	a.storeMgr = newManager()
	a.memory = memory.NewStorage()
	a.memMgr = newManager()
	if err := a.load(); err != nil {
		return nil, err
	}
	return &a, nil
}

func (a *App) load() error {
	bookLoans, err := getAllBookLoans(a.storage)
	if err != nil {
		return fmt.Errorf("app: %w", err)
	}
	for _, book := range bookLoans {
		_, err = createBookLoan(a.memory, &book)
		if errors.Is(err, storage.ErrEntryExists) {
			err = updateBookLoan(a.memory, &book)
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
	a.storeMgr.undos.Clear() // clear undos from executed commands
	return a.load()
}

func (a *App) GetBookLoans() ([]BookLoan, error) {
	bookLoans, err := getAllBookLoans(a.memory)
	if err != nil {
		return nil, err
	}
	return bookLoans, nil
}

func (a *App) ImportBookLoans(bookLoans []BookLoan) error {
	cmd := newCommandImportBookLoans(bookLoans)
	err := a.memMgr.execute(cmd(a.memory))
	if err != nil {
		return err
	}
	a.storeMgr.enqueue(cmd(a.storage))
	return nil
}


func (a *App) CreateBookLoan(book *BookLoan) error {
	cmd := newCommandCreateBookLoan(book)
	err := a.memMgr.execute(cmd(a.memory))
	if err != nil {
		return err
	}
	a.storeMgr.enqueue(cmd(a.storage))
	return nil
}

func (a *App) UpdateBookLoan(book *BookLoan) error {
	cmd := newCommandUpdateBookLoan(book)
	err := a.memMgr.execute(cmd(a.memory))
	if err != nil {
		return err
	}
	a.storeMgr.enqueue(cmd(a.storage))
	return nil
}

func (a *App) DeleteBookLoan(book *BookLoan) error {
	cmd := newCommandDeleteBookLoan(book)
	err := a.memMgr.execute(cmd(a.memory))
	if err != nil {
		return err
	}
	a.storeMgr.enqueue(cmd(a.storage))
	return nil
}

func (a *App) Undo() error {
	return a.memMgr.unExecute()
}
func (a *App) UndoIsEmpty() bool {
	return a.memMgr.undos.Length() == 0
}

func (a *App) Redo() error {
	return a.memMgr.reExecute()
}
func (a *App) RedoIsEmpty() bool {
	return a.memMgr.redos.Length() == 0
}


func createBookLoan(s storage.BookLoanStore, bookLoan *BookLoan) (int64, error) {
	id, err := s.CreateBook(bookLoan.Title, bookLoan.Author, bookLoan.Genre, bookLoan.Ratting)
	if err != nil {
		return -1, err
	}
	if !bookLoan.IsOnLoan {
		return id, nil
	}
	err = s.CreateLoan(id, bookLoan.Borrower, bookLoan.Date)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func getAllBookLoans(s storage.BookLoanStore) ([]BookLoan, error) {
	books, err := s.GetBooks()
	if err != nil {
		return nil, err
	}

	bookLoans := make([]BookLoan, 0)
	for _, book := range books {
		bookLoan := BookLoan{
			Title: book.Title,
			Author: book.Author,
			Genre: book.Genre,
			Ratting: book.Ratting,
		}

		loan, err := s.GetLoan(book.ID)
		if errors.Is(err, storage.ErrEntryNotFound) {
			bookLoans = append(bookLoans, bookLoan)
			continue
		}
		if err == nil {
			bookLoan.Borrower = loan.Borrower
			bookLoan.Date = loan.Date
			bookLoan.IsOnLoan = true
			bookLoans = append(bookLoans, bookLoan)
			continue
		}
		return nil, err
	}
	return bookLoans, nil
}

func getBookLoanByID(s storage.BookLoanStore, id int64) (*BookLoan, error) {
	bookLoan := BookLoan{}
	book, err := s.GetBookByID(id)
	if err != nil {
		return nil, err
	}
	bookLoan.Title = book.Title
	bookLoan.Author = book.Author
	bookLoan.Genre = book.Genre
	bookLoan.Ratting = book.Ratting
	bookLoan.IsOnLoan = false

	loan, err := s.GetLoan(id)
	if errors.Is(err, storage.ErrEntryNotFound) {
		return &bookLoan, nil
	} else if err != nil {
		return nil, err
	}
	bookLoan.IsOnLoan = true
	bookLoan.Borrower = loan.Borrower
	bookLoan.Date = loan.Date
	return &bookLoan, nil
}

func updateBookLoan(s storage.BookLoanStore, bookLoan *BookLoan) error {
	book, err := s.GetBookByID(bookLoan.ID)
	if err != nil {
		return err
	}
	book.ID = bookLoan.ID
	book.Title = bookLoan.Title
	book.Author = bookLoan.Author
	book.Genre = bookLoan.Genre
	book.Ratting = bookLoan.Ratting
	err = s.UpdateBook(book)
	if err != nil {
		return err
	}

	var (
		LoanIsInStore bool
	)
	_, err = s.GetLoan(bookLoan.ID)
	if errors.Is(err, storage.ErrEntryNotFound) {
		LoanIsInStore = false
	} else if err != nil {
		return err
	} else {
		LoanIsInStore = true
	}

	loan := storage.Loan{
		ID: bookLoan.ID,
		Borrower: bookLoan.Borrower,
		Date: bookLoan.Date,
	}

	if !LoanIsInStore && !bookLoan.IsOnLoan {
		return nil
	}

	if LoanIsInStore && bookLoan.IsOnLoan {
		err := s.UpdateLoan(&loan)
		if err != nil {
			return err
		}
		return nil
	}
	if LoanIsInStore && !bookLoan.IsOnLoan {
		err := s.DeleteLoan(&loan)
		if err != nil {
			return err
		}
		return nil
	}
	return s.CreateLoan(loan.ID, loan.Borrower, loan.Date)
}

func deleteBookLoan(s storage.BookLoanStore, bookLoan *BookLoan) error {
	book := storage.Book{
		ID: bookLoan.ID,
	}
	err := s.DeleteBook(&book)
	if err != nil {
		return err
	}
	loan, err := s.GetLoan(book.ID)
	if errors.Is(err, storage.ErrEntryNotFound) {
		return nil
	} else if (err != nil ) {
		return err
	}
	return s.DeleteLoan(loan)
}









