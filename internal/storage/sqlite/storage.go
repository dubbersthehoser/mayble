package sqlite

import (
	"fmt"
	"context"
	"time"
	"errors"
	"database/sql"

	"github.com/dubbersthehoser/mayble/internal/storage"
	"github.com/dubbersthehoser/mayble/internal/sqlite"
	"github.com/dubbersthehoser/mayble/internal/sqlite/database"
)

type Storage struct {
	sqlite.Database
}


func NewStorage(path string) (*Storage, error) {
	db := sqlite.NewDatabase()
	err := db.Open(path)
	if err != nil {
		return nil, err
	}
	err = db.MigrateUp()
	if err != nil {
		return nil, err
	}
	return &Storage{Database: *db}, nil

}



/*
        BookStore
*/

func (d *Storage) GetBookByID(id int64) (*storage.Book, error) {
	book, err := d.Queries.GetBookByID(context.Background(), id)
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

	if err != nil {
		return nil, err
	}

	ret := &storage.Book{
		ID: book.ID,
		Title: book.Title,
		Author: book.Author,
		Genre: book.Genre,
		Ratting: int(book.Ratting),
	}
	return ret, nil
}



func (d *Storage) GetBooks() ([]storage.Book, error) {	

	books, err := d.Queries.GetAllBooks(context.Background())
	if err != nil {
		return nil, err
	}

	ret := make([]storage.Book, len(books))
	for i := range books {
		book := storage.Book{
			ID: books[i].ID,
			Title: books[i].Title,
			Author: books[i].Author,
			Genre: books[i].Genre,
			Ratting: int(books[i].Ratting),
		}
		ret[i] = book
		
	}
	return ret, nil
}

func (d *Storage) CreateBook(id int64, title, author, genre string, ratting int) (int64, error) { 
	
	if storage.IsZeroID(id) {
		books, err := d.GetBooks()
		if err != nil {
			return id, err
		}
		id = int64(len(books) + 1)
	}

	params := database.CreateBookParams{
		ID:      id,
		Title:   title,
		Author:  author,
		Genre:   genre,
		Ratting: int64(ratting),
	}

	_, err := d.Queries.CreateBook(context.Background(), params)
	if err != nil {
		return id, err
	}
	return id, nil
}



func (d *Storage) UpdateBook(book *storage.Book) error {
	params := database.UpdateBookParams{
		ID: book.ID,
		Title: book.Title,
		Author: book.Author,
		Genre: book.Genre,
		Ratting: int64(book.Ratting),
	}
	_, err := d.Queries.UpdateBook(context.Background(), params)
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
	return nil
}

func (d *Storage) DeleteBook(book *storage.Book) error {
	
	if storage.IsZeroID(book.ID) {
		return fmt.Errorf("%w: given zero id", storage.ErrInvalidValue)
	}

	err := d.Queries.DeleteBook(context.Background(), book.ID)
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

	return nil
}


/*
        LoanStore
*/


func (d *Storage) GetLoan(id int64) (*storage.Loan, error) {
	loan, err := d.Queries.GetLoanByBookID(context.Background(), id)
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

	date, err := time.Parse(time.DateOnly, loan.Date)

	if err != nil {
		return nil, err
	}
	ret := &storage.Loan{
		ID: loan.BookID,
		Borrower: loan.Name,
		Date: date,
	}
	return ret, nil
}
