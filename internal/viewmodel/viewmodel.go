package viewmodel

import (
	"time"

	"fyne.io/fyne/v2/data/binding"

	repo "github.com/dubbersthehoser/mayble/internal/repository"
)


type mockBookSearcher struct {}

var _ repo.BookQuerier = &mockBookSearcher{}

func (m *mockBookSearcher) BookQuery(q *repo.Query) ([]repo.BookEntry, error) {
	es := []repo.BookEntry{
		{
			Variant: repo.Book,
			Title: "Example Title",
			Author: "Example Author",
			Genre: "Example Genre",
		},
		{
			Variant: repo.Book,
			Title: "Example Title",
			Author: "Example Author",
			Genre: "Example Genre",
		},
		{
			Variant: repo.Book,
			Title: "Example Title",
			Author: "Example Author",
			Genre: "Example Genre",
		},
		{
			Variant: repo.Book,
			Title: "Example Title",
			Author: "Example Author",
			Genre: "Example Genre",
		},
	}
	return es, nil
}


const (
	BodyData int = iota
	BodyForm
	BodyMenu
)

type MainUI struct {	

	OpenedBody binding.Int

	Repo       repo.BookQuerier

	Error      binding.String
	Success    binding.String
	Info       binding.String

}

func NewMainUI() *MainUI {
	mu := &MainUI{
		
		OpenedBody:  binding.NewInt(),

		Repo: &mockBookSearcher{},

		Error: binding.NewString(),
		Success: binding.NewString(),
		Info: binding.NewString(),
	}
	return mu
}


type BookVM struct {
	id int64
	Title binding.String
	Author binding.String
	Genre binding.String
}
func NewBookVM(id int64, title, author, genre string) *BookVM {
	vm := &BookVM{
		id: id,
		Title: binding.NewString(),
		Author: binding.NewString(),
		Genre: binding.NewString(),
	}
	_ = vm.Title.Set(title)
	_ = vm.Author.Set(author)
	_ = vm.Genre.Set(genre)
	return vm
}


const dateFormat = "02/01/2006"

func formatDate(t *time.Time) string {
	return t.Format(dateFormat)
}

func parseDate(t string) (*time.Time, error) {
	ret, err := time.Parse(t, dateFormat)
	return &ret, err
}


func formatRating(r int) string {
	switch r {
	case 1:
		return "⭐"
	case 2:
		return "⭐⭐"
	case 3:
		return "⭐⭐⭐"
	case 4:
		return "⭐⭐⭐⭐"
	case 5:
		return "⭐⭐⭐⭐⭐"
	default:
		return "N/A"
	}
}

func RatingsStrings() []string {
	s := 6
	r := make([]string, s)
	for i := range s {
		r[i] = formatRating(i+1)
	}
	return r
}
