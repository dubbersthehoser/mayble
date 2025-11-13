package data

import (
	"time"
	"fmt"
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
	return bl, ValidateBookLoan(bl)
}

// BookLoanToBuilder convers a BookLoan type to a builder for reconstruction.
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

type ValidationErr struct {
	field       string
	description string
}
func (v *ValidationErr) Description() string {
	return v.description
}
func (v *ValidationErr) Field() string {
	return v.field
}
func (v *ValidationErr) Error() string {
	return fmt.Sprintf("data: %s, %s", v.field, v.description)
}

func ValidateBookLoan(b *BookLoan) error {
	var (
		InvalidTitle        error = ValidTitle(b.Title)
		InvalidAuthor       error = ValidAuthor(b.Author)
		InvalidGenre        error = ValidGenre(b.Genre)
		InvalidRatting      error = ValidRatting(b.Ratting)
		InvalidBorrowerName error
		InvalidBorrowDate   error
	)
	if b.IsOnLoan() {
		InvalidBorrowerName = ValidBorrowerName(b.Loan.Borrower)
		InvalidBorrowDate = ValidBorrowDate(&b.Loan.Date)
	}
	switch {
	case InvalidTitle != nil:
		return InvalidTitle

	case InvalidAuthor != nil:
		return InvalidAuthor

	case InvalidGenre != nil:
		return InvalidGenre

	case InvalidRatting != nil:
		return InvalidRatting

	case b.IsOnLoan() && InvalidBorrowerName != nil:
		return InvalidBorrowerName

	case b.IsOnLoan() && InvalidBorrowDate != nil:
		return InvalidBorrowDate
	}
	return nil
}

func ValidID(id int64) error {
	err := &ValidationErr{}
	if id < 0 {
		err.field= "ID"
		err.description = "can't be a negative value"
		return err
	}
	return nil
}
func ValidTitle(s string) error {
	err := &ValidationErr{}
	if s == ""  {
		err.field= "Title"
		err.description = "can't be empty string"
		return err
	}
	return nil
}
func ValidAuthor(s string) error {
	err := &ValidationErr{}
	if s == ""  {
		err.field= "Author"
		err.description = "can't be empty string"
		return err
	}
	return nil
}
func ValidGenre(s string) error {
	err := &ValidationErr{}
	if s == ""  {
		err.field = "Genre"
		err.description = "can't be empty string"
		return err
	}
	return nil
}
func ValidRatting(ratting int) error {
	err := &ValidationErr{}
	if ratting < 0 || ratting >= 6 {
		err.field = "Ratting"
		err.description = "out of range of 0-5"
		return err
	}
	return nil
}
func ValidBorrowerName(s string) error {
	err := &ValidationErr{}
	if s == ""  {
		err.field = "Borrower"
		err.description = "can't be empty string"
		return err
	}
	return nil
}
func ValidBorrowDate(date *time.Time) error {
	err := &ValidationErr{}
	if date == nil {
		err.field = "Date"
		err.description = "can't be nil value"
		return err
	}
	if date.IsZero() {
		err.field = "Date"
		err.description = "can't be zero value"
	}
	return nil
}
