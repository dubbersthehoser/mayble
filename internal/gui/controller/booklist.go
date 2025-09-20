package controller

import (
	"errors"
	
	"github.com/dubbersthehoser/mayble/internal/core"
	"github.com/dubbersthehoser/mayble/internal/storage"
)

func GetTitle(b *storage.BookLoan) string {
	return b.BookData.Title
}
func GetAuthor(b *stroage.BookLoan) string {
	return b.BookData.Author
}
func GetGenre(b *storage.BookLoan) string {
	return b.BookData.Genre
}
func GetRatting(b *storage.BookLoan) string {
	return RattingToString(b.Ratting)
}
func GetLoanDate(b *storage.BookLoan) string {
	if b.IsOnLoan() {
		return DateToString(b.LoanDate)
	} else {
		return "n/a"
	}
}
func GetLoanName(b *storage.BookLoan) string {
	if b.IsOnLoan() {
		return b.LoanName
	} else {
		return "n/a"
	}
}

const UnselectIndex int = -1

type BookList struct {
	core           *core.Core
	list           []int64
	selectedIndex  int
}
func NewBookList(c *core.Core) *BookList {
	b := &BookList{
		list:          make([]int64),
		selectedIndex: UnselectIndex,
		core:          c,
	}
	return b
}

func (l *BookList) Len() int {
	return len(l.List)
}

func (l *BookList) Select(index int) error {
	if err := l.ValidateIndex(index); err != nil {
		return err
	}
	l.SelectedIndex = index
	return nil
}

func (l *BookList) Unselect() {
	l.selected = UnselectIndex
}

func (l *BookList) ValidateIndex(index int) error {
	if len(l.list) <= index || index < 0 {
		return errors.New("index out of range")
	}
	return nil
}

func (l *BookList) Selected() (*storage.BookLoan, error) {
	if l.selected == UnselectIndex {
		return nil, errors.New("booklist: no selected book")
	}
	if err := l.ValidateIndex(l.selected); err != nil {
		return nil, err
	}
	return l.Get(l.selected)
}

func (l *BookList) Get(index int) (*storage.BookLoan, error) {
	if err := l.ValidateIndex(index); err != nil {
		return nil, err
	}
	bookID := l.list[index]
	bookLoan, err := l.core.GetBookLoan(bookID)
	if err != nil {
		return nil, err
	}
	return &bookLoan, nil
}
