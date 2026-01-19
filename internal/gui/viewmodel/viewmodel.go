package viewmodel

import (
	"errors"

	"fyne.io/fyne/v2/data/binding"
)


// API
type Book struct {}
type BookRead struct {}
type BookLoan struct {}

type Store[T any] interface {
	GetAllItems() ([]T, error)
	GetItemByID(int64) (T, error)
	CreateItem(*T) (int64, error)
	DeleteItem(*T) (error)
	UpdateItem(*T) (error)
}

type BookStore = Store[Book] 
type LoanStore = Store[BookLoan] 
type ReadStore = Store[BookRead]

type event struct {
	name string
	data any
}

type eventBus struct {
	bus map[string][]func(event)
}

func (b *eventBus) subscribe(eventName string, handler func(event)) {
	if b.bus == nil {
		b.bus = make(map[string][]func(event))
	}
	handlers, ok := b.bus[eventName]
	if !ok {
		handlers = make([]func(event), 0)
	}
	handlers = append(handlers, handler)
	b.bus[eventName] = handlers
}
func (b *eventBus) publish(e event) {
	handlers, ok := b.bus[e.name]
	if !ok {
		return
	}
	for _, handler := range handlers {
		handler(e)
	}
}



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

type BookLoanVM struct {
	bookID   int64
	Date     binding.String
	Borrower binding.String
}
type BookReadVM struct {
	bookID     int64
	Completion binding.String
	Rating     binding.String
}

type Table struct {
	Header    []string
	Items     []binding.Struct
	listeners []binding.DataListener
	
	OrderField binding.String
}
func NewTable() *Table {
	t := &Table{
		Items: make([]binding.Struct, 0),
		listeners: make([]binding.DataListener, 0),
		Header: []string{"Title", "Author", "Genre"},
	}
	t.Items = append(t.Items,
		binding.BindStruct(NewBookVM(0, "Example Title", "Example Author", "Example Genre")),
		binding.BindStruct(NewBookVM(1, "Example Title", "Example Author", "Example Genre")),
	)
	return t
}
func (t *Table) Length() int {
	return len(t.Items)
}
func (t *Table) GetItem(index int) (binding.DataItem, error) {
	return t.GetValue(index)
}
func (t *Table) GetValue(index int) (binding.Struct, error) {
	if len(t.Items) >= index || 0 > index {
		return nil, errors.New("index out of range")
	}
	return t.Items[index], nil
}
func (t *Table) AddListener(l binding.DataListener) {
	t.listeners = append(t.listeners, l)
}
func (t *Table) RemoveListener(binding.DataListener) {
	panic("not-implemented")
}

