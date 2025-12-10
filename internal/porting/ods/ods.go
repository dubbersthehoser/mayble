// Package ods provides Open Document Spreadsheet as a 
// `github.com/dubbersthehoser/mayble/porting` implementation.
//
//  NOTE: ExportBookLoans method is incomplete. Need a template file to add data to.
//        Not having the appropate surounding data / metadata will create a broken file.
//        This package WILL NOT BE USED until fixed. Or just deleted. 
//        The ImportBookLoans is untested but confedented it could work even if bugged.
//        Test before used.
//
package ods



import (
	"io"
	"errors"
	"strconv"
	"bytes"
	"archive/zip"

	ods "github.com/AlexJarrah/go-ods"

	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/porting/util"
)

type BookLoanPorter struct {}

const MaxZipSize int64 = 1 << 30 // max size of the spreadsheet file when being unziped

func (c BookLoanPorter) ImportBookLoans(r io.Reader) ([]app.BookLoan, error) {

	data, _, err := ods.ReadFrom(r, MaxZipSize)

	if err != nil {
		return nil, err
	}

	data.Content = ods.Uncompress(data.Content, 20) // I don't understand the ignore param. Going off docs.

	sheet := data.Content.Body.Spreadsheet.Table[0]
	
	if len(sheet.TableRow) == 0 { // when there no data just return empty slice.
		return []app.BookLoan{}, nil
	}

	var (
		MaxCellPoint int = util.BookLoanFieldCount
		MinCellPoint int = util.BookLoanFieldCount // when sheet has no BORROWER and DATE
	)

	first := sheet.TableRow[0]
	if len(first.TableCell) != MinCellPoint || len(first.TableCell) != MaxCellPoint {
		return nil, errors.New("invalid number of row cells in spreadsheet")
	}

	books := make([]app.BookLoan, 0)

	fields := make([]string, 6)
	for _, row := range sheet.TableRow {
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
		books = append(books, *book)
	}

	return books, nil

}


func (c BookLoanPorter ) ExportBookLoans(w io.Writer, books []app.BookLoan) error {

	sheet := ods.Table{}
	row := ods.TableRow{
		TableCell:  make([]ods.TableCell, 6),
	}
	for _, book := range books {
		fields, err := util.BookLoanToFields(book)
		if err != nil {
			return err
		} 
		for i := range fields {
			row.TableCell[i].P = fields[i]
		}
		sheet.TableRow = append(sheet.TableRow, row)
	}

	data := ods.ODS{
		Content: ods.Content{
			 Body: ods.Body{
				Spreadsheet: ods.Spreadsheet{
					Table: []ods.Table{
						sheet,
					},
				}, 
			 },
		},
	}

	data.Meta.Meta.DocumentStatistic.CellCount = strconv.Itoa(len(books) * 6)

	buf := new(bytes.Buffer)
	files, err := zip.NewReader(bytes.NewReader(buf.Bytes()), MaxZipSize)
	if err != nil {
		return err
	}

	_ = files
	err = ods.WriteTo(w, data, &zip.ReadCloser{}) 
	if err != nil {
		return err
	}

	return nil
}

