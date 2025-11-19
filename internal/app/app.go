package app

import (
	"fmt"
	"errors"
	"github.com/dubbersthehoser/mayble/internal/storage"
	"github.com/dubbersthehoser/mayble/internal/storage/memory"
	"github.com/dubbersthehoser/mayble/internal/command"
)


type App struct {
	storage  storage.Storage
	storeMgr *manager 
	memMgr   *manager 
}

func New(store storage.BookLoanStore) (*App, error) {
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

func (a *App) ImportBookLoans(bookLoans []BookLoan) error {
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


func (a *App) CreateBookLoan(book *BookLoan) error {
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

func (a *App) UpdateBookLoan(book *BookLoan) error {
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

func (a *App) DeleteBookLoan(book *BookLoan) error {
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

func createBookLoan(s storage.BookLoanStore, bookLoan *BookLoan) (int64, error) {
	id, err := s.CreateBook(bookLoan.Title, bookLoan.Author, bookLoan.Genre, bookLoan.Ratting)
	if err != nil {
		return -1, err
	}
	if !bookLoan.IsOnLoan{
		return id, nil
	}
	err := s.CreateLoan(id, bookLoan.Borrower, bookLoan.Date)
	if err != nil {
		return -1, err
	}
	return id, nil
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
	}
	else if err != nil {
		return nil, err
	}
	bookLoan.Borrower = loan.Borrower
	bookLoan.Date = loan.Date
	return &bookLoan, nil
}

func updateBookLoan(s storage.BookLoanStore, bookLoan *BookLoan) error {
	book, err := s.GetBookByID(bookLoan.ID)
	if err != nil {
		return err
	}
	book.Title = bookLoan.Title
	book.Author = bookLoan.Author
	book.Genre = bookLoan.Genre
	book.Ratting = bookLoan.Ratting
	err = s.UpdateBook(&book)
	if err != nil {
		return err
	}

	var (
		LoanIsInStore bool
	)
	_, err := s.GetLoan(bookLoan.ID)
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
		Date: bookLoan.Date
	}
	if LoanIsInStore {
		if bookLoan.IsOnLoan {
			err := s.UpdateLoan(&loan)
			if err != nil {
				return err
			}
		} else {
			err := s.DeleteLoan(&loan)
			if err != nil {
				return err
			}
		}
		return nil
	} else { 
		return s.CreateLoan(loan.ID, loan.Borrower, loan.Date)
	}
}

func deleteBookLoan(s storage.BookLoanStore, bookLoan *BookLoan) error {
	book := storage.Book{
		ID: bookLoan.ID,
	}
	err := s.DeleteBook(&book)
	if err != nil {
		return err
	}
	loan, err : s.GetLoan(book.ID)
	if errors.Is(err, storage.EntryNotFound) {
		return
	} else if (err != nil ) {
		return err
	}
	loan := storage.Loan{
		ID: bookLoan.ID,
	}
	return s.DeleteLoan(&loan)
}
