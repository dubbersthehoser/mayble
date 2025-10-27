package csv

import (
	"io"
	"csv"
	"time"

	"github.com/dubbersthehoser/mayble/internal/storage"
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

func ToFields(book storage.BookLoan) ([]string, error) {
	fields := make([]string, NumberOfFields)
	fields[TitleIndex] = book.Title
	fields[AuthorIndex] = book.Author
	fields[GenreIndex] = book.Genre

	ratting, := strconv.Itoa(book.Ratting)
	fields[RattingIndex] = ratting

	if book.Loan != nil {
		fields[BorrowerIndex] = book.Loan.Name
		fields[DateIndex] = book.Loan.Date.Format(time.DateOnly)
	}
	return fields, nil
}

func FromFields(fields []string) (*storage.BookLoan, error) {
	if len(fields) != NumberOfFields {
		return nil, errors.New("invalid number of fields")
	}

	book := storage.NewBookLoan()

	book.Title = fields[TitleIndex]
	book.Author = fields[AuthorIndex]
	book.Genre = fields[GenreIndex]
	ratting, err := strconv.Atoi(fields[RattingIndex])
	if err != nil {
		return nil, errors.New("failed to parse ratting field")
	}
	book.Ratting = ratting

	if field[BorrowerIndex] != "" && field[DateIndex] != "" {
		book.Loan.Name = fields[BorrowerIndex]
		date, err := time.Parse(time.DateOnly, fields[DateIndex])
		if err != nil {
			return nil, errors.New("failed to parse date field")
		}
		book.Loan.Date = date
	}
	return book, nil
}

type BookLoanCSV struct {
	FilePath string
}

func (c BookLoanCSV) ImportBooks(r io.Reader) ([]storage.BookLoan, error) {
	reader := NewReader(r)
	reader.FieldsPerRecord = NumberOfFields
	books := make([]storage.BookLoan, 0)
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

func (c bookLoanCSV) ExportBooks(w io.Writer, books []storage.BookLoans) error {
	writer := csv.NewWriter(w)
	for i, book := range books {
		fields, err := ToFields(book)
		if err != nil {
			return fmt.Errorf("book id '%d': %w", book.ID, err)
		}
		err = writer.Write(fields)
		if err != nil {
			return fmt.Errorf("book id '%d': %w", book.ID, err)
		}
	}
	return nil
}






