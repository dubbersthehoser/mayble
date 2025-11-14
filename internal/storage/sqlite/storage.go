package sqlite

import (
	"fmt"
	"context"
	"time"
	"errors"
	"database/sql"

	"github.com/dubbersthehoser/mayble/internal/data"
	"github.com/dubbersthehoser/mayble/internal/storage"
	"github.com/dubbersthehoser/mayble/internal/sqlite"
	"github.com/dubbersthehoser/mayble/internal/sqlite/database"
)

type Storage struct {
	sqlite.Storage
}

func NewStorage() *Storage {
	s := Storage{}
}

// GetBookLoanByID
func (d *Storage) GetBookLoanBy(bookID int64) (data.BookLoan, error) {
	book, err := d.Queries.GetBookByID(context.Background(), bookID)
	var (
		unknownErr bool = err != nil && !errors.Is(err, sql.ErrNoRows)
		notFound   bool = err != nil && errors.Is(err, sql.ErrNoRows)
	)
	if unknownErr {
		return data.BookLoan{}, err
	}
	if notFound {
		return data.BookLoan{}, storage.ErrEntryNotFound
	}
	bookLoan := data.BookLoan{
		Book: storage.Book{
			ID:      book.ID,
			Title:   book.Title,
			Author:  book.Author,
			Genre:   book.Genre,
			Ratting: int(book.Ratting),
		},
	}
	dbLoan, err := d.Queries.GetLoanByBookID(context.Background(), bookLoan.ID)
	if err == nil {
		date, err := time.Parse(time.DateOnly, dbLoan.Date)
		if err != nil {
			return data.BookLoan{}, err
		}
		loan := storage.Loan{
			ID:   dbLoan.ID,
			Name: dbLoan.Name,
			Date: date,
		}
		bookLoan.Loan = &loan
	}
	return bookLoan, nil
}

// GetAllBookLoans
func (d *Storage) GetAllBookLoans() ([]data.BookLoan, error) {
	ctx := context.Background()
	books, err := d.Queries.GetAllBooks(ctx)
	if err != nil {
		return nil, err
	}
	storeBooks := make([]data.BookLoan, len(books))
	for i, b := range books {
		bookLoan, err := d.getBookLoanAsStorage(b.ID)
		if err != nil {
			return nil, err
		}
		storeBooks[i] = bookLoan
	}
	return storeBooks, nil
}

// CreateBookLoan 
func (d *Storage) CreateBookLoan(book *data.BookLoan) (error) {

	if book == nil {
		return fmt.Errorf("%w: book pointer is nil", storage.ErrInvalidValue)
	}

	if book.ID != storage.ZeroID {
		return fmt.Errorf("%w: book id is non zero", storage.ErrInvalidValue)
	}

	ctx := context.Background()

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
		_, err = d.Queries.GetLoanByBookID(ctx, book.ID)
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
func (d *Storage) UpdateBookLoan(book *data.BookLoan) (error) {
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

	_, err = d.Queries.GetLoanByBookID(ctx, book.ID)
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
func (d *Storage) DeleteBookLoan(bookID int64) (error) {

	if bookID == storage.ZeroID {
		return fmt.Errorf("%w: given zero id", storage.ErrInvalidValue)
	}

	ctx := context.Background()
	err := d.Queries.DeleteBook(ctx, bookID)

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

	loan, err := d.Queries.GetLoanByBookID(ctx, bookID)
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
	err = d.Queries.DeleteLoan(ctx, loan.ID)
	return err
}




