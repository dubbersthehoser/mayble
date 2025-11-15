package controller

import (
	//"fmt"
	"errors"
	
	"github.com/dubbersthehoser/mayble/internal/core"
	"github.com/dubbersthehoser/mayble/internal/storage"
	"github.com/dubbersthehoser/mayble/internal/listing"
	"github.com/dubbersthehoser/mayble/internal/searching"
)

/*************************
	Book List
**************************/

const UnselectIndex int = -1

type BookList struct {
	broker         *broker.Broker
	list           []storage.BookLoan
	ordering       listring.Ordering
	orderBy        listing.OrderBy
	searchBy       searching.Field
	searchPattern  string
	SelectedIndex  int                // index element in list
	selection      searching.Ring
}
func NewBookList(b *broker.Broker) *BookList {
	b := &BookList{
		list:          make([]storage.BookLoan, 0),
		SelectedIndex: UnselectIndex,
		broker:          b,
	}
	b.searcher.Refresh(b.list)
	return b
}

func (l *BookList) SetOrderBy(by listing.OrderBy) {
	if err != nil {
		panic("invalid order by value")
	}
	l.orderBy = orderBy
	l.Search()
}

func (l *BookList) SetOrdering(o listing.Ordering) {
	l.ordering = o
}

func (l *BookList) Update() error {
	err := l.broker.Request(core.KeyRequestBookLoans, func(data any) {
		bookLoans, ok := data.([]data.BookLoan)
		if !ok {
			panic("request book loan given bad data")
		}
		l.list = listing.OrderBookLoans(bookLoans, l.orderBy, l.ordering)
	})
	if err != nil {
		return err
	}
	l.Unselect()
	l.searcher.Refresh(l.list)
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

func (l *BookList) Get(index int) (*listed.BookLoan, error) {
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
	selection := searching.SearchBookLoan(l.list, l.searchBy, l.searchPattern)
	if len(selection) != 0 {
		l.selection = searching.NewSelectionRing(selection)
		return true
	}
	return false
} 

func (l *BookList) SetSearchBy(by seaching.Field) {
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
