package controller

import (
	//"fmt"
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

func toBookLoanListed(bookLoan *storage.BookLoan) *BookLoanListed {
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

type Ordering core.Order
const DEC Ordering = Ordering(core.DEC)
const ASC Ordering = Ordering(core.ASC)

type BookList struct {
	ordering       Ordering
	orderBy        core.OrderBy
	core           *core.Core
	list           []storage.BookLoan
	SelectedIndex  int
	searcher       *Searcher
}
func NewBookList(c *core.Core) *BookList {
	b := &BookList{
		list:          make([]storage.BookLoan, 0),
		SelectedIndex: UnselectIndex,
		core:          c,
		searcher:      NewSearcher(),
	}
	b.searcher.Refresh(b.list)
	return b
}

func SortByList() []string {
	return []string{"Title", "Author", "Genre", "Ratting", "Borrower", "Date"}
}

func (l *BookList) SetOrderBy(by string) {
	switch by {
	case "Title":
		l.orderBy = core.ByTitle
	case "Author":
		l.orderBy = core.ByAuthor
	case "Genre":
		l.orderBy = core.ByGenre
	case "Ratting":
		l.orderBy = core.ByRatting
	case "Borrower":
		l.orderBy = core.ByBorrower
	case "Date":
		l.orderBy = core.ByDate
	default:
		panic("invalid order by value")
	}
	l.Search()
}
func (l *BookList) SetOrdering(o Ordering) {
	l.ordering = o
}


func (l *BookList) Update() error {
	bookLoans, err := l.core.ListBookLoans(l.orderBy, core.Order(l.ordering))
	if err != nil {
		return err
	}
	l.Unselect()
	l.list = bookLoans
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

func (l *BookList) Selected() (*storage.BookLoan, error) {
	if l.SelectedIndex == UnselectIndex {
		return nil, errors.New("booklist: no book selected")
	}
	if err := l.ValidateIndex(l.SelectedIndex); err != nil {
		return nil, err
	}
	bookLoan := l.list[l.SelectedIndex]
	return &bookLoan, nil
}

func (l *BookList) Get(index int) (*BookLoanListed, error) {
	if err := l.ValidateIndex(index); err != nil {
		return nil, err
	}
	bookLoan := l.list[index]
	bookView := toBookLoanListed(&bookLoan)
	return bookView, nil
}

func (l *BookList) SetSearch(pattern string) {
	if pattern == "" {
		l.SearchUndo()
		return 
	}
	l.searcher.SetSearch(pattern)
}

func (l *BookList) Search() bool {
	l.searcher.Search()
	if l.searcher.IsSelection() {
		l.Select(l.searcher.Selected())
		return true
	} else {
		l.Unselect()
		return false
	}
} 

func (l *BookList) SetSearchBy(by Field) {
	l.searcher.SetByField(by)
	l.searcher.Refresh(l.list)
}

func (l *BookList) SearchUndo() {
	l.searcher.Refresh(l.list)
}

func (l *BookList) SelectNext() {
	if !l.IsSelected() {
		err := l.Select(0)
		if err != nil {
			return 
		}
	}
	l.searcher.SelectedNext()
	l.Select(l.searcher.Selected())
}

func (l *BookList) SelectPrev() {
	if !l.IsSelected() {
		err := l.Select(0)
		if err != nil {
			return
		}
	}
	l.searcher.SelectedPrev()
	l.Select(l.searcher.Selected())
}
