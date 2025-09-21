package controller

import (
	"github.com/dubbersthehoser/mayble/internal/storage"
)

type Selector {
	controller *Controller
	bookLoan *storage.BookLoan
}
func NewSelector(c *Controller) *Selecter {
	var s Selector
	s.controller = c
	bookLoan = nil

}

func (s *Selector) SetBookLoan(bookLoan *storage.BookLoan) {
	if bookLoan != nil {
		tmp := *bookLoan
		bookLoan = &tmp
	}
	s.bookLoan = bookLoan
}

func (s *Selector) BookLoan() *storage.BookLoan {
	return s.bookLoan
}

func (s *Selector) NoBookLoan() bool {
	return s.bookLoan == nil
}
