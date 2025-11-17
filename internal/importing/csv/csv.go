package csv

import (
	"io"
	"encoding/csv"
	"time"
	"strconv"
	"errors"
	"fmt"

	"github.com/dubbersthehoser/mayble/internal/data"
)


// fields: TITLE,AUTHOR,GENRE,RATTING,BORROWER,DATE

const NumberOfFields int = 6

const (
	TitleIndex int = 0
	AuthorIndex    = 1
	GenreIndex     = 2
	RattingIndex   = 3
	BorrowerIndex  = 4
	DateIndex      = 5

)

func ToFields(book data.BookLoan) ([]string, error) {
	fields := make([]string, NumberOfFields)
	fields[TitleIndex] = book.Title
	fields[AuthorIndex] = book.Author
	fields[GenreIndex] = book.Genre

	ratting := strconv.Itoa(book.Ratting)
	fields[RattingIndex] = ratting

	if book.Loan != nil {
		fields[BorrowerIndex] = book.Loan.Borrower
		fields[DateIndex] = book.Loan.Date.Format(time.DateOnly)
	}
	return fields, nil
}

func FromFields(fields []string) (*data.BookLoan, error) {
	if len(fields) != NumberOfFields {
		return nil, errors.New("invalid number of fields")
	}

	book := data.NewBookLoan()

	book.Title = fields[TitleIndex]
	book.Author = fields[AuthorIndex]
	book.Genre = fields[GenreIndex]
	ratting, err := strconv.Atoi(fields[RattingIndex])
	if err != nil {
		return nil, errors.New("failed to parse ratting field")
	}

	book.Ratting = ratting
	if fields[BorrowerIndex] != "" && fields[DateIndex] != "" {
		book.Loan.Borrower = fields[BorrowerIndex]
		date, err := time.Parse(time.DateOnly, fields[DateIndex])
		if err != nil {
			return nil, errors.New("failed to parse date field")
		}
		book.Loan.Date = date
	} else {
		book.Loan = nil
	}
	return book, nil
}

type BookLoanCSV struct {}

func (c BookLoanCSV) ImportBooks(r io.Reader) ([]data.BookLoan, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = NumberOfFields
	books := make([]data.BookLoan, 0)
	recordCount := 0
	for {
		fields, err := reader.Read()
		recordCount += 1
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("record %d: %w", recordCount, err )
		}

		book, err := FromFields(fields)
		if err != nil {
			return nil, fmt.Errorf("record %d: %w", recordCount, err)
		}
		books = append(books, *book)
	}

	return books, nil

}

func (c BookLoanCSV) ExportBooks(w io.Writer, books []data.BookLoan) error {
	writer := csv.NewWriter(w)
	for _, book := range books {
		fields, err := ToFields(book)
		if err != nil {
			return fmt.Errorf("book id '%d': %w", book.ID, err)
		}
		err = writer.Write(fields)
		if err != nil {
			return fmt.Errorf("book id '%d': %w", book.ID, err)
		}
	}
	writer.Flush()
	return nil
}






