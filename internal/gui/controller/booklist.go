package controller

import (
	//"fmt"
	
	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/listing"
	"github.com/dubbersthehoser/mayble/internal/searching"
	"github.com/dubbersthehoser/mayble/internal/emiter"
	"github.com/dubbersthehoser/mayble/internal/gui"
)

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
		selection: nil,
		list: list,
		broker: b,
	}

	bs.broker.Subscribe(&emiter.Listener{
		Handler: func(e *emiter.Event) {
			switch e.Name {

			case gui.EventSearchPattern:
				pattern := e.Data.(string)
				bs.pattern = pattern
				bs.Search()

			case gui.EventSearchBy:
				by := e.Data.(searching.Field)
				bs.by = by

			case gui.EventSelectNext:
				if bs.selection == nil || len(*bs.list) == 0 {
					return
				}
				idx := bs.selection.Next()
				bs.broker.Notify(emiter.Event{
					Name: gui.EventEntrySelected,
					Data: idx,
				})

			case gui.EventSelectPrev:
				if bs.selection == nil || len(*bs.list) == 0 {
					return
				}
				idx := bs.selection.Prev()
				bs.broker.Notify(emiter.Event{
					Name: gui.EventEntrySelected,
					Data: idx,
				})

			case gui.EventSearch:
				bs.Search()

			case gui.EventSelectionAll:
				bs.selection = searching.NewRangeRing(len(*bs.list))
			}
		},
	},
		gui.EventSearchPattern,
		gui.EventSearchBy,
		gui.EventSelectNext,
		gui.EventSelectPrev,
		gui.EventSearch,
		gui.EventSelectionAll,
	)
	return bs
}

func (bs *BookLoanSearcher) HasSelection() bool {
	return bs.selection != nil 
}


func (bs *BookLoanSearcher) Search() {
	selection := searching.SearchBookLoans(*bs.list, bs.by, bs.pattern)
	if len(selection) == 0 || bs.pattern == "" {
		bs.broker.Notify(emiter.Event{
			Name: gui.EventSelectionNone,
		})
		bs.selection = nil
		return
	}
	bs.selection = searching.NewSelectionRing(selection)
	idx := bs.selection.Selected()
	bs.broker.Notify(emiter.Event{
		Name: gui.EventSelection,
	})
	bs.broker.Notify(emiter.Event{
		Name: gui.EventEntrySelected,
		Data: idx,
	})
}


type BookLoanList struct {
	app      *app.App
	list     []app.BookLoan
	broker   *emiter.Broker
	selected int // index item in list. When index < 0 then it's unselected.
	orderBy  listing.OrderBy
	ordering listing.Ordering
}

func NewBookLoanList(b *emiter.Broker, a *app.App) *BookLoanList {

	list, err := a.GetBookLoans()
	if err != nil {
		b.Notify(emiter.Event{
			Name: gui.EventDisplayErr,
			Data: err,
		})
	}

	bl := &BookLoanList{
		app:    a,
		broker: b,   
		selected: -1,
		list: list,
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

			case gui.EventEntrySelected:
				idx, ok := e.Data.(int)
				if !ok {
					panic("event: entry selected: data is not int")
				}
				if idx >= len(bl.list) || idx < 0 {
					panic("event: entry selected: invalid index")
				}
				bl.selected = idx

			case gui.EventEntryUnselected:
				bl.selected = -1

			case gui.EventDocumentModified:
				l, err := bl.app.GetBookLoans()
				if err != nil {
					bl.broker.Notify(emiter.Event{
						Name: gui.EventDisplayErr,
						Data: err,
					})
					return
				}
				bl.list = l
				bl.Sort()
			}
		},
	},
		gui.EventListOrdering,
		gui.EventListOrderBy,
		gui.EventEntryUnselected,
		gui.EventEntrySelected,
		gui.EventDocumentModified,
	)
	return bl
}

func (bl *BookLoanList) Selected() *app.BookLoan {
	if bl.selected < 0 || len(bl.list) <= bl.selected {
		panic("booklist: selected index is invalid")
	}
	bookLoan := bl.list[bl.selected]
	return &bookLoan
}

func (bl *BookLoanList) HasSelected() bool {
	return bl.selected >= 0
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
