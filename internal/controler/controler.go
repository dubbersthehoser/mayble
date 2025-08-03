package controler

import (
	"fmt"
	"time"
	"log"
	"errors"
	"strings"
	"slices"

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
	bookList      []BookVM
	BookView      []*BookVM
	SelectedIndex int
	model         *model.Model
	orderBy       string
	ordering      string
	onError   func(error)
}
func NewVM() *VM {
	vm := &VM{
		bookList: []BookVM{},
		BookView: []*BookVM{},
		ordering: "ASC",
	}
	return vm
}
func (v *VM) NewBook() *BookVM {
	b := NewBook()
	b.id = noID
	return b
}
func (v *VM) GetBook(index int) *BookVM {
	if len(v.BookView) > index {
		log.Fatal("VM.GetBook(): index out of range")
	}
	return v.BookView[index]
}
func (v *VM) ViewListSize() int {
	return len(v.BookView)
}
func (v *VM) UniqueAuthors() []string {
	return []string{"AUTHOR_1", "AUTHOR_2"}
}

func (v *VM) UniqueGenres() []string {
	return []string{"GENRE_1", "GENRE_2"}
}

func (v *VM) AddBook(b *BookVM) {
	if v.BookView == nil {
		log.Fatal("Fatal: VM.BookView is nil")
	}
	if v.bookList == nil {
		log.Fatal("Fatal: VM.bookList is nil")
	}
	b.id = len(v.bookList)
	v.bookList = append(v.bookList, *b)
}

func (v *VM) SortBookView() {
	
	titleCompare := func(a *BookVM, b *BookVM) int {
		return strings.Compare(a.book.Title, b.book.Title)
	}
	authorCompare := func(a *BookVM, b *BookVM) int {
		return strings.Compare(a.book.Author, b.book.Author)
	}
	genreCompare := func(a *BookVM, b *BookVM) int {
		return strings.Compare(a.book.Genre, b.book.Genre)
	}
	rattingCompare := func(a *BookVM, b *BookVM) int {
		if a.book.Ratting > b.book.Ratting {
			return 1
		} else if a.book.Ratting < b.book.Ratting {
			return -1
		} else {
			return 0
		}
	}

	applyOrdering := func(fn func(*BookVM, *BookVM) int) func(*BookVM, *BookVM) int {
		if v.ordering == "DESC" {
			return func(a *BookVM, b *BookVM) int {
				return fn(a, b) * -1
				
			}
		} else {
			return fn
		}
	}

	switch v.orderBy {
	case "titles":
		slices.SortFunc(v.BookView, applyOrdering(titleCompare))
	case "authors":
		slices.SortFunc(v.BookView, applyOrdering(authorCompare))
	case "genres":
		slices.SortFunc(v.BookView, applyOrdering(genreCompare))
	case "rattings":
		slices.SortFunc(v.BookView, applyOrdering(rattingCompare))
	default:
		fmt.Printf("%s: ordering not implemented\n", v.orderBy)
	}
}

func (v *VM) UpdateBookView() {
	viewList := []*BookVM{}
	for _, book := range v.bookList {
		if book.flag != DeleteFlag {
			viewList = append(viewList, &book)
		}
	}
	v.BookView = viewList
}

func (v *VM) UpdateAndSortBookView() {
	v.UpdateBookView()
	v.SortBookView()
}

func (v *VM) RemoveBook(b *BookVM) {
	v.bookList = append(v.bookList[:b.id],  v.bookList[b.id+1:]...)
}
func (v *VM) ChangeOrderBy(label string) {
	fmt.Printf("Order by %s\n", label)
	v.orderBy = label
}
func (v *VM) SetOrderingASC() {
	v.ordering = "ASC"
	fmt.Printf("Ordering is ASC\n")
}
func (v *VM) SetOrderingDESC() {
	v.ordering = "DESC"
	fmt.Printf("Ordering is DESC\n")
}
func (v *VM) SetSelectedBook(index int) {
	v.SelectedIndex = index
	fmt.Printf("Book of index '%d' was selected\n",index)
}
func (v *VM) SelectedBook() *BookVM {
	return v.BookView[v.SelectedIndex]
}
func (v *VM) SetOnError(callback func(error)) {
	v.onError = callback
}




/*
	Loaning Entry
*/
type LoanVM struct {
	loan model.Loan
}
func NewLoan() *LoanVM {
	return &LoanVM{loan: model.Loan{}}
}
func (l *LoanVM) SetName(s string) {
	l.loan.Name = s
}
func (l *LoanVM) Name() string {
	if l.loan.Name == "" {
		return "N/A"
	}
	return l.loan.Name
}
func (l *LoanVM) SetDate(t time.Time) {
	l.loan.Date = t
}
func (l *LoanVM) Date() time.Time {
	return l.loan.Date
}
func (l *LoanVM) DateString() string {
	if *new(time.Time) == l.Date() {
		return "N/A"
	}
	t := l.Date()
	fmtdate := fmt.Sprintf("%d/%d/%d", t.Day(), t.Month(), t.Year())
	return fmtdate
}

type BookFlag int
const (
	NothingFlag BookFlag = iota
	DeleteFlag
	UpdateFlag
	NewFlag
)

/*
	Book Entry
*/
type BookVM struct {
	id   int
	flag BookFlag
	book model.Book
	loan *LoanVM
}

func NewBook() *BookVM {
	b := &BookVM{
		id: noID,
		loan: NewLoan(),
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

func (b *BookVM) SetTitle(s string) {
	b.setUpdateFlag()
	b.book.Title = s
}
func (b *BookVM) Title() string {
	return b.book.Title
}

func (b *BookVM) SetAuthor(s string) {
	b.setUpdateFlag()
	b.book.Author = s
}
func (b *BookVM) Author() string {
	return b.book.Author
}

func (b *BookVM) SetRatting(s string) {
	i := RattingToInt(s)
	b.setUpdateFlag()
	b.book.Ratting = i
}
func (b *BookVM) Ratting() string {
	return RattingToString(b.book.Ratting)
}

func (b *BookVM) SetGenre(s string) {
	b.setUpdateFlag()
	b.book.Genre = s
}
func (b *BookVM) Genre() string {
	if b.book.Genre == "" {
		return "N/A"
	}
	return b.book.Genre
}

func (b *BookVM) SetLoan(l *LoanVM) {
	b.setUpdateFlag()
	b.loan = l
}
func (b *BookVM) UnsetLoan() {
	b.setUpdateFlag()
	b.loan = NewLoan()
}
func (b *BookVM) Loan() *LoanVM {
	if b.loan == nil {
		log.Fatal("BookVM.Loan(): loan is nil")
	}
	return b.loan
}
func (b *BookVM) LoanName() string {
	return b.loan.Name()
}
func (b *BookVM) LoanDate() string {
	return b.loan.DateString()
}



/*
	Data Validation
*/

func (b *BookVM) TitleValidator(s string) error {
	if len(s) == 0 {
		return errors.New("Must have a Title")
	}
	return nil
}
func (b *BookVM) AuthorValidator(s string) error {
	if len(s) == 0 {
		return errors.New("Must have an Author")
	}
	return nil
}

func (l *LoanVM) NameValidator(s string) error {
	if len(s) == 0 {
		return errors.New("Must have an name")
	}
	return nil
}


