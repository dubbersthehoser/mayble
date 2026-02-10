package bus

import (
	"slices"
)

type Event struct {
	On string
	Data any
}

type Listener struct{
	Handler func(e Event)
	On      string
	id      int
}

type Bus struct {
	listeners map[string][]Listener
	lastID int
}

func (b *Bus) Publish(e Event) {
	lst, ok := b.listeners[e.On]
	if !ok {
		return
	}

	for i := range lst {
		lst[i].Handler(e)
	}
}

func (b *Bus) Subscribe(l Listener) int {
	id := b.lastID + 1
	b.lastID = id

	l.id = id

	lst, ok := b.listeners[l.On]

	if !ok {
		lst = make([]Listener, 0)
	}

	lst = append(lst, l)

	b.listeners[l.On] = lst

	return id
}

func (b *Bus) Unsubscribe(id int, off string) {

	_, ok := b.listeners[off]
	if !ok {
		return
	}
	
	b.listeners[off] = slices.DeleteFunc(b.listeners[off], func(l Listener) bool {
		return id == l.id
	})
}
