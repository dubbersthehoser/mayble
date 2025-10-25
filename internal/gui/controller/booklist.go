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
	selectedIndex  int
	searcher       *Searcher
}
func NewBookList(c *core.Core) *BookList {
	b := &BookList{
		list:          make([]storage.BookLoan, 0),
		selectedIndex: UnselectIndex,
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

func (l *BookList) IsSelected() bool {
	return l.selectedIndex != UnselectIndex
}

func (l *BookList) Selected() (*storage.BookLoan, error) {
	if l.selectedIndex == UnselectIndex {
		return nil, errors.New("booklist: no book selected")
	}
	if err := l.ValidateIndex(l.selectedIndex); err != nil {
		return nil, err
	}
	bookLoan := l.list[l.selectedIndex]
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

func (l *BookList) Search(pattern string) {
	l.searcher.Search(pattern)
	result := l.searcher.selection
	if len(result) != 0 {
		l.Select(result[0])
	}
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

	if (l.searcher.selection != nil) && (len(l.searcher.selection) > 0) {
		idx := (l.searcher.selected + 1) % len(l.searcher.selection)
		idx = l.searcher.selection[idx]
		err := l.Select(idx)
		if err != nil {
			return
		}
		return
	}

	idx := (l.selectedIndex + 1) % len(l.list)
	err := l.Select(idx)
	if err != nil {
		return
	}
	return


}

func (l *BookList) SelectPrev() {
	if !l.IsSelected() {
		err := l.Select(len(l.list)-1)
		if err != nil {
			return
		}
	}

	if (l.searcher.selection != nil) || (len(l.searcher.selection) > 0) {
		idx := (l.searcher.selected - 1) % len(l.searcher.selection)
		idx = l.searcher.selection[idx]
		err := l.Select(idx)
		if err != nil {
			return
		}
		return
	}

	idx := (l.selectedIndex - 1) % len(l.list)
	err := l.Select(idx)
	if err != nil {
		return
	}
	return


}
