package controller

import (
	"time"
	"fmt"
	"errors"

	"github.com/dubbersthehoser/mayble/internal/core"
	"github.com/dubbersthehoser/mayble/internal/storage"
)

type EditType int
const (
	Updating EditType = iota
	Creating 
	Deleting
)
func (t EditType) String() string {
	switch t {
	case Updating:
		return "Updating"
	case Creating:
		return "Creating"
	case Deleting:
		return "Deleting"
	default:
		panic("invalid builder build type")
	}
	return ""
}

type BookEditor struct {
	core *core.Core
}

func NewBookEditor(c *core.Core) *BookEditor {
	return &BookEditor{
		core: c,
	}
}

func (be *BookEditor) Submit(builder *BookLoanBuilder) error {
	bookLoan := builder.Build()
	switch builder.Type {
	case Updating:
		return be.core.UpdateBookLoan(bookLoan)
	case Creating:
		return be.core.CreateBookLoan(bookLoan)
	case Deleting:
		return be.core.DeleteBookLoan(bookLoan)
	default:
		return fmt.Errorf("submit type not found: %s", builder.Type)
	}
	return nil
}

func NewBookLoanBuilder() *BookLoanBuilder {
	return &BookLoanBuilder{
		Type: Creating,
	}
} 

func NewBuilderWithBookLoan(b *storage.BookLoan) *BookLoanBuilder {
	builder := NewBookLoanBuilder()
	builder.id = b.ID
	builder.SetToUpdate()
	builder.SetTitle(b.Title)
	builder.SetAuthor(b.Author)
	builder.SetGenre(b.Genre)
	builder.SetRatting(b.Ratting)
	
	if b.IsOnLoan() {
		builder.SetIsOnLoan(true)
		builder.SetBorrower(b.Loan.Name)
		builder.SetDate(&b.Loan.Date)
	}

	return builder
}

type BookLoanBuilder struct {
	Type     EditType
	id       int64
	Title    string
	Author   string
	Genre    string
	Ratting  int 
	IsOnLoan bool
	Borrower string
	Date     time.Time
}

func (b *BookLoanBuilder) Validate() error {
	if err := ValidateTitle(b.Title); err != nil {
		return errors.New("must have an title.")
	}
	if err := ValidateAuthor(b.Author); err != nil {
		return errors.New("must have an author.")
	}
	if err := ValidateGenre(b.Genre); err != nil {
		return errors.New("must have an genre.")
	}
	if !b.IsOnLoan {
		return nil
	}
	if err := ValidateLoanName(b.Borrower); err != nil {
		return errors.New("must have borrower.")
	}
	if err := ValidateLoanDate(&b.Date); err != nil {
		return errors.New("must have borrow date.")
	}
	return nil
}

func (b *BookLoanBuilder) SetToDelete() {
	b.Type = Deleting
}
func (b *BookLoanBuilder) SetToUpdate() {
	b.Type = Updating
}
func (b *BookLoanBuilder) SetToCreate() {
	b.Type = Creating
}


func (b *BookLoanBuilder) SetIsOnLoan(onLoan bool) {
	b.IsOnLoan = onLoan
}
func (b *BookLoanBuilder) SetTitle(title string) {
	b.Title = title
}
func (b *BookLoanBuilder) SetAuthor(author string) {
	b.Author = author
}
func (b *BookLoanBuilder) SetGenre(genre string) {
	b.Genre = genre
}
func (b *BookLoanBuilder) SetRatting(ratting int){
	b.Ratting = ratting
}
// SetRattingAsString sets ratting from string value returned by RattingToStirng
func (b *BookLoanBuilder) SetRattingAsString(ratting string){
	b.Ratting = RattingToInt(ratting)
}
func (b *BookLoanBuilder) SetBorrower(name string) {
	b.Borrower = name
}
func (b *BookLoanBuilder) SetDate(date *time.Time) {
	if date == nil {
		b.Date = time.Time{}
	} else {
		b.Date = *date
	}
}
func (b *BookLoanBuilder) SetDateAsString(date string) {
	t, err := time.Parse(time.DateOnly, date)
	if err != nil {
		panic("SetDateAsString() invalid date string")
	}
	b.Date = t
}
func (b *BookLoanBuilder) Build() *storage.BookLoan {
	bl := storage.NewBookLoan()
	bl.Title = b.Title
	bl.Author = b.Author
	bl.Genre = b.Genre
	bl.Ratting = b.Ratting
	bl.ID = b.id
	if b.IsOnLoan {
		bl.Loan.Name = b.Borrower
		bl.Loan.Date = b.Date
	} else {
		bl.UnsetLoan()
	}
	return bl
}






