package data

import (
	"time"
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



/********************
	Builder
*********************/

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

// Build with builder. Returns BookLoan and a validation error.
//func (b *BookLoanBuilder) Build() (*BookLoan, error) {
//	bl := NewBookLoan()
//	bl.Title = b.title
//	bl.Author = b.author
//	bl.Genre = b.genre
//	bl.Ratting = b.ratting
//	if b.date != nil && b.borrower != "" {
//		bl.Loan.Borrower = b.borrower
//		bl.Loan.Date = *b.date
//	} else {
//		bl.UnsetLoan()
//	}
//	return bl, ValidateBookLoan(bl)
//}

// BookLoanToBuilder convers a BookLoan type to a builder for reconstruction.
//func BookLoanToBuilder(bl *BookLoan) *BookLoanBuilder {
//	builder := NewBookLoanBuilder()
//	builder.bookID = bl.Book.ID
//	if bl.IsOnLoan() {
//		builder.loanID = bl.Loan.ID
//		builder.WithTitle(bl.Title).
//			WithAuthor(bl.Author).
//			WithGenre(bl.Genre).
//			WithRatting(bl.Ratting).
//			WithBorrower(bl.Loan.Borrower).
//			WithDate(bl.Loan.Date)
//	} else {
//		builder.WithTitle(bl.Title).
//			WithAuthor(bl.Author).
//			WithGenre(bl.Genre).
//			WithRatting(bl.Ratting)
//	}
//	return builder
//}



