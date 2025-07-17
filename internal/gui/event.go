package gui

import (
	"log"
	"fmt"
)

func (u *UIState) EventNewBook(data any) {
	b := &Book{}
	d := u.NewBookDialog(b)
	d.Show()

	fmt.Printf("%#v\n", *b)
	if (Book{}) != *b {
		fmt.Printf("%#v\n", *b)
		return
	}
}

func (u *UIState) EventUpdateBook(data any) {
	b, ok := data.(*Book)
	if !ok {
		log.Fatal("No book was given to EventUpdateBook()...")
	}

	d := u.NewBookDialog(b)
	d.Show()
}

func (u *UIState) EventNewOnLoan(data any) {
	b, ok := data.(*Book)
	if !ok {
		log.Fatal("No book was given to EventNewOnLoan()...")
	}
	l := &Loan{}
	d := u.NewOnLoanDialog(l)
	d.Show()
	if (Loan{}) != *l {
		b.IsOnLoan = true
		b.OnLoan = *l
		fmt.Printf("%#v\n", *b)
	}
}
