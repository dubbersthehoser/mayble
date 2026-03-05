package csv

import (
	"io"
	"fmt"
	"time"
	"strconv"
	"encoding/csv"
	"errors"

	repo "github.com/dubbersthehoser/mayble/internal/repository"
)


var (
	ErrInvalidFormat error = errors.New("invalid format")
)

func fieldLength() int {
	return len(repo.BookEntryFields())
}

func entryToFields(book *repo.BookEntry) []string {
	f := make([]string, fieldLength())

	f[repo.IdxTitle] = book.Title
	f[repo.IdxAuthor] = book.Author
	f[repo.IdxGenre] = book.Genre

	if book.Variant & repo.Read != 0 {
		f[repo.IdxRead] = book.Read.Format(time.DateOnly)
		rating := strconv.Itoa(book.Rating)
		f[repo.IdxRating] = rating
	}
	
	if book.Variant & repo.Loaned != 0 {
		f[repo.IdxLoaned] = book.Loaned.Format(time.DateOnly)
		f[repo.IdxBorrower] = book.Borrower
	}
	return f

}

func fieldsToEntry(f []string) (*repo.BookEntry, error) {
	
	if len(f) != fieldLength() {
		return nil, fmt.Errorf("invalid length of slice: %d != %d", fieldLength(), len(f))
	}

	book := &repo.BookEntry{}

	book.Title = f[repo.IdxTitle]
	book.Author = f[repo.IdxAuthor]
	book.Genre = f[repo.IdxGenre]

	loaned := f[repo.IdxLoaned]
	borrower := f[repo.IdxBorrower]

	if loaned != "" && borrower != "" {
		book.Variant |= repo.Loaned
		book.Borrower = borrower
		date, err := time.Parse(time.DateOnly, loaned)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid loaned date '%s'", ErrInvalidFormat, loaned)
		}
		book.Loaned = date
	}

	read := f[repo.IdxRead]
	rating := f[repo.IdxRating]

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
