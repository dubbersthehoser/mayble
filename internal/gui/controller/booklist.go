package controller

import (
	"errors"
	
	"github.com/dubbersthehoser/mayble/internal/core"
	"github.com/dubbersthehoser/mayble/internal/storage"
)


type BookLoanListed struct {
	Title    string
	Author   string
	Genre    string
	Ratting  string
	Borrower string
	Date     string
}


func toBookLoanView(bookLoan *storage.BookLoan) *BookLoanListed {
	view := BookLoanListed{
		Title:   bookLoan.Title,
		Author:  bookLoan.Author,
		Genre:   bookLoan.Genre,
		Ratting: rattingToString(bookLoan.Ratting),
		Borrower: "n/a",
		Date:     "n/a",
	}
	if bookLoan.IsOnLoan() {
		view.Borrower = bookLoan.Loan.Name
		view.Date  = dateToString(&bookLoan.Loan.Date)
	}
	return &view
}


/*************************
	Book List
**************************/

const UnselectIndex int = -1

type BookList struct {
	order          core.Order
	orderBy        core.OrderBy
	core           *core.Core
	list           []storage.BookLoan
	selectedIndex  int
}
func NewBookList(c *core.Core) *BookList {
	b := &BookList{
		list:          make([]storage.BookLoan, 0),
		selectedIndex: UnselectIndex,
		core:          c,
	}
	return b
}

func (l *BookList) List() error {
	bookLoans, err := l.core.ListBookLoans(core.ByTitle, core.ASC)
	if err != nil {
		return err
	}
	l.Unselect()
	l.list = bookLoans
	return nil
}

func (l *BookList) Len() int {
	return len(l.list)
}

func (l *BookList) Select(index int) error {
	if err := l.ValidateIndex(index); err != nil {
		return err
	}
	l.selectedIndex = index
	return nil
}

func (l *BookList) Unselect() {
	l.selectedIndex = UnselectIndex
}

func (l *BookList) ValidateIndex(index int) error {
	if len(l.list) <= index || index < 0 {
		return errors.New("index out of range")
	}
	return nil
}

func (l *BookList) Selected() (*BookLoanListed, error) {
	if l.selectedIndex == UnselectIndex {
		return nil, errors.New("booklist: no book selected")
	}
	if err := l.ValidateIndex(l.selectedIndex); err != nil {
		return nil, err
	}
	bookListed, err := l.Get(l.selectedIndex)
	if err != nil {
		return nil, err
	}
	return bookListed, nil
}

func (l *BookList) Get(index int) (*BookLoanListed, error) {
	if err := l.ValidateIndex(index); err != nil {
		return nil, err
	}
	bookLoan := l.list[index]
	bookView := toBookLoanView(&bookLoan)
	return bookView, nil
}














