package csv

import (
	"encoding/csv"
	"slices"
	"strings"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	models "github.com/dubbersthehoser/mayble/internal/models"
)

var (
	ErrInvalidFormat error = errors.New("invalid format")
)

const (
	idxTitle int = iota
	idxAuthor
	idxGenre
	idxCompletedAt
	idxRating
	idxLoanedAt
	idxBorrower
)

func Import(r io.Reader) ([]models.BookEntry, error) {
	
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = 7

	entries, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	books := make([]models.BookEntry, 0)

	if len(entries) == 0 {
		return books, nil
	}

	mapping, ok := mapSchema(entries[0])
	if !ok {
		mapping = make([]int, fieldLength())
		for i := range mapping {
			mapping[i] = i
		}
	} else {
		entries = entries[1:]
	}

	for i := range entries {
		book, err := fieldsToEntry(entries[i], mapping)
		if err != nil {
			return nil, err
		}
		books = append(books, *book)
	}
	return books, nil
}

func Export(w io.Writer, entries []models.BookEntry) error {

	writer := csv.NewWriter(w)

	fields := make([][]string, len(entries)+1)

	fields[0] = schemaHeaders()

	for i := range entries {
		fields[i+1] = entryToFields(&entries[i])
	}
	return writer.WriteAll(fields)
}

func fieldLength() int {
	return len(models.BookEntryFields())-1 // neg one to avoid ID Field
}

func entryToFields(book *models.BookEntry) []string {
	f := make([]string, fieldLength())

	f[idxTitle] = book.Title
	f[idxAuthor] = book.Author
	f[idxGenre] = book.Genre

	if book.IsCompleted {
		f[idxCompletedAt] = book.CompletedAt.Format(time.DateOnly)
		rating := strconv.Itoa(book.Rating)
		f[idxRating] = rating
	}

	if book.IsLoaned {
		f[idxLoanedAt] = book.LoanedAt.Format(time.DateOnly)
		f[idxBorrower] = book.Borrower
	}
	return f
}

func fieldsToEntry(f []string, m []int) (*models.BookEntry, error) {

	if len(f) != fieldLength() {
		return nil, fmt.Errorf("invalid length of slice: %d != %d", fieldLength(), len(f))
	}

	builder := models.NewBookEntryBuilder()

	builder.SetTitle(f[m[models.IdxTitle]]).
		SetAuthor(f[m[models.IdxAuthor]]).
		SetGenre(f[m[models.IdxGenre]]).
		SetBorrower(f[m[models.IdxBorrower]]).
		SetLoaned(f[m[models.IdxLoanedAt]]).
		SetCompleted(f[m[models.IdxCompletedAt]]).
		SetRating(f[m[models.IdxRating]]).
		SetID(0)

	book, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return book, nil
}

func schemaHeaders() []string {
	headers := models.BookEntryFields()
	for i := range headers {
		if i != 0 {
			headers[i] = strings.ToUpper(headers[i])
		}
	}
	return headers
}

func mapSchema(f []string) ([]int, bool) {
	mapping := make([]int, fieldLength())
	headers := schemaHeaders()
	if len(headers) != len(f) {
		return nil, false
	}
	for i, field := range f {
		idx := slices.Index(headers, field)
		if idx == -1 {
			return nil, false
		}
		mapping[i] = idx
	}
	return mapping, true
}
