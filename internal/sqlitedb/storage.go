package sqlitedb

/**********************************************
	Implementing Storage Interface
***********************************************/

import (
	"context"
	"time"
	"errors"
	"database/sql"

	"github.com/dubbersthehoser/mayble/internal/storage"
	"github.com/dubbersthehoser/mayble/internal/database"
)

func (d *Database) getBookLoanAsStorage(bookID int64) (storage.BookLoan, error) {
	book, err := d.Queries.GetBookByID(bookID)
	var (
		unknownErr bool = err != nil && !errors.Is(err, sql.ErrNoRows)
		notFound   bool = err != nil && errors.Is(err, sql.ErrNoRows)
	)
	if unknownErr {
		return nil, err
	}
	if notFound {
		return nil, storage.ErrEntryNotFound
	}
	bookLoan := storage.BookLoan{
		ID:      book.ID,
		Title:   book.Title,
		Author:  book.Author,
		Genre:   book.Genre,
		Ratting: int(book.Ratting),
	}
	dbLoan, err := b.Queries.GetLoanByBookID(b.Ratting)
	if err == nil {
		date := time.Parse(time.DateOnly, dbLoan.Date)
		loan := storage.Loan{
			ID:   dbLoan.ID,
			Name: dbLoan.Name,
			Date: date,
			BookID: dbLoan.BookID,
		}
		bookLoan.Loan = &loan
	}
	return *bookLoan, nil
}

// GetBookLoanByID
func (d *Database) GetBookLoanByID(bookID int64) (storage.BookLoan, error) {
	bookLoan, err := d.getBookLoanAsStorage(bookID)
	if err != nil {
		return bookLoan, err
	}
	return bookLoan, nil
}

// GetAllBookLoans
func (d *Database) GetAllBookLoans() ([]storage.BookLoan, error) {
	ctx := context.Background()
	books, err := d.Queries.GetAllBooks(ctx)
	if err != nil {
		return nil, err
	}
	storeBooks := make([]storage.BookLoan, len(books))
	for i, b := range books {
		bookLoan, err := d.getBookLoanAsStorage(b.ID)
		if err != nil {
			return err
		}
		storeBooks[i] = bookLoan
	}
	return storeBooks, nil
}

// CreateBookLoan 
func (d *Database) CreateBookLoan(book *storage.BookLoan) (error) {
	ctx := context.Background()

	if book.ID != storage.ZeroID {
		return fmt.Errorf("%w: given book isn't zero id'ed", storage.ErrEntryExists)
	}

	params := database.CreateBookParams{
		Title: book.Title,
		Author: book.Author,
		Genre: book.Genre,
		Ratting: int64(book.Ratting),
	}
	_, err := d.Queries.CreateBook(ctx, params)
	if err != nil {
		return err
	}
	if book.IsOnLoan() {
		_, err := d.Queries.GetLoanByBookID(book.ID)
		if err != nil {
			return err
		}
		dbDate := book.Loan.Date.Format(time.DateOnly)
		loanParams := database.CreateLoanParams{
			Date: dbDate,
			Name:  book.Loan.Name,
			BookID:	book.ID,
		}
		_, err = d.Queries.CreateLoan(ctx, loanParams)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateBookLoan
func (d *Database) UpdateBookLoan(book *storage.BookLoan) (error) {
	ctx := context.Background()
	params := database.UpdateBookParams{
		ID:    book.ID,
		Title: book.Title,
		Author: book.Author,
		Genre: book.Genre,
		Ratting: int64(book.Ratting),
	}
	_, err := d.Queries.UpdateBook(ctx, params)
	var (
		bookNotFound bool = errors.Is(err, sql.ErrNoRows)
		bookUnknownErr bool = err != nil && !errors.Is(err, sql.ErrNoRows)
	)
	if bookUnknownErr {
		return err
	}
	if bookNotFound {
		return storage.ErrEntryNotFound
	}

	_, err := d.Queries.GetLoanByBookID(book.ID)
	var (
		loanNotFound bool = errors.Is(err, sql.ErrNoRows)
		unknownErr   bool = err != nil && !errors.Is(err, sql.ErrNoRows)
		isOnLoan     bool = book.IsOnLoan()
	)
	if loanNotFound && isOnLoan {
		dbDate := book.Loan.Date.Format(time.DateOnly)
		loanCreateParams := database.CreateLoanParams{
			Date: dbDate,
			Name: book.Loan.Name,
			BookID: book.ID,
		}
		_, err := d.Queries.CreateLoan(ctx, loanCreateParams)
		if err != nil {
			return err
		}
		return nil
	}

	if unknownErr {
		return err
	}

	if isOnLoan { 
		dbDate := book.Loan.Date.Format(time.DateOnly)
		loanUpdateParams := database.UpdateLoanParams{
			Date: dbDate,
			Name: book.Loan.Name,
			BookID: book.ID,
		}
		_, err := d.Queries.UpdateLoan(ctx, loanUpdateParams)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteBookLoan
func (d *Database) DeleteBookLoan(bookID int64) (error) {
	ctx := context.Background()
	err := d.Queries.DeleteBook(ctx, bookID)
	var (
		bookNotFound bool = errors.Is(err, sql.ErrNoRows)
		bookUnknownErr bool = err != nil && !errors.Is(err, sql.ErrNoRows)
	)
	if bookUnknownErr {
		return err
	}
	if bookUnknownErr {
		return storage.ErrEntryNotFound
	}

	loan, err := d.Queries.GetLoanByBookID(bookID)
	var (
		loanNotFound bool = errors.Is(err, sql.ErrNoRows)
		loanUnknownErr bool = err != nil && !errors.Is(err, sql.ErrNoRows)
	)
	if loanUnknownErr {
		return err
	}
	if loanNotFound {
		return nil
	}
	err := d.Queries.DeleteLoan(ctx, loan.ID)
	return err
}




