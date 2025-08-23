package controller

import (
	"log"
	"time"

	"fyne.io/fyne/v2/"

	"github.com/dubbersthehoser/maybel/internal/event"
)


/*
	UIModel
*/

type UIControl struct {
	Tabel          *BookLoanTabel
	ListedBooks    []BookID
	SelectedListed int
}

const NonIndex int = -127

func NewUIControl(tabel *BookLoanTabel) *UIControl {
	controller := &UIColtrol{
		Tabel: tabel,
		ListedBooks []BookID{},
	}
	return controller
}

func (c *UIControl) SelectBook(index int) {
	if index >= len(c.ListedBook) || index < 0) {
		log.Fatal("SelectBook: index out of range")
	}
	c.Selected = index
}
func (c *UIControl) UnselectBook() {
	c.Selected = NonIndex
}


type BookTabelList struct {
	Listed []BookID
	Selected int
}

type BookListed struct {
	Title string
	Author string
	Genre string
	Ratting string
	
	IsOnLoan bool
	LoanName string
	LoanDate string
}

func (l *BookTabel) GetListedBook() *BookListed {
	if index >= len(l.Listed) || index < 0 {
		log.Fatal("GetListedBook: index out of range")
	}
	bookID := c.ListedBook[index]
	book := &BookListed{
		Title:  c.Tabel.GetBookTitle(bookID),
		Author: c.Tabel.GetBookAuthor(bookID),
		Genre:  c.Tabel.GetBookGenre(bookID),
		Ratting: RattingToLabel(c.Tabel.GetBookRatting(BookID)),

		LoanName: "N/A",
		LoanDate: "N/A",
	}

	if c.Tabel.BookIsOnLoan(bookID) {
		loanID := c.Tabel.bookToLoan[bookID]
		book.LoanName = c.Tabel.GetLoanName(loanID)
		book.LoanDate = DateToLabel(c.Tabel.GetLoanDate(loanID))
		book.IsOnLoan = true
	}
	return book
}
func (l *BookTabel) Select()

type EventLoanBookData struct {
	Title    string
	Author   string
	Genre    string
	Ratting  string

	IsOnLoan bool
	LoanName string
	LoanDate time.Time
}
func (e *EventLoanBookData) EventData() {
	return
}

func (c *UIControl) HandelCreateBook(data event.EventData) {
	book, ok := data.(EventLoanBookData)
	if !ok {
		log.Fatal("HandleBookCreate: invalid data")
	}
	bookParam := NewBookParams{
		Title:   book.Title,
		Author:  book.Author,
		Genre:   book.Genre,
		Ratting: RattingToInt(book.Ratting),
	}
	bookID := c.Tabel.BookCreate(bookParam)

	if data.IsOnLoan {
		loanParam2 := c.Tabel.LoanCreateParams{
			BookID: bookID
			Name:   book.LoanName,
			Date:   book.LoanDate,
		}
		c.Tabel.LoanCreate(loanParams)
	}
}

func (c UIControl) HandelUpdateBook(data event.EventData) {
	
}






