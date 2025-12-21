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

func (d *Storage) Close() error {
	return d.Database.Close()
}

/*
        BookStore
*/

func (d *Storage) GetBookByID(id int64) (*storage.Book, error) {
	
	if storage.IsZeroID(id) {
		return nil, storage.ErrInvalidValue
	}
	
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
	
	var createdID bool
	if storage.IsZeroID(id) {
		books, err := d.GetBooks()
		if err != nil {
			return id, err
		}
		id = int64(len(books) + 1)
		createdID = true
	}

	_, err := d.Queries.GetBookByID(context.Background(), id)
	if err == nil && !createdID {
		return id, storage.ErrEntryExists
	} 
	for err == nil && createdID {
		id += 1
		_, err = d.Queries.GetBookByID(context.Background(), id)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return id, err
	}


	params := database.CreateBookParams{
		ID:      id,
		Title:   title,
		Author:  author,
		Genre:   genre,
		Ratting: int64(ratting),
	}

	_, err = d.Queries.CreateBook(context.Background(), params)
	if err != nil { 
		return id, err
	}
	return id, nil
}



func (d *Storage) UpdateBook(book *storage.Book) error {
	if storage.IsZeroID(book.ID) {
		return storage.ErrInvalidValue
	}

	params := database.UpdateBookParams{
		ID: book.ID,
		Title: book.Title,
		Author: book.Author,
		Genre: book.Genre,
		Ratting: int64(book.Ratting),
	}

	_, err := d.Queries.GetBookByID(context.Background(), book.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrEntryNotFound
	}
	if err != nil {
		return err
	}
	_, err = d.Queries.UpdateBook(context.Background(), params)
	return err
}

func (d *Storage) DeleteBook(book *storage.Book) error {
	
	if storage.IsZeroID(book.ID) {
		return fmt.Errorf("%w: given zero id", storage.ErrInvalidValue)
	}
	_, err := d.Queries.GetBookByID(context.Background(), book.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrEntryNotFound
	}
	if err != nil {
		return err
	}

	err = d.Queries.DeleteBook(context.Background(), book.ID)
	return err
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

func (d *Storage) CreateLoan(bookID int64, borrower string, date time.Time) error {
	
	if storage.IsZeroID(bookID) {
		return storage.ErrInvalidValue
	}

	_, err := d.Queries.GetBookByID(context.Background(), bookID)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrEntryNotFound
	}
	if err != nil {
		return err
	}

	_, err = d.Queries.GetLoanByBookID(context.Background(), bookID)
	if err == nil {
		return storage.ErrEntryExists
	} 
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	params := database.CreateLoanParams{
		BookID: bookID,
		Name: borrower,
		Date: date.Format(time.DateOnly),
	}

	_, err = d.Queries.CreateLoan(context.Background(), params)
	return err
}

func (d *Storage) UpdateLoan(loan *storage.Loan) error {

	if storage.IsZeroID(loan.ID) {
		return storage.ErrInvalidValue
	}

	_, err := d.Queries.GetLoanByBookID(context.Background(), loan.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrEntryNotFound
	}
	if err != nil {
		return err
	}

	params := database.UpdateLoanParams{
		BookID: loan.ID,
		Name: loan.Borrower,
		Date: loan.Date.Format(time.DateOnly),
	}
	_, err = d.Queries.UpdateLoan(context.Background(), params)
	return err
}

func (d *Storage) DeleteLoan(loan *storage.Loan) error {
	if storage.IsZeroID(loan.ID) {
		return storage.ErrInvalidValue
	}

	_, err := d.Queries.GetLoanByBookID(context.Background(), loan.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrEntryNotFound
	}
	if err != nil {
		return err
	}
	
	return d.Queries.DeleteLoan(context.Background(), loan.ID)
}



