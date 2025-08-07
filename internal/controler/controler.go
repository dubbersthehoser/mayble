package controler

import (
	"time"
	"log"

	"github.com/dubbersthehoser/mayble/internal/model"
)


const noID int = -127
type VM struct {
	bookList      []BookVM
	model         *model.Model
	onError       func(error)
}
func NewVM() *VM {
	vm := &VM{}
	return vm
}
func (v *VM) UniqueAuthors() []string {
	return []string{"AUTHOR_1", "AUTHOR_2"}
}

func (v *VM) UniqueGenres() []string {
	return []string{"GENRE_1", "GENRE_2"}
}

func (v *VM) AddBook(b *BookVM) {
	b.id = len(v.bookList)
	v.bookList = append(v.bookList, *b)
}

func (v *VM) RemoveBook(b *BookVM) {
	v.bookList = append(v.bookList[:b.id],  v.bookList[b.id+1:]...)
}

func (v *VM) SetOnError(callback func(error)) {
	v.onError = callback
}

/*
	Book Entry
*/
type BookFlag int
const (
	NothingFlag BookFlag = iota
	DeleteFlag
	UpdateFlag
	NewFlag
)

type BookVM struct {
	id   int
	flag BookFlag
	book model.Book
	loan model.Loan

}

func NewBook() *BookVM {
	b := &BookVM{
		id: noID,
		flag: NewFlag,
	}
	return b
}

func (b *BookVM) setUpdateFlag() {
	if b.flag == DeleteFlag {
		log.Fatal("updating a deleted vmbook!")
	}
	if b.flag != NewFlag {
		b.flag = UpdateFlag
	}
}

func (b *BookVM) Title() string {
	return b.book.Title
}

func (b *BookVM) Author() string {
	return b.book.Author
}

func (b *BookVM) Ratting() int {
	return b.book.Ratting
}

func (b *BookVM) Genre() string {
	return b.book.Genre
}

func (b *BookVM) LoanName() string {
	return b.loan.Name
}

func (b *BookVM) LoanDate() *time.Time {
	date := b.loan.Date
	return &date
}

func (b *BookVM) SetTitle(s string) {
	b.book.Title = s
}

func (b *BookVM) SetAuthor(s string) {
	b.book.Author = s
} 

func (b *BookVM) SetGenre(s string) {
	b.book.Genre = s
}

func (b *BookVM) SetRatting(r int) {
	b.book.Ratting = r
}

func (b *BookVM) SetLoanName(s string) {
	b.loan.Name = s
}

func (b *BookVM) SetLoanDate(t *time.Time) {
	b.loan.Date = *t
}

