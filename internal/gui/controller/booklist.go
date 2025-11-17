package controller

import (
	"errors"
	
	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/data"
	"github.com/dubbersthehoser/mayble/internal/listing"
	"github.com/dubbersthehoser/mayble/internal/searching"
)

/*************************
	Book List
**************************/

const UnselectIndex int = -1

type BookList struct {
	app            app.BookLoaning
	list           []data.BookLoan
	ordering       listing.Ordering
	orderBy        listing.OrderBy
	searchBy       searching.Field
	searchPattern  string
	SelectedIndex  int // selected element in list
	selection      searching.Ring
}
func NewBookList(app app.BookLoaning) *BookList {
	b := &BookList{
		list:          make([]data.BookLoan, 0),
		SelectedIndex: UnselectIndex,
		app:           app,
	}
	return b
}

func (l *BookList) SetOrderBy(by listing.OrderBy) {
	l.orderBy = by
	l.Search()
}

func (l *BookList) SetOrdering(o listing.Ordering) {
	l.ordering = o
}

func (l *BookList) Update() error {
	bookLoans, err := l.app.GetBookLoans()
	if err != nil {
		return err
	}
	l.list = bookLoans
	l.Unselect()
	return nil
}

func (l *BookList) Len() int {
	return len(l.list)
}

func (l *BookList) Select(index int) error {
	if err := l.ValidateIndex(index); err != nil {
		return err
	}
	l.SelectedIndex = index
	return nil
}

func (l *BookList) Unselect() {
	l.SelectedIndex = UnselectIndex
}

func (l *BookList) ValidateIndex(index int) error {
	if len(l.list) <= index || index < 0 {
		return errors.New("index out of range")
	}
	return nil
}

func (l *BookList) IsSelected() bool {
	return l.SelectedIndex != UnselectIndex
}

func (l *BookList) Selected() (*data.BookLoan, error) {
	if l.SelectedIndex == UnselectIndex {
		return nil, errors.New("booklist: no book selected")
	}
	if err := l.ValidateIndex(l.SelectedIndex); err != nil {
		return nil, err
	}
	bookLoan := l.list[l.SelectedIndex]
	return &bookLoan, nil
}

func (l *BookList) Get(index int) (*listing.BookLoan, error) {
	if err := l.ValidateIndex(index); err != nil {
		return nil, err
	}
	bookLoan := l.list[index]
	bookView := listing.BookLoanToListed(&bookLoan)
	return bookView, nil
}


func (l *BookList) SetSearch(pattern string) {
	l.searchPattern = pattern
}

func (l *BookList) Search() bool {
	if l.searchPattern == "" {
		l.selection = searching.NewRangeRing(len(l.list))
		return true
	}
	selection := searching.SearchBookLoans(l.list, l.searchBy, l.searchPattern)
	if len(selection) != 0 {
		l.selection = searching.NewSelectionRing(selection)
		return true
	}
	return false
} 

func (l *BookList) SetSearchBy(by searching.Field) {
	l.searchBy = by
}

func (l *BookList) SelectNext() {
	if !l.IsSelected() {
		err := l.Select(0)
		if err != nil {
			return 
		}
	}
	l.Select(l.selection.Next())
}

func (l *BookList) SelectPrev() {
	if !l.IsSelected() {
		err := l.Select(0)
		if err != nil {
			return
		}
	}
	l.Select(l.selection.Prev())
}
