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


type bookBuilder struct {
	id                   int64
	title, author, genre string
	loanedDate, readDate string
	borrowerName         string
	rating               int64
}

func (b *bookBuilder) Build() (*repo.BookEntry, error) {
	
	book := &repo.BookEntry{
		ID: b.id,
		Title: b.title,
		Author: b.author,
		Genre: b.genre,
	}

	if b.loanedDate != "" && b.borrowerName != "" {
		book.Variant |= repo.Loaned
		date, err := time.Parse(time.DateOnly, b.loanedDate)
		if err != nil {
			return nil, err
		}
		book.Loaned = date
		book.Borrower = b.borrowerName
	}

	if b.readDate != "" {
		book.Variant |= repo.Read
		date, err := time.Parse(time.DateOnly, b.readDate)
		if err != nil {
			return nil, err
		}

		book.Read = date
		book.Rating = int(b.rating)
	}

	return book, nil

}

func (b *bookBuilder) SetID(id int64) *bookBuilder {
	b.id = id
	return b
}

func (b *bookBuilder) SetTitle(t string) *bookBuilder {
	b.title = t
	return b
}

func (b *bookBuilder) SetAuthor(a string) *bookBuilder {
	b.author = a
	return b
}

func (b *bookBuilder) SetGenre(g string) *bookBuilder {
	b.genre = g
	return b
}

func (b *bookBuilder) SetLoanedDate(d string) *bookBuilder {
	b.loanedDate = d
	return b
}

func (b *bookBuilder) SetReadDate(d string) *bookBuilder {
	b.readDate = d
	return b
}

func (b *bookBuilder) SetBorrower(n string) *bookBuilder {
	b.borrowerName = n
	return b
}

func (b *bookBuilder) SetRating(r int64) *bookBuilder {
	b.rating = r
	return b
}



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
	const op status.Op = "database.UpdateBook"

	var (
		hasLoaned bool = true
		hasRead   bool = true

	)

	_, err := db.Queries.GetLoanByBookID(context.Background(), b.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return status.E(op, status.LevelError, err)
		}
		hasLoaned = false
	}
	_, err = db.Queries.GetReadByBookID(context.Background(), b.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return status.E(op, status.LevelError, err)
		}
		hasRead = false
	}

	if b.Variant & repo.Loaned != 0 {

		loanUpdateParams := database.UpdateLoanParams{
			BookID: b.ID,
			Name: b.Borrower,
			Date: b.Loaned.Format(time.DateOnly),
		}

		loanCreateParams := database.CreateLoanParams{
			BookID: b.ID,
			Name: b.Borrower,
			Date: b.Loaned.Format(time.DateOnly),
		}

		if hasLoaned {
			_, err := db.Queries.UpdateLoan(context.Background(), loanUpdateParams)
			if err != nil {
				return status.E(op, status.LevelError, err)
			}
		} else {
			_, err := db.Queries.CreateLoan(context.Background(), loanCreateParams)
			if err != nil {
				return status.E(op, status.LevelError, err)
			}
		}
	} else {
		if hasLoaned {
			err := db.Queries.DeleteLoan(context.Background(), b.ID)
			if err != nil {
				return status.E(op, status.LevelError, err)
			}
		}
	}

	if b.Variant & repo.Read != 0 {

		readUpdateParams := database.UpdateReadParams{
			BookID: b.ID,
			Rating: int64(b.Rating),
			DateCompleted: b.Read.Format(time.DateOnly),
		}

		readCreateParams := database.CreateReadParams{
			BookID: b.ID,
			Rating: int64(b.Rating),
			DateCompleted: b.Read.Format(time.DateOnly),
		}

		if hasRead {
			_, err := db.Queries.UpdateRead(context.Background(), readUpdateParams)
			if err != nil {
				return status.E(op, status.LevelError, err)
			}
		} else {
			_, err := db.Queries.CreateRead(context.Background(), readCreateParams)
			if err != nil {
				return status.E(op, status.LevelError, err)
			}
		}
	} else {
		if hasRead {
			err := db.Queries.DeleteRead(context.Background(), b.ID)
			if err != nil {
				return status.E(op, status.LevelError, err)
			}
		}
	}

	updateBookParams := database.UpdateBookParams{
		ID: b.ID,
		Title: b.Title,
		Author: b.Author,
		Genre: b.Genre,
	}

	_, err = db.Queries.UpdateBook(context.Background(), updateBookParams)
	if err != nil {
		return status.E(op, status.LevelError, err)
	}
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

		builder := &bookBuilder{}
		builder.SetID(book.ID).
			SetTitle(book.Title).
			SetAuthor(book.Author).
			SetGenre(book.Genre)

		if hasLoaned {
			builder.SetLoanedDate(loan.Date).
				SetBorrower(loan.Name)
		}

		if hasRead {
			builder.SetReadDate(read.DateCompleted).
				SetRating(read.Rating)
		}
		book, err := builder.Build()
		if err != nil {
			return nil, status.E(op, status.LevelWarn, err)
		}
		entries[i] = *book
	}
	return entries, nil
}


func (db *Database) GetBookByID(id int64) (repo.BookEntry, error) {
	const op status.Op = "database.GetBookByID"
	
	bookRow, err := db.Queries.GetBookByID(context.Background(), id)
	if err != nil {
		return repo.BookEntry{}, status.E(op, status.LevelError, err)
	}


	hasLoan := true
	loanRow, err := db.Queries.GetLoanByBookID(context.Background(), id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return repo.BookEntry{}, status.E(op, status.LevelError, err)
		}
		hasLoan = false
	}
	
	hasRead := true
	readRow, err := db.Queries.GetReadByBookID(context.Background(), id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return repo.BookEntry{}, status.E(op, status.LevelError, err)
		}
		hasRead = false
	}

	builder := &bookBuilder{}

	builder.SetID(id).
		SetTitle(bookRow.Title).
		SetAuthor(bookRow.Author).
		SetGenre(bookRow.Genre)

	if hasLoan {
		builder.SetLoanedDate(loanRow.Date).
			SetBorrower(loanRow.Name)
	}

	if hasRead {
		builder.SetReadDate(readRow.DateCompleted).
			SetRating(readRow.Rating)
	}

	book, err := builder.Build()
	if err != nil {
		return repo.BookEntry{}, status.E(op, status.LevelError, err)
	}

	return *book, nil
}

func (db *Database) GetUniqueGenres() ([]string, error) {
	const op status.Op = "database.GetUniqueGenres"
	genres, err := db.Queries.GetUniqueGenres(context.Background())
	if err != nil {
		return nil, status.E(op, status.LevelError, err)
	}
	return genres, nil
}
