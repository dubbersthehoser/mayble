package controller

import (
	"time"
	"fmt"

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
		fmt.Printf("%#v\n", bookLoan)
		return be.core.CreateBookLoan(bookLoan)
	case Deleting:
		return be.core.DeleteBookLoan(bookLoan)
	default:
		return fmt.Errorf("submit type not found: %s", builder.Type)
	}
	return nil
}

func NewBookLoanBuilder(editType EditType) *BookLoanBuilder {
	return &BookLoanBuilder{
		Type: editType,
	}
} 

type BookLoanBuilder struct {
	Type     EditType
	id       int64
	title    string
	author   string
	genre    string
	ratting  int 
	isOnLoan bool
	borrower string
	date     time.Time
}

// TODO add validate function that returns an error

func (b *BookLoanBuilder) SetIsOnLoan(onLoan bool) {
	b.isOnLoan = onLoan
}
func (b *BookLoanBuilder) SetTitle(title string) {
	b.title = title
}
func (b *BookLoanBuilder) SetAuthor(author string) {
	b.author = author
}
func (b *BookLoanBuilder) SetGenre(genre string) {
	b.genre = genre
}
func (b *BookLoanBuilder) SetRatting(ratting int){
	b.ratting = ratting
}
// SetRattingAsString sets ratting from string value returned by RattingToStirng
func (b *BookLoanBuilder) SetRattingAsString(ratting string){
	b.ratting = RattingToInt(ratting)
}
func (b *BookLoanBuilder) SetBorrower(name string) {
	b.borrower = name
}
func (b *BookLoanBuilder) SetDate(date *time.Time) {
	b.date = *date
}
func (b *BookLoanBuilder) SetDateAsString(date string) {
	t, err := time.Parse(time.DateOnly, date)
	if err != nil {
		panic("SetDateAsString() invalid date string")
	}
	b.date = t
}
func (b *BookLoanBuilder) Build() *storage.BookLoan {
	bl := storage.NewBookLoan()
	bl.Title = b.title
	bl.Author = b.author
	bl.Genre = b.genre
	bl.Ratting = b.ratting
	bl.ID = b.id
	if b.isOnLoan {
		bl.Loan.Name = b.borrower
		bl.Loan.Date = b.date
	} else {
		bl.UnsetLoan()
	}
	return bl
}






