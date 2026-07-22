package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"database/sql"
	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/sqlite/database"
)

// CreateBook insert b into database.
func (db *Database) CreateBook(b *models.BookEntry) (int64, error) {

	params := database.CreateBookParams{
		Title:  b.Title,
		Author: b.Author,
		Genre:  b.Genre,
	}

	row, err := db.Queries.CreateBook(context.Background(), params)
	if err != nil {
		return -1, fmt.Errorf("database: %w", err)
	}
	bookID := row.ID

	if b.IsLoaned {
		date := b.LoanedAt.Format(time.DateOnly)
		params := database.CreateLoanParams{
			BookID: bookID,
			Name:   b.Borrower,
			Date:   date,
		}
		_, err := db.Queries.CreateLoan(context.Background(), params)
		if err != nil {
			return -1, fmt.Errorf("database: %w", err)
		}
	}

	if b.IsCompleted {
		date := b.CompletedAt.Format(time.DateOnly)
		params := database.CreateReadParams{
			BookID:        bookID,
			Rating:        int64(b.Rating),
			DateCompleted: date,
		}
		_, err := db.Queries.CreateRead(context.Background(), params)
		if err != nil {
			return -1, fmt.Errorf("database: %w", err)
		}
	}
	return bookID, nil
}

// DeleteBook remove book from database errors a database related error.
func (db *Database) DeleteBook(id int64) error {

	err := db.Queries.DeleteBook(context.Background(), id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("database: %w", err)
		}
	}
	err = db.Queries.DeleteLoan(context.Background(), id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("database: %w", err)
		}
	}
	err = db.Queries.DeleteRead(context.Background(), id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("database: %w", err)
		}
	}
	return nil
}

// UpdateBook update book from database.
func (db *Database) UpdateBook(b *models.BookEntry) error {

	var (
		hasLoaned bool = true
		hasRead   bool = true
	)

	_, err := db.Queries.GetLoanByBookID(context.Background(), b.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("database: %w", err)
		}
		hasLoaned = false
	}
	_, err = db.Queries.GetReadByBookID(context.Background(), b.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("database: %w", err)
		}
		hasRead = false
	}

	if b.IsLoaned {

		loanUpdateParams := database.UpdateLoanParams{
			BookID: b.ID,
			Name:   b.Borrower,
			Date:   b.LoanedAt.Format(time.DateOnly),
		}

		loanCreateParams := database.CreateLoanParams{
			BookID: b.ID,
			Name:   b.Borrower,
			Date:   b.LoanedAt.Format(time.DateOnly),
		}

		if hasLoaned {
			_, err := db.Queries.UpdateLoan(context.Background(), loanUpdateParams)
			if err != nil {
				return fmt.Errorf("database: %w", err)
			}
		} else {
			_, err := db.Queries.CreateLoan(context.Background(), loanCreateParams)
			if err != nil {
				return fmt.Errorf("database: %w", err)
			}
		}
	}
	if hasLoaned && !b.IsLoaned {
		err := db.Queries.DeleteLoan(context.Background(), b.ID)
		if err != nil {
			return fmt.Errorf("database: %w", err)
		}
	}

	if b.IsCompleted {

		readUpdateParams := database.UpdateReadParams{
			BookID:        b.ID,
			Rating:        int64(b.Rating),
			DateCompleted: b.CompletedAt.Format(time.DateOnly),
		}

		readCreateParams := database.CreateReadParams{
			BookID:        b.ID,
			Rating:        int64(b.Rating),
			DateCompleted: b.CompletedAt.Format(time.DateOnly),
		}

		if hasRead {
			_, err := db.Queries.UpdateRead(context.Background(), readUpdateParams)
			if err != nil {
				return fmt.Errorf("database: %w", err)
			}
		} else {
			_, err := db.Queries.CreateRead(context.Background(), readCreateParams)
			if err != nil {
				return fmt.Errorf("database: %w", err)
			}
		}
	} else {
		if hasRead {
			err := db.Queries.DeleteRead(context.Background(), b.ID)
			if err != nil {
				return fmt.Errorf("database: %w", err)
			}
		}
	}

	updateBookParams := database.UpdateBookParams{
		ID:     b.ID,
		Title:  b.Title,
		Author: b.Author,
		Genre:  b.Genre,
	}

	_, err = db.Queries.UpdateBook(context.Background(), updateBookParams)
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}
	return nil
}

// GetAllBooks returns all books from database when v is zero, otherwise filters for variant.
func (db *Database) GetAllBooks() ([]models.BookEntry, error) {

	books, err := db.Queries.GetAllBooks(context.Background())
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}
	var entries []models.BookEntry
	for _, book := range books {
		hasLoaned := true
		hasRead := true
		loan, err := db.Queries.GetLoanByBookID(context.Background(), book.ID)
		if err != nil { // only return error when err is not ErrNoRows.
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("database: %w", err)
			}
			hasLoaned = false
		}
		read, err := db.Queries.GetReadByBookID(context.Background(), book.ID)
		if err != nil { // only return error when err is not ErrNoRows.
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("database: %w", err)
			}
			hasRead = false
		}

		builder := models.NewBookEntryBuilder()

		builder.SetID(book.ID).
			SetTitle(book.Title).
			SetAuthor(book.Author).
			SetGenre(book.Genre)

		if hasLoaned {
			builder.SetLoaned(loan.Date).
				SetBorrower(loan.Name)
		}

		if hasRead {
			builder.SetCompleted(read.DateCompleted).
				SetRating(int(read.Rating))
		}
		book, err := builder.Build()
		if err != nil {
			return nil, fmt.Errorf("database: %w", err)
		}
		entries = append(entries, *book)
	}
	return entries, nil
}

// GetBookByID returns book entry by id.
func (db *Database) GetBookByID(id int64) (models.BookEntry, error) {

	bookRow, err := db.Queries.GetBookByID(context.Background(), id)
	if err != nil {
		return models.BookEntry{}, fmt.Errorf("database: %w", err)
	}

	hasLoan := true
	loanRow, err := db.Queries.GetLoanByBookID(context.Background(), id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return models.BookEntry{}, fmt.Errorf("database: %w", err)
		}
		hasLoan = false
	}

	hasRead := true
	readRow, err := db.Queries.GetReadByBookID(context.Background(), id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return models.BookEntry{}, fmt.Errorf("database: %w", err)
		}
		hasRead = false
	}

	builder := models.NewBookEntryBuilder()

	builder.SetID(id).
		SetTitle(bookRow.Title).
		SetAuthor(bookRow.Author).
		SetGenre(bookRow.Genre)

	if hasLoan {
		builder.SetLoaned(loanRow.Date).
			SetBorrower(loanRow.Name)
	}

	if hasRead {
		builder.SetCompleted(readRow.DateCompleted).
			SetRating(int(readRow.Rating))
	}

	book, err := builder.Build()
	if err != nil {
		return models.BookEntry{}, fmt.Errorf("database: %w", err)
	}

	return *book, nil
}

// GetUniqueGenres return unique set of genres from database.
func (db *Database) GetUniqueGenres() ([]string, error) {
	genres, err := db.Queries.GetUniqueGenres(context.Background())
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}
	return genres, nil
}
