package csv

import (
	"io"
	"encoding/csv"
	"errors"
	"fmt"

	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/porting/util"
)

type BookLoanPorter struct {}

func (c BookLoanPorter ) ImportBookLoans(r io.Reader) ([]app.BookLoan, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = util.BookLoanFieldCount 
	books := make([]app.BookLoan, 0)
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

		book, err := util.BookLoanFromFields(fields)
		if err != nil {
			return nil, fmt.Errorf("record %d: %w", recordCount, err)
		}
		books = append(books, *book)
	}

	return books, nil

}

func (c BookLoanPorter ) ExportBookLoans(w io.Writer, books []app.BookLoan) error {
	writer := csv.NewWriter(w)
	for _, book := range books {
		fields, err := util.BookLoanToFields(book)
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
