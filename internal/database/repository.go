package database

import (
	"context"
	"time"
	"errors"
	
	"database/sql"
	repo "github.com/dubbersthehoser/mayble/internal/repository"
	"github.com/dubbersthehoser/mayble/internal/status"
	"github.com/dubbersthehoser/mayble/internal/sqlite/database"
)


func (db *Database) CreateBook(b *repo.BookEntry) error {
	const op status.Op = "database.CreateBook"

	params := database.CreateBookParams{
		Title: b.Title,
		Author: b.Author,
		Genre: b.Genre,
	}
	row, err := db.Queries.CreateBook(context.Background(), params)
	if err != nil {
		return status.E(op, status.LevelWarn, err)
	}
	bookID := row.ID

	if b.Variant & repo.Loaned != 0 {
		date := b.Loaned.Format(time.DateOnly)
		params := database.CreateLoanParams{
			BookID: bookID,
			Name: b.Borrower,
			Date: date,
		}
		_, err := db.Queries.CreateLoan(context.Background(), params)
		if err != nil {
			return status.E(op, status.LevelWarn, err)
		}
	}

	
	if b.Variant & repo.Read != 0 {
		date := b.Loaned.Format(time.DateOnly)
		params := database.CreateReadParams{
			BookID: bookID,
			Rating: int64(b.Rating),
			DateCompleted: date,
		}
		_, err := db.Queries.CreateRead(context.Background(), params)
		if err != nil {
			return status.E(op, status.LevelWarn, err)
		}
	}
	return nil
}

func (db *Database) DeleteBook(b *repo.BookEntry) error {
	const op status.Op = "database.DeleteBook"
	id := b.ID
	err := db.Queries.DeleteBook(context.Background(), id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return status.E(op, status.LevelError, err)
		}
	}
	err = db.Queries.DeleteLoan(context.Background(), id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return status.E(op, status.LevelError, err)
		}
	}
	err = db.Queries.DeleteRead(context.Background(), id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return status.E(op, status.LevelError, err)
		}
	}
	return nil
}

func (db *Database) UpdateBook(b *repo.BookEntry) error {
	return nil
}

func (db *Database) GetAllBooks(v repo.Variant) ([]repo.BookEntry, error) {
	const op status.Op = "database.GetAllBooks"

	books, err := db.Queries.GetAllBooks(context.Background())
	entries := make([]repo.BookEntry, len(books))
	if err != nil {
		return nil, status.E(op, status.LevelWarn, err)
	}
	for i, book := range books {
		hasLoaned := true
		hasRead := true
		loan, err := db.Queries.GetLoanByBookID(context.Background(), book.ID)
		if err != nil { // return error when err is not ErrNoRows. 
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, status.E(op, status.LevelWarn, err)
			}
			hasLoaned = false
		}
		read, err := db.Queries.GetReadByBookID(context.Background(), book.ID)
		if err != nil { // return error when err is not ErrNoRows. 
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, status.E(op, status.LevelWarn, err)
			}
			hasRead = false
		}
		
		entries[i].Variant |= repo.Book
		entries[i].Title = book.Title
		entries[i].Author = book.Author
		entries[i].Genre = book.Genre

		if hasLoaned {
			entries[i].Variant |= repo.Loaned
			date, err := time.Parse(time.DateOnly, loan.Date)
			if err != nil {
				return nil, status.E(op, status.LevelError, err)
			}
			entries[i].Loaned = date
			entries[i].Borrower = loan.Name
		}

		if hasRead {
			entries[i].Variant |= repo.Read
			date, err := time.Parse(time.DateOnly, read.DateCompleted)
			if err != nil {
				return nil, status.E(op, status.LevelError, err)
			}
			entries[i].Read = date
			entries[i].Rating = int(read.Rating)
		}
	}
	return entries, nil
}

func (db *Database) GetBookByID(id int64) (repo.BookEntry, error) {
	return repo.BookEntry{}, nil
}

func (a *Database) GetUniqueGenres() ([]string, error) {
	return []string{
		"Cat",
		"Dog",
		"Bird",
	}, nil
}




