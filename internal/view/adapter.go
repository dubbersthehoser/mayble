package gui

import (
	"time"
	"fmt"
	"strings"
	"errors"
	"log"

	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/controler"
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

func DateToString(date *time.Time) string {
	if date == nil {
		return "N/A"
	}
	if *date == *new(time.Time) {
		return "N/A"
	}
	return fmt.Sprintf("%d/%d/%d", date.Day(), date.Month(), date.Year())
}

type SortOrder int
const (
	ASC  SortOrder = iota
	DESC
)

func GetOrderingComp(orderBy string, ordering SortOrder) func(*ListedBook, *ListedBook) int {

	titleCompare := func(a *ListedBook, b *ListedBook) int {
		return strings.Compare(a.Data.Title(), b.Data.Title())
	}
	authorCompare := func(a *ListedBook, b *ListedBook) int {
		return strings.Compare(a.Data.Author(), b.Data.Author())
	}
	genreCompare := func(a *ListedBook, b *ListedBook) int {
		return strings.Compare(a.Data.Genre(), b.Data.Genre())
	}
	rattingCompare := func(a *ListedBook, b *ListedBook) int {
		if a.Data.Ratting() > b.Data.Ratting() {
			return 1
		} else if a.Data.Ratting() < b.Data.Ratting() {
			return -1
		} else {
			return 0
		}
	}

	applyOrdering := func(fn func(*ListedBook, *ListedBook) int) func(*ListedBook, *ListedBook) int {
		if ordering == DESC {
			return func(a *ListedBook, b *ListedBook) int {
				return fn(a, b) * -1
			}
		}
		return fn
	}

	switch orderBy {
	case "title", "Title":
		return applyOrdering(titleCompare)
	case "author", "Author":
		return applyOrdering(authorCompare)
	case "genre", "Genre":
		return applyOrdering(genreCompare)
	case "ratting", "Rattting":
		return applyOrdering(rattingCompare)
	default:
		fmt.Printf("%s: ordering not implemented\n", orderBy)
		log.Fatal("ERROR")
	}
	return nil
}

type FyneVM struct {
	VM           *controler.VM
	SelectedBook *ListedBook
	ListedBooks  []ListedBook
	SortOrder    SortOrder
	OrderComp    func(a *ListedBook, b *ListedBook) int
}

func NewFyneVM() *FyneVM {
	vm := &FyneVM{
		VM: controler.NewVM(),
	}
	return vm
}

func (vm *FyneVM) ListedCount() int {
	return len(vm.ListedBooks)
}
func (vm *FyneVM) SetSelectedBook(index int) {
	vm.SelectedBook = vm.GetBook(index)
}

func (vm *FyneVM) UniqueAuthors() []string {
	return vm.VM.UniqueAuthors()
}
func (vm *FyneVM) UniqueGenres() []string {
	return vm.VM.UniqueGenres()
}

func (vm *FyneVM) GetBook(index int) *ListedBook {
	if index >= vm.ListedCount() {
		log.Fatal("FyneVM.GetBook(): index out of range")
	}
	return &vm.ListedBooks[index]
}

func (vm *FyneVM) NewBook() *ListedBook {
	book := controler.NewBook()
	listed := NewListedBook(book)
	return listed
}

func (vm *FyneVM) AddBook(listed *ListedBook) {
	vm.VM.AddBook(listed.Data)
	vm.ListedBooks = append(vm.ListedBooks, *listed)
}

func (vm *FyneVM) SetSorting(orderBy string, ordering SortOrder) {
	vm.OrderComp = GetOrderingComp(orderBy, ordering)
}
func (vm *FyneVM) Sort() {
	println("Sorting")
}


type ListedBook struct {
	Data   *controler.BookVM
	Label  *BookLabel
	Entry  *BookEntry
}

func NewListedBook(data *controler.BookVM) *ListedBook {
	listed := &ListedBook{
		Data: data,
		Label: NewBookLabel(),
		Entry: NewBookEntry(),
	}
	return listed
} 

func (l *ListedBook) EntryToData() {
	title, _ := l.Entry.Title.Get()
	author, _ := l.Entry.Author.Get()
	genre, _ := l.Entry.Genre.Get()
	ratting, _ := l.Entry.Ratting.Get()
	loanName, _ := l.Entry.LoanName.Get()
	loanDate := l.Entry.LoanDate

	fmt.Printf("%#v\n", loanDate)
	
	l.Data.SetTitle(title)
	l.Data.SetAuthor(author)
	l.Data.SetGenre(genre)
	l.Data.SetRatting(RattingToInt(ratting))
	l.Data.SetLoanName(loanName)
	l.Data.SetLoanDate(loanDate)
}

func (l *ListedBook) DataToLabel() {

	title := l.Data.Title()
	author := l.Data.Author()
	genre := l.Data.Genre()
	ratting := RattingToString(l.Data.Ratting())
	loanName := l.Data.LoanName()
	loanDate := DateToString(l.Data.LoanDate())

	if genre == "" {
		genre = "N/A"
	}
	if loanName == "" {
		loanName = "N/A"
	}

	l.Label.Title.Set(title)
	l.Label.Author.Set(author)
	l.Label.Genre.Set(genre)
	l.Label.Ratting.Set(ratting)
	l.Label.LoanName.Set(loanName)
	l.Label.LoanDate.Set(loanDate)
}

func (l *ListedBook) IsOnLoan() bool {
	if l.Data.LoanName() == "" {
		return false
	}
	return true
}

// Primary use is for the labels on the rows
type BookLabel struct {
	Title    binding.String
	Author   binding.String
	Genre    binding.String
	Ratting  binding.String
	LoanName binding.String
	LoanDate binding.String
}

func NewBookLabel() *BookLabel {
	b := &BookLabel{
		Title:    binding.NewString(),
		Author:   binding.NewString(),
		Genre:    binding.NewString(),
		Ratting:  binding.NewString(),
		LoanName: binding.NewString(),
		LoanDate: binding.NewString(),
	}
	return b
}

type BookEntry struct {
	Title    binding.String
	Author   binding.String
	Genre    binding.String
	Ratting  binding.String
	LoanName binding.String
	LoanDate *time.Time
}

func NewBookEntry() *BookEntry {
	b := &BookEntry{
		Title:    binding.NewString(),
		Author:   binding.NewString(),
		Genre:    binding.NewString(),
		Ratting:  binding.NewString(),
		LoanName: binding.NewString(),
		LoanDate: &time.Time{},
	}
	return b
}

func (e *BookEntry) UnsetLoan() {
	e.LoanName.Set("")
	e.LoanDate = &time.Time{}
}

/*
	Data Validation
*/
func (e *BookEntry) TitleValidator(s string) error {
	if len(s) == 0 {
		return errors.New("Must have a Title")
	}
	return nil
}

func (e *BookEntry) AuthorValidator(s string) error {
	if len(s) == 0 {
		return errors.New("Must have an Author")
	}
	return nil
}

func (e *BookEntry) LoanNameValidator(s string) error {
	if len(s) == 0 {
		return errors.New("Must have an name")
	}
	return nil
}








