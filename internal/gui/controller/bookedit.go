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

type BookEditor struct {
	controller Controller
	undo builderStack
	redo builderStack
}
func NewBookEditor(c *Controller) *BookEditor {
	return &BookEditor{
		controller: c,
	}
}

func (be *BookEditor) Submit(builder *BookLoanBuilder) {
	switch builder.Type {
	case Updating:
		
	case Creating:
		
	}
}

func NewBookLoanBuilder(editType EditType) *BookLoanBuilder {
	return &BookLoanbuilder{
		Type: editType,
	}
} 


type BookLoanBuilder struct {
	Type     EditType
	title    string
	author   string
	genre    string
	ratting  int 
	isOnLoan bool
	loanName string
	loanDate time.Time
}
func (b *BookLoanBuilder) SetIsLoan(onLoan bool) {
	b.isOnLoan = onLoan
}
func (b *BookLoanBuilder) SetTitle(title string) {
	b.title = title
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
func (b *BookLoanBuilder) SetLoanName(name string) {
	b.loanName = name
}
func (b *BookLoanBuilder) SetLoanDate(date *time.Time) {
	b.loanDate = *date
}
func (b *BookLoanBuilder) Build() *storage.BookLoan {
	bl := storage.NewBookLoan()
	bl.Title = b.title
	bl.Author = b.author
	bl.Genre = b.genre
	bl.Ratting = RattingToInt(b.Ratting)
	if b.isOnLoan {
		bl.Loan.Name = b.loanName
		bl.Loan.Date = b.loanDate
	} else {
		bl.SetUnloan()
	}
	return &bl
}

type builderStack struct {
	stack []BookLoanBuilder
}

func (b *builderStack) Push(bl *BookLoanBuilder) {
	if bl == nil {
		panic("stack can't push a nil value")
	}
	b.stack = append(b.stack, *bl)
}
func (b *builderStack) Pop() *BookLoanBuilder {
	if len(b.stack) == 0 {
		return nil
	}
	bl := b.Peek()
	b.stack = b.stack[:len(b.stack)-1]
	return &bl
}
func (b *builderStack) Peek() *BookLoanBuilder {
	bl := b.stack[len(b.stack)-1]
	return &bl
}










