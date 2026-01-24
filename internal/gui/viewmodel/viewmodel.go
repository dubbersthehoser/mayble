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






type TableData struct {
	Header []string
	Items  [][]string
	Sizes  []int
}
func newTableData() *TableData {
	return &TableData{
		Header: make([]string, 0),
		Items:  make([][]string, 0),
		Sizes:  make([]int, 0),
	}
}

func (t *TableData) SetHeader(h []string) {
	t.Header = h
	t.Items = make([][]string, 0)
	t.Sizes = make([]int, len(h))
}


func (t *TableData) Length() int {
	return len(t.Items)
}

func (t *TableData) Append(values... string) error {
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


func translateExcludedIndex(idx int, header, exclude []string) int {
	indexes  := []int{}
	for i, h := range header {
		if tmp := slices.Index(exclude, h); tmp == -1 {
			indexes = append(indexes, i)
		}
	}
	return indexes[idx]
}






type ColumnExcluder struct {
	selected  []string
	listeners []binding.DataListener
}
func NewColumnExcluder() *ColumnExcluder {
	cs := &ColumnExcluder{}
	return cs
}
func (cs *ColumnExcluder) SetHeader(header []string) {
	cs.selected = header
	cs.notify()
}
func (cs *ColumnExcluder) GetHeader() []string {
	return cs.selected
}

func (cs *ColumnExcluder) notify() {
	for _, l := range cs.listeners {
		l.DataChanged()
	}
}

func (cs *ColumnExcluder) AddListener(l binding.DataListener) {
	cs.listeners = append(cs.listeners, l)
}
func (cs *ColumnExcluder) RemoveListener(l binding.DataListener) {
	index := -1
	for i, listener := range cs.listeners {
		if listener == l {
			index = i
		}
	}
	if index == -1 {
		return
	}
	cs.listeners = append(cs.listeners[:index], cs.listeners[index-1:]...)
}








type Table struct {
	
	Data *TableData
	Excluder *ColumnExcluder

	ShowColumns binding.BoolList

	OrderField binding.String
	OrderASC   binding.Bool

	listeners []binding.DataListener
}

func NewTable() *Table {
	t := &Table{
		Data: newTableData(),
		Excluder: NewColumnExcluder(),
		ShowColumns: binding.NewBoolList(),

		OrderField: binding.NewString(),
		OrderASC: binding.NewBool(),

		listeners: make([]binding.DataListener, 0),
	}
	t.Data.SetHeader([]string{"Title", "Author", "Genre"})
	t.Data.Append("Example Title", "Author", "Example Genre")
	t.Data.Append("Example Title", "Author", "Example Genre")
	t.Data.Append("Example Title", "Author", "Example Genre")
	_ = t.OrderField.Set(t.Data.Header[0])


	t.Excluder.AddListener(binding.NewDataListener(func(){
		if len(t.Excluder.selected) != len(t.Data.Header) {
			idx := translateExcludedIndex(0, t.Data.Header, t.Excluder.selected)
			t.OrderField.Set(t.Data.Header[idx])
		}
		t.notify()
	}))


	return t
}

func (t *Table) Size() (int, int) {
	return len(t.Data.Items), len(t.Data.Header) - len(t.Excluder.selected)
}

func (t *Table) GetValue(row, column int) string {
	col := translateExcludedIndex(column, t.Data.Header, t.Excluder.selected)
	return t.Data.Items[row][col]
}
func (t *Table) GetHeader(column int) string {
	idx := translateExcludedIndex(column, t.Data.Header, t.Excluder.selected)
	return t.Data.Header[idx]
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
