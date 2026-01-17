package viewmodel

import (
	"errors"

	"fyne.io/fyne/v2/data/binding"
)

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


type MainUI struct {
	
}
func NewMainUI() *MainUI {
	mu := &MainUI{
		
	}
	return mu
}




type BookEntry struct {
	Title  binding.String
	Author binding.String
	Genre  binding.String
}

var _ binding.DataList = &Table{}

type Table struct {
	HeaderLables []string
	entries []binding.DataItem
}

func (t *Table) AddListener(binding.DataListener) {

}
func (t *Table) RemoveListener(binding.DataListener) {

}

func (t *Table) GetItem(index int) (binding.DataItem, error) {
	if len(t.entries) <= index || 0 > index {
		return nil, errors.New("index out of range")
	}
	return t.entries[index], nil
}

func (t *Table) Length() int {
	return len(t.entries)
}







type ReadEntry struct {
	Rating   binding.String
	Competion binding.String
}

type LoanEntry struct {
	Name binding.String
	Date binding.String
}





