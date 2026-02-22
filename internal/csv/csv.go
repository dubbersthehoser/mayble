package csv

import (
	"io"
	"fmt"
	"time"
	"strconv"
	"encoding/csv"
	"errors"

	repo "github.com/dubbersthehoser/mayble/internal/repository"
	"github.com/dubbersthehoser/mayble/internal/fields"
)


var (
	ErrInvalidFormat error = errors.New("invalid format")
)


func entryToFields(book *repo.BookEntry) []string {
	f := make([]string, fields.Length)

	f[fields.Title] = book.Title
	f[fields.Author] = book.Author
	f[fields.Genre] = book.Genre

	if book.Variant & repo.Read != 0 {
		f[fields.Read] = book.Read.Format(time.DateOnly)
		rating := strconv.Itoa(book.Rating)
		f[fields.Rating] = rating
	}
	
	if book.Variant & repo.Loaned != 0 {
		f[fields.Loaned] = book.Loaned.Format(time.DateOnly)
		f[fields.Borrower] = book.Borrower
	}
	return f

}

func fieldsToEntry(f []string) (*repo.BookEntry, error) {
	
	if len(f) != fields.Length {
		return nil, fmt.Errorf("invalid length of slice: %d != %d", fields.Length, len(f))
	}

	book := repo.NewBookEntry()

	book.Title = f[fields.Title]
	book.Author = f[fields.Author]
	book.Genre = f[fields.Genre]

	loaned := f[fields.Loaned]
	borrower := f[fields.Borrower]

	if loaned != "" && borrower != "" {
		book.Variant |= repo.Loaned
		book.Borrower = borrower
		date, err := time.Parse(time.DateOnly, loaned)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid loaned date '%s'", ErrInvalidFormat, loaned)
		}
		book.Loaned = date
	}

	read := f[fields.Read]
	rating := f[fields.Rating]

	if read != "" && rating != "" {
		book.Variant |= repo.Read

		var err error
		book.Rating, err = strconv.Atoi(rating)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid rating '%s'", ErrInvalidFormat, rating)
		}
		book.Read, err = time.Parse(time.DateOnly, read)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid read date '%s'", ErrInvalidFormat, read)
		}
	}

	return book, nil
}

func Import(r io.Reader) ([]repo.BookEntry, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = 7

	entries, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	books := make([]repo.BookEntry, len(entries))

	for i := range entries {
		book, err := fieldsToEntry(entries[i])
		if err != nil {
			return nil, err
		}
		books[i] = *book
	}
	return books, nil
}

func Export(w io.Writer, entries []repo.BookEntry) error {
	
	writer := csv.NewWriter(w)

	fields := make([][]string, len(entries)) 

	for i := range entries {
		fields[i] = entryToFields(&entries[i])
	}
	return writer.WriteAll(fields)
} 











