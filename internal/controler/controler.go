package controler

import (
	"fmt"
	"time"
	"errors"

	"github.com/dubbersthehoser/mayble/internal/model"
)

// GetRattinStrings return ratting labels that map to ratting int.
func GetRattingStrings() []string {
	return []string{"TBR", "⭐", "⭐⭐", "⭐⭐⭐", "⭐⭐⭐⭐", "⭐⭐⭐⭐⭐"}
}

// RattingToString given ratting int return string for ui. If int is out of bounds will truncate with in range.
func RattingToString(i int) string {
	r := GetRattingStrings()
	if len(r) <= i {
		i = len(r) - 1
	} else if i < 0 {
		i = 0
	}
	return r[i]
}
func RattingToInt(s string) int {
	for i, l := range GetRattingStrings() {
		if l == s {
			return i
		}
	}
	return 0
}

const noID int = -127
type VM struct {
	BookList []BookVM
	model    *model.Model
}
func (v *VM) NewBook() *BookVM {
	b := &BookVM{id: noID}
	return b
}
func (v *VM) NewLoan() *LoanVM {
	l := &LoanVM{}
	return l
}
func (v *VM) UniqueAuthors() []string {
	return []string{"PLACEHOLDER_1", "PLACEHOLDER_2"}
}
func (v *VM) UniqueGenres() []string {
	return []string{"PLACEHOLDER_1", "PLACEHOLDER_2"}
}
func (v *VM) AddBook(b *BookVM) {
	fmt.Println("add book")
}
func (v *VM) ChangeOrderBy(label string) {
	fmt.Printf("Order by %s\n", label)
}
func (v *VM) SetOrderingASC() {
	fmt.Printf("Ordering is ASC\n")
}
func (v *VM) SetOrderingDESC() {
	fmt.Printf("Ordering is DESC\n")
}
func (v *VM) SetBookSelected(id int) {
	fmt.Printf("Book of id '%d' was selected\n", id)
}

type LoanVM struct {
	loan model.Loan
}
func NewLoan() *LoanVM {
	return &LoanVM{}
}
func (l *LoanVM) SetName(s string) {
	l.loan.Name = s
}
func (l *LoanVM) SetDate(t time.Time) {
	l.loan.Date = t
}
func (l *LoanVM) Date() time.Time {
	return l.loan.Date
}
func (b *LoanVM) NameValidator(s string) error {
	if len(s) == 0 {
		return errors.New("Must have an name")
	}
	return nil
}

type BookVM struct {
	id   int
	book model.Book
	loan *LoanVM
}


func (b *BookVM) SetTitle(s string) {
	b.book.Title = s
}
func (b *BookVM) Title() string {
	return b.book.Title
}

func (b *BookVM) SetAuthor(s string) {
	b.book.Author = s
}
func (b *BookVM) Author() string {
	return b.book.Author
}

func (b *BookVM) SetRatting(s string) {
	i := RattingToInt(s)
	b.book.Ratting = i
}
func (b *BookVM) Ratting() string {
	return RattingToString(b.book.Ratting)
}

func (b *BookVM) SetGenre(s string) {
	b.book.Genre = s
}
func (b *BookVM) Genre() string {
	return b.book.Genre
}

func (b *BookVM) SetLoan(l *LoanVM) {
	b.loan = l
}
func (b *BookVM) UnsetLoan() {
	b.loan = nil
}


func (b *BookVM) TitleValidator(s string) error {
	fmt.Printf("Title Validation: '%s'\n", s)
	if len(s) == 0 {
		return errors.New("Must have a Title")
	}
	return nil
}
func (b *BookVM) AuthorValidator(s string) error {
	fmt.Printf("Author Validation: '%s'\n", s)
	if len(s) == 0 {
		fmt.Println("FAILED")
		return errors.New("Must have an Author")
	}
	fmt.Println("SUCCSESS")
	return nil
}




