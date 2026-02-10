package viewmodel

import (
	"time"

	"fyne.io/fyne/v2/data/binding"
)

const (
	BodyData int = iota
	BodyForm
	BodyMenu
)
type MainUI struct {

	OpenedBody binding.Int

	Error      binding.String
	Success    binding.String
	Info       binding.String

}
func NewMainUI() *MainUI {
	mu := &MainUI{
		OpenedBody:  binding.NewInt(),

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


