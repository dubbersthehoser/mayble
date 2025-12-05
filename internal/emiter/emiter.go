package emiter

import (
	"errors"
	"sync"
)


type Emiter struct {
	event map[string][]func(any)
	mu sync.RWMutex
}

func NewEmiter() *Emiter {
	return &Emiter{
		event: make(map[string][]func(any), 0),
	}
}

func (e *Emiter) OnEvent(key string, handler func(any)) {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, ok := e.event[key]
	if !ok {
		e.event[key] = make([]func(any), 0)
	}
	e.event[key] = append(e.event[key], handler)
}

func (e *Emiter) Emit(key string, data any) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	handlers, ok := e.event[key]
	if !ok {
		return errors.New("emiter: key not found")
	}
	for _, handle := range handlers {
		handle(data)
	}
	return nil
}


type Event struct {
	Name  string
	Data interface{}
}

type Listener struct {
	Handler func(*Event)
	id      int
}

type Broker struct {
	idCount int
	items map[string][]*Listener
	mu sync.RWMutex
}


func (b *Broker) Subscribe(l *Listener, events ...string) (int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.items == nil {
		b.items = make(map[string][]*Listener)
	}
	b.idCount++
	
	for _, e := range events {
		_, ok := b.items[e]
		if !ok {
			b.items[e] = make([]*Listener, 0)
		}
		l.id = b.idCount
		b.items[e] = append(b.items[e], l)
	}
	return l.id
}

func (b *Broker) Unsubscribe(id int, events ...string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, e := range events {
		listeners, ok := b.items[e]
		if !ok {
			continue
		}
		for i, l := range listeners {
			if l.id == id {
				listeners = append(listeners[:i], listeners[i+1:]...)
				b.items[e] = listeners
				break
			}
		}
	}
	return nil
}

func (b *Broker) Notify(e Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	listeners, ok := b.items[e.Name]
	if !ok {
		return errors.New("event name not found: " + e.Name)
	}
	for _, l := range listeners {
		l.Handler(&e)
	}
	return nil
}




