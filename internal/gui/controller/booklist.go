package controller

import (
	"errors"
	"log"
	//"fmt"
	
	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/listing"
	"github.com/dubbersthehoser/mayble/internal/searching"
)

const UnselectIndex int = -1

type BookList struct {
	app            app.BookLoaning
	list           []app.BookLoan
	ordering       listing.Ordering
	orderBy        listing.OrderBy
	searchBy       searching.Field
	searchPattern  string
	SelectedIndex  int // selected element in list
	selection      searching.Ring
}
func NewBookList(a app.BookLoaning) *BookList {
	b := &BookList{
		list:          make([]app.BookLoan, 0),
		SelectedIndex: UnselectIndex,
		app:           a,
		selection:     searching.NewRangeRing(0),

		// needs to be the same to the init state of table header
		ordering: listing.ASC,   
		orderBy:  listing.ByTitle,

	}
	b.Update()
	return b
}

func (l *BookList) Update() error {
	bookLoans, err := l.app.GetBookLoans()
	if err != nil {
		return err
	}
	l.list = listing.OrderBookLoans(bookLoans, l.orderBy, l.ordering)
	l.Unselect()
	l.selection = searching.NewRangeRing(len(bookLoans))
	return nil
}

func (l *BookList) Len() int {
	return len(l.list)
}

func (l *BookList) SetOrderBy(by listing.OrderBy) {
	l.orderBy = by
}
func (l *BookList) OrderBy() listing.OrderBy {
	return l.orderBy
}

func (l *BookList) SetOrdering(o listing.Ordering) {
	l.ordering = o
}
func (l *BookList) Ordering() listing.Ordering {
	return l.ordering
}

func (l *BookList) SetSearchPattern(pattern string) {
	l.searchPattern = pattern
}

func (l *BookList) SetSearchBy(by searching.Field) {
	l.searchBy = by
}

func (l *BookList) Search() bool {
	if l.searchPattern == "" {
		l.selection = searching.NewRangeRing(len(l.list))
		return true
	}
	selection := searching.SearchBookLoans(l.list, l.searchBy, l.searchPattern)
	println(selection)
	if len(selection) != 0 {
		l.selection = searching.NewSelectionRing(selection)
		return true
	}
	return false
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
	return l.SelectedIndex > UnselectIndex
}

func (l *BookList) Selected() (*app.BookLoan, error) {
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

func (l *BookList) SelectNext() {
	var err error
	if !l.IsSelected()  {
		err = l.Select(l.selection.Selected())
	} else {
		err = l.Select(l.selection.Next())
		println(l.SelectedIndex)
	}
	if err != nil {
		log.Println("booklist: unexpected error: ", err)
	}
}

func (l *BookList) SelectPrev() {
	var err error
	if !l.IsSelected() {
		err = l.Select(l.selection.Selected())
	} else {
		err = l.Select(l.selection.Prev())
	}
	if err != nil {
		log.Println("booklist: unexpected error: ", err)
	}
}
