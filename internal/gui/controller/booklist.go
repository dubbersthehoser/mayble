package controller

import (
	"errors"
	"log"
	//"fmt"
	
	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/listing"
	"github.com/dubbersthehoser/mayble/internal/searching"
	"github.com/dubbersthehoser/mayble/internal/emiter"
	"github.com/dubbersthehoser/mayble/internal/gui"
)

const UnselectIndex int = -1

type BookLoanSearcher struct {
	pattern   string
	by        searching.Field
	selected  int
	selection searching.Ring
	list      *[]app.BookLoan
	broker    *emiter.Broker
}

func NewBookLoanSearcher(b *emiter.Broker, list *[]app.BookLoan) *BookLoanSearcher {
	bs := &BookLoanSearcher{
		selected: UnselectIndex,
		selection: searching.NewRangeRing(len(*list)),
		list: list,
		broker: b,
	}

	bs.broker.Subscribe(&emiter.Listener{
		Handler: func(e *emiter.Event) {
			switch e.Name {

			case gui.EventSearchPattern:
				pattern := e.Data.(string)
				bs.pattern = pattern

			case gui.EventSearchBy:
				by := e.Data.(searching.Field)
				bs.by = by

			case gui.EventSelectNext:
				idx := bs.selection.Next()
				bs.broker.Notify(emiter.Event{
					Name: gui.EventEntrySelected,
					Data: idx,
				})

			case gui.EventSelectPrev:
				idx := bs.selection.Prev()
				bs.broker.Notify(emiter.Event{
					Name: gui.EventEntrySelected,
					Data: idx,
				})

			case gui.EventSearch:
				selection := searching.SearchBookLoans(*bs.list, bs.by, bs.pattern)
				if len(selection) == 0 {
					return
				}
				bs.selection = searching.NewSelectionRing(selection)
				idx := bs.selection.Selected()
				bs.broker.Notify(emiter.Event{
					Name: gui.EventEntrySelected,
					Data: idx,
				})
			}
		},
	},
		gui.EventSearchPattern,
		gui.EventSearchBy,
		gui.EventSelectNext,
		gui.EventSelectPrev,
		gui.EventSearch,
	)
	return bs
}
func (bs *BookLoanSearcher) HasSelected() bool {
	return bs.selected > -1
}

func (bs *BookLoanSearcher) Search() {
	if bs.pattern == "" {
		bs.selection = searching.NewRangeRing(len(*bs.list))
		return
	}
	selection := searching.SearchBookLoans(*bs.list, bs.by, bs.pattern)
	if len(selection) != 0 {
		bs.selection = searching.NewSelectionRing(selection)
	}
}


type BookLoanList struct {
	app      app.BookLoaning
	list     []app.BookLoan
	broker   *emiter.Broker
	orderBy  listing.OrderBy
	ordering listing.Ordering
}

func NewBookLoanList(a app.BookLoaning, b *emiter.Broker) *BookLoanList {
	bl := &BookLoanList{
		app:    a,
		broker: b,   
		list:   make([]app.BookLoan, 0),
	}

	bl.broker.Subscribe(&emiter.Listener{
		Handler: func(e *emiter.Event) {
			switch e.Name {
			case gui.EventListOrderBy:
				by := e.Data.(listing.OrderBy)
				bl.orderBy = by

			case gui.EventListOrdering:
				o := e.Data.(listing.Ordering)
				bl.ordering = o
				bl.Sort()
			}
		},
	},
		gui.EventListOrdering,
		gui.EventListOrderBy,
	)
	return bl
}

func (bl *BookLoanList) Sort() {
	bl.list = listing.OrderBookLoans(bl.list, bl.orderBy, bl.ordering)
	bl.broker.Notify(emiter.Event{
		Name: gui.EventListOrdered,
	})
}

func (bl *BookLoanList) Get(index int) *listing.BookLoan {
	bookLoan := bl.list[index]
	bookView := listing.BookLoanToListed(&bookLoan)
	return bookView
}

func (bl *BookLoanList) Len() int {
	return len(bl.list)
}


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
