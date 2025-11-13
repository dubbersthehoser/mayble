package data

import (
	"time"
	"errors"
)

const ZeroID int64 = -127 

type Book struct {
	ID      int64  
	Title   string 
	Author  string 
	Genre   string 
	Ratting int    
}

type Loan struct {
	ID       int64     
	Borrower string   
	Date     time.Time 
}

type BookLoan struct {
	Book
	Loan *Loan
}

func NewBookLoan() *BookLoan {
	b := &BookLoan{
		Book: Book{
			ID: ZeroID,
		},
		Loan: &Loan{
			ID: ZeroID,
		},
	}
	return b
}

func (bl *BookLoan) IsOnLoan() bool {
	return bl.Loan != nil
}
func (bl *BookLoan) UnsetLoan() {
	bl.Loan = nil
}

type BookLoanBuilder struct {
	bookID   int64
	title    string
	author   string
	genre    string
	ratting  int
	loanID   int64
	borrower string
	date     *time.Time
}
func NewBookLoanBuilder() *BookLoanBuilder{
	return &BookLoanBuilder{
		bookID: ZeroID,
		loanID: ZeroID,
	}
}

func (b *BookLoanBuilder) WithTitle(s string) *BookLoanBuilder {
	b.title = s
	return b
}
func (b *BookLoanBuilder) WithAuthor(s string) *BookLoanBuilder {
	b.author = s
	return b
}
func (b *BookLoanBuilder) WithGenre(s string) *BookLoanBuilder {
	b.genre = s
	return b
}
func (b *BookLoanBuilder) WithRatting(r int) *BookLoanBuilder {
	b.ratting = r
	return b
}
func(b *BookLoanBuilder) WithBorrower(s string) *BookLoanBuilder {
	b.borrower = s
	return b
}
func (b *BookLoanBuilder) WithDate(t time.Time) *BookLoanBuilder {
	b.date = &t
	return b
}
func (b *BookLoanBuilder) Build() (*BookLoan, error) {

	bl := NewBookLoan()
	bl.Title = b.title
	bl.Author = b.author
	bl.Genre = b.genre
	bl.Ratting = b.ratting
	if b.date != nil && b.borrower != "" {
		bl.Loan.Borrower = b.borrower
		bl.Loan.Date = *b.date
	} else {
		bl.UnsetLoan()
	}

	var (
		InvalidTitle        bool = !ValidTitle(b.title)
		InvalidAuthor       bool = !ValidAuthor(b.author)
		InvalidGenre        bool = !ValidGenre(b.genre)
		InvalidRatting      bool = !ValidRatting(b.ratting)
		InvalidBorrowerName bool = !ValidBorrowerName(b.borrower)
		InvalidBorrowDate   bool = !ValidBorrowDate(b.date)
	)
	switch {
	case InvalidTitle:
		return nil, errors.New("data: invalid title")

	case InvalidAuthor:
		return nil, errors.New("data: invalid author")

	case InvalidGenre:
		return nil, errors.New("data: invalid genre")

	case InvalidRatting:
		return nil, errors.New("data: invalid ratting")

	case InvalidBorrowerName && bl.IsOnLoan():
		return nil , errors.New("data: invalid borrower name")

	case InvalidBorrowDate && bl.IsOnLoan():
		return nil, errors.New("data: invalid borrow date")
	}
	return bl, nil
}

func BookLoanToBuilder(bl *BookLoan) *BookLoanBuilder {
	builder := NewBookLoanBuilder()
	builder.bookID = bl.Book.ID
	if bl.IsOnLoan() {
		builder.loanID = bl.Loan.ID
		builder.WithTitle(bl.Title).
			WithAuthor(bl.Author).
			WithGenre(bl.Genre).
			WithRatting(bl.Ratting).
			WithBorrower(bl.Loan.Borrower).
			WithDate(bl.Loan.Date)
	} else {
		builder.WithTitle(bl.Title).
			WithAuthor(bl.Author).
			WithGenre(bl.Genre).
			WithRatting(bl.Ratting)
	}
	return builder
}



/***********************
	Validation
************************/

func ValidID(id int64) bool {
	return id >= 0
}
func ValidTitle(s string) bool {
	return "" != s
}
func ValidAuthor(s string) bool {
	return "" != s
}
func ValidGenre(s string) bool {
	return "" != s
}
func ValidRatting(ratting int) bool {
	return ratting >= 0 && ratting < 6
}
func ValidBorrowerName(s string) bool {
	return "" != s
}
func ValidBorrowDate(date *time.Time) bool {
	if date == nil {
		return false
	}
	return !date.IsZero()
}




