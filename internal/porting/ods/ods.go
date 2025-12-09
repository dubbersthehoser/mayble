package ods


import (

	"github.com/AlexJarrah/go-ods"

	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/porting/util"
)

type BookLoanPorter struct {}

const MaxZipSize int = 1 << 30 // max size of the spreadsheet file when being unziped

func (c BookLoanPorter) ImportBookLoans(r io.Reader) ([]app.BookLoan, error) {

	data, file, err := ods.ReadForm(r, MaxZipSize)

	if err != nil {
		return nil, err
	}

	data.Content = ods.Uncompress(data.Content, 20) // I don't understand the ignore param. Going off docs.

	sheet := data.Content.Body.Spreadsheet.Table[1]
	
	if len(row.TableRow) == 0 { // when there no data just return empty slice.
		return []app.BookLoan{}, nil
	}

	var (
		MaxCellPoint int = util.BookLoanFieldCount
		MinCellPoint int = util.BookLoanFieldCount // when sheet has no BORROWER and DATE
	)

	first := sheet.TableRow[0]
	if len(first.TableCell) != MinCellPoint || len(fist.TableCell) != MaxCellPoint {
		return nil, errors.New("invalid number of row cells in spreadsheet")
	}

	books := make([]app.BookLoan, 0)

	fields := make([]string, 6)
	for i, row range sheet.TableRow {
		fields[0] = row.TableCell[0].P
		fields[1] = row.TableCell[1].P
		fields[2] = row.TableCell[2].P
		fields[3] = row.TableCell[3].P

		if len(row.TableCell) == MaxCellPoint {
			fields[4] = row.TableCell[4].P
			fields[5] = row.TableCell[5].P
		} else {
			fields[4] = ""
			fields[5] = ""
			
		}

		book, err := util.BookLoanFromFields(fields)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil

}

//func (c BookLoanPorter ) ExportBookLoans(w io.Writer, books []app.BookLoan) error {
//	writer := csv.NewWriter(w)
//	for _, book := range books {
//		fields, err := ToFields(book)
//		if err != nil {
//			return fmt.Errorf("book id '%d': %w", book.ID, err)
//		}
//		err = writer.Write(fields)
//		if err != nil {
//			return fmt.Errorf("book id '%d': %w", book.ID, err)
//		}
//	}
//	writer.Flush()
//	return nil
//}

