package controller

import (
	"time"
)

type BookForm struct {
	Title string
	
	onloan bool
	LoanName string
	LoanDate *time.Time
}

func NewBookForm() *BookForm {
	b := &BookForm{}
	return b
}

func (b *BookForm) IsOnLoan() bool {
	return b.onLoan
}

func (b *BookForm) LoanTo(loan *LoanData) {
	b.Loan.
}
func (b *BookForm) BookIs(book *BookData)

func (b *BookForm) Unloan() {
	b.onloan = false
}



