package viewmodel

import (
	"slices"
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
	Completion string
	Rating     string
}

type Table struct {
	Header    []string
	Sizes     []int
	Items     [][]string
	listeners []binding.DataListener
	
	OrderField binding.String
}

func NewTable() *Table {
	t := &Table{
		Items: make([][]string, 0),
		listeners: make([]binding.DataListener, 0),
		Header: []string{"Title", "Author", "Genre"},
	}
	t.Sizes = make([]int, len(t.Header))
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	t.Append("Example Title", "Example Author", "Example Genre")
	return t
}

func (t *Table) Append(values... string) error {
	if len(values) != len(t.Header) {
		return errors.New("length of values missmatch header length")
	}
	for i := range len(t.Header) {
		if t.Sizes[i] < len(values[i]) {
			t.Sizes[i] = len(values[i])
		}
	}
	t.Items = append(t.Items, values)
	return nil
}

func (t *Table) Length() int {
	return len(t.Items)
}
func (t *Table) GetValue(index int) ([]string, error) {
	if len(t.Items) >= index || 0 > index {
		return nil, errors.New("index out of range")
	}
	return t.Items[index], nil
}
func (t *Table) SetValue(index int, key string, value string) error {
	if len(t.Items) >= index || 0 > index {
		return errors.New("index out of range")
	}
	field := slices.Index(t.Header, key)
	if field == -1 {
		return errors.New("key not found")
	}
	t.Items[index][field] = value
	t.notify()
	return nil
}

func (t *Table) notify() {
	for _, listener := range t.listeners {
		listener.DataChanged()
	}
}

func (t *Table) AddListener(l binding.DataListener) {
	t.listeners = append(t.listeners, l)
}
func (t *Table) RemoveListener(l binding.DataListener) {
	index := -1
	for i, listener := range t.listeners {
		if listener == l {
			index = i
		}
	}
	if index == -1 {
		return
	}
	t.listeners = append(t.listeners[:index], t.listeners[index-1:]...)
}

