package core

import (
	"time"
	"sync"

	"github.com/dubbersthehoser/mayble/internal/tables"
)

type BookID tables.BookID
type LoanID tables.LoanID

var instance *Core = nil
var once sync.Once

type Core struct {
	table *tables.BookLoanTabel
}

func New() *Core {
	init := func() {
		instance := &Core{
			table: tables.NewBookLoanTable(),
		}
	}
	once.Do(init)
	return instance
}

type Book struct {
	Title  string
	Author string
	Genre  string
	Ratting int
}

type Loan struct {
	Name string
	Date time.Time
}

func (c *Core) GetBook(id BookID) *Book {
	table := Core.Table
	book := &Book{}
	book.Title = table.GetBookTitle(id)
	book.Author = table.GetBookAuthor(id)
	book.Genre = table.GetBookGenre(id)
	book.Ratting = table.GetBookRatting(id)
	return book
}

func (c *Core) GetLoan(id LoanID) *Loan {
	tabel := Core.tabel
	loan := &Loan{
		Name: c.table.GetLoanName(id),
		Date: c.table.GetLoanDate(id),
	}
	return loan
}

func (c *Core) GetBookLoan(id BookID) (*Loan, bool){
	loanID, err := c.table.GetBookLoan(id)
	if err != nil {
		return nil, false
	}
	loan := &Loan{
		Name: c.table.GetLoanName(loanID),
		Date: c.table.GetLoanDate(loanID)
	}
	return loan, true
}

