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

	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

var (
	ErrInvalidFormat error = errors.New("invalid format")
)

func Import(r io.Reader) ([]repo.BookEntry, error) {
	
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = 7

	entries, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	books := make([]repo.BookEntry, len(entries))

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
		books[i] = *book
	}
	return books, nil
}

func Export(w io.Writer, entries []repo.BookEntry) error {

	writer := csv.NewWriter(w)

	fields := make([][]string, len(entries)+1)

	fields[0] = schemaHeaders()

	for i := range entries {
		fields[i+1] = entryToFields(&entries[i])
	}
	return writer.WriteAll(fields)
}

func fieldLength() int {
	return len(repo.BookEntryFields())
}

func entryToFields(book *repo.BookEntry) []string {
	f := make([]string, fieldLength())

	f[repo.IdxTitle] = book.Title
	f[repo.IdxAuthor] = book.Author
	f[repo.IdxGenre] = book.Genre

	if book.Variant&repo.Read != 0 {
		f[repo.IdxRead] = book.Read.Format(time.DateOnly)
		rating := strconv.Itoa(book.Rating)
		f[repo.IdxRating] = rating
	}

	if book.Variant&repo.Loaned != 0 {
		f[repo.IdxLoaned] = book.Loaned.Format(time.DateOnly)
		f[repo.IdxBorrower] = book.Borrower
	}
	return f

}

func fieldsToEntry(f []string, m []int) (*repo.BookEntry, error) {

	if len(f) != fieldLength() {
		return nil, fmt.Errorf("invalid length of slice: %d != %d", fieldLength(), len(f))
	}

	builder := repo.NewBookEntryBuilder()

	builder.SetTitle(f[m[repo.IdxTitle]]).
		SetAuthor(f[m[repo.IdxAuthor]]).
		SetGenre(f[m[repo.IdxGenre]]).
		SetBorrower(f[m[repo.IdxBorrower]]).
		SetLoaned(f[m[repo.IdxLoaned]]).
		SetCompleted(f[m[repo.IdxRead]]).
		SetRating(f[m[repo.IdxRating]]).
		SetID(0)

	return builder.Build()
}

func schemaHeaders() []string {
	headers := repo.BookEntryFields()
	for i := range headers {
		headers[i] = strings.ToUpper(headers[i])
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
