package app

import (
	"fmt"
	"errors"
	"time"
	"github.com/dubbersthehoser/mayble/internal/storage"
)

type BookLoan struct {
	ID       int64
	Title    string
	Author   string
	Genre    string
	Ratting  int
	IsOnLoan   bool
	Borrower string
	Date     time.Time
}

var (
	ErrInvalidValue error = errors.New("invalid value")
)

var ZeroID int64 = storage.ZeroID

type App struct {
	storage  storage.BookLoanStore
	storeMgr *manager 
}

func New(store storage.BookLoanStore) *App {
	var a App
	a.storage = store
	a.storeMgr = newManager()
	return &a
}

func (a *App) GetBookLoans() ([]BookLoan, error) {
	bookLoans, err := getAllBookLoans(a.storage)
	if err != nil {
		return nil, err
	}
	return bookLoans, nil
}

func (a *App) ImportBookLoans(bookLoans []BookLoan) error {
	if bookLoans == nil {
		return fmt.Errorf("%w: nil pointer", ErrInvalidValue)
	}
	cmd := newCommandImportBookLoans(bookLoans)
	return a.storeMgr.execute(cmd(a.storage))
}


func (a *App) CreateBookLoan(book *BookLoan) error {
	if book == nil {
		return fmt.Errorf("%w: nil pointer", ErrInvalidValue)
	}
	cmd := newCommandCreateBookLoan(book)
	return a.storeMgr.execute(cmd(a.storage))
}

func (a *App) UpdateBookLoan(book *BookLoan) error {
	if book == nil {
		return fmt.Errorf("%w: nil pointer", ErrInvalidValue)
	}
	cmd := newCommandUpdateBookLoan(book)
	return a.storeMgr.execute(cmd(a.storage))
}

func (a *App) DeleteBookLoan(book *BookLoan) error {
	if book == nil {
		return fmt.Errorf("%w: nil pointer", ErrInvalidValue)
	}
	cmd := newCommandDeleteBookLoan(book)
	return a.storeMgr.execute(cmd(a.storage))
}




func (a *App) Undo() error {
	if err := a.storeMgr.unExecute(); err != nil {
		return err
	}
	return nil
}

func (a *App) UndoIsEmpty() bool {
	return a.storeMgr.undos.Length() == 0
}

func (a *App) Redo() error {
	if err := a.storeMgr.reExecute(); err != nil {
		return err
	}
	return nil
}

func (a *App) RedoIsEmpty() bool {
	return a.storeMgr.redos.Length() == 0
}

func createBookLoan(s storage.BookLoanStore, bookLoan *BookLoan) (int64, error) {
	id, err := s.CreateBook(bookLoan.ID, bookLoan.Title, bookLoan.Author, bookLoan.Genre, bookLoan.Ratting)
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
			ID: book.ID,
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
	bookLoan.ID = id
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
