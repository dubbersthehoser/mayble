package sqlitedb


// Implementing the storage interface for sqlite.

import (
	"context"
	"time"

	"github.com/dubbersthehoser/mayble/internal/storage"
	"github.com/dubbersthehoser/mayble/internal/database"
)
/*
	Implement storage.BookStore Interface
*/


// GetAllBooks returns all books from database.
func (d *Database) GetAllBooks() ([]storage.Book, error) {
	ctx := context.Background()
	books, err := d.Queries.GetAllBooks(ctx)
	if err != nil {
		return nil, fmtError(err)
	}

	storeBooks := make([]storage.Book, len(books))
	for i, b := range books {
		sBook := storage.Book{
			ID:      b.ID,
			Title:   b.Title,
			Author:  b.Author,
			Genre:   b.Genre,
			Ratting: int(b.Ratting),
		}
		storeBooks[i] = sBook
	}
	return storeBooks, nil
}

// createBookWithRet used for testing and is the main logic for CreateBook().
func (d *Database) createBookWithRet(book *storage.Book) (database.Book, error) {
	ctx := context.Background()
	params := database.CreateBookParams{
		Title: book.Title,
		Author: book.Author,
		Genre: book.Genre,
		Ratting: int64(book.Ratting),
	}
	return d.Queries.CreateBook(ctx, params)
}

// CreateBook add given book to database.
func (d *Database) CreateBook(book *storage.Book) error {
	_, err := d.createBookWithRet(book)
	return err
}

// updateBookWithRet used for testing, and is the main logic of UpdateBook().
func (d *Database) updateBookWithRet(book *storage.Book) (*database.Book, error) {
	ctx := context.Background()
	params := database.UpdateBookParams{
		ID:    book.ID,
		Title: book.Title,
		Author: book.Author,
		Genre: book.Genre,
		Ratting: int64(book.Ratting),
	}
	b, err := d.Queries.UpdateBook(ctx, params)
	return &b, err
}

// UpdateBook update given book in the database.
func (d *Database) UpdateBook(book *storage.Book) error {
	_, err := d.updateBookWithRet(book)
	return err
}

// DeleteBook remove given book from database.
func (d *Database) DeleteBook(book *storage.Book) error {
	ctx := context.Background()
	return d.Queries.DeleteBook(ctx, book.ID)
}


/*
	Implement storage.LoanStore Interface
*/


// GetAllLoans returns all loans from database.
func (d *Database) GetAllLoans() ([]storage.Loan, error) {
	ctx := context.Background()
	loans, err := d.Queries.GetAllLoans(ctx)
	if err != nil {
		return nil, fmtError(err)
	}

	storeLoans := make([]storage.Loan, len(loans))
	for i, l := range loans {
		date := time.Unix(l.Date, 0)
		sLoan := storage.Loan{
			ID:     l.ID,
			Name:   l.Name,
			Date:   date,
			BookID: l.BookID,
		}
		storeLoans[i] = sLoan
	}
	return storeLoans, nil
}

// createLoanWithRet used for testing and is the main logic for CreateLoan().
func (d *Database) createLoanWithRet(loan *storage.Loan) (*database.LoanedBook, error) {
	ctx := context.Background()

	params := database.CreateLoanParams{
		Name:   loan.Name,
		Date:   loan.Date.Unix(),
		BookID: loan.BookID,
	}

	l, err := d.Queries.CreateLoan(ctx, params)
	return &l, err
}


// CreateLona add new loan to database
func (d *Database) CreateLoan(loan *storage.Loan) error {
	_, err := d.createLoanWithRet(loan)
	return err
}

// updateLoanWithRet use for testing and is the main logic for UpdateLoan().
func (d *Database) updateLoanWithRet(loan *storage.Loan) (*database.LoanedBook, error) { 
	ctx := context.Background()

	params := database.UpdateLoanParams{
		ID:    loan.ID,
		Name:  loan.Name,
		Date:  loan.Date.Unix(),
	}

	l, err := d.Queries.UpdateLoan(ctx, params)
	return &l, err
}

// UpdateLoan update given loan.
func (d *Database) UpdateLoan(loan *storage.Loan) error {
	_, err := d.updateLoanWithRet(loan)
	return err
}

// DeleteLoan remove loan from database
func (d *Database) DeleteLoan(loan *storage.Loan) error {
	ctx := context.Background()
	return d.Queries.DeleteLoan(ctx, int64(loan.ID))
}
