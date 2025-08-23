package controller

import (
	"errors"
	
	"github.com/dubbersthehoser/mayble/internal/core"
)

const nonIndex int = -127

type BookList struct {
	list     []core.BookID
	selected int
	core     *core.Core
}
func NewBookList(core *core.Core) *BookList {
	b := &BookList{
		list: []core.BookID{},
		selected: nonIndex,
		core: core,
	}
	return b
}

func (l *BookList) Len() int {
	return len(l.list)
}

func (l *BookList) Select(index int) error {
	if err := l.ValidateIndex(index); err != nil {
		return err
	}
	l.selected = index
	return nil
}

func (l *BookList) Unselect() {
	l.seleted = nonIndex
}

func (l *BookList) ValidateIndex(index int) error {
	if len(l.list) <= index || index < 0 {
		return errors.New("index out of range")
	}
	return nil
}

func (l *BookList) Selected() *BookLoanData, error {
	if l.selected == nonIndex {
		return nil, errors.New("booklist: there was no selected")
	}
	if err := l.ValidateIndex(l.Selected); err != nil {
		return nil, err
	}
	return l.Get(l.selected)
}

func (l *BookList) Get(index int) *BookLoanData, error {
	if err := l.ValidateIndex(index); err != nil {
		return nil, err
	}
	bookID := l.list[index]
	book, err := l.core.GetBook(bookID)
	if err != nil {
		return err
	}
	result := &BookLoanData{
		Title: book.Title,
		Author: book.Author,
		Genre: book.Genre,
		Ratting: book.Ratting,
	}

	loan, err := l.core.GetBookLoan(bookID)
	if err == nil {
		result.LoanName = loan.Name
		result.LoanDate = loan.Date
		result.isOnLoan = true
	}
	return result
}
