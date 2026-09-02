package event

import (
	"sync"
	"fmt"
)

type NameHandler func(v any) (string, error)

type Event any

type EventBus struct {
	mu sync.Mutex
	handlers map[string][]func(v Event)
}

func NewEventBus() *EventBus {
	eb := &EventBus{
		handlers: make(map[string][]func(v Event)),
	}
	return eb
}

func (eb *EventBus) Subscribe(e Event, h func(v Event)) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	name := fmt.Sprintf("%T", e)
	handlers, ok := eb.handlers[name]
	if !ok {
		handlers = make([]func(Event), 0)
	}

	handlers = append(handlers, h)
	eb.handlers[name] = handlers
	return nil
}

func (eb *EventBus) Notify(v Event) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	name := fmt.Sprintf("%T", v)
	handlers, ok := eb.handlers[name]
	if !ok {
		return fmt.Errorf("notify %s: event not found", name)
	}
	for _, h := range handlers {
		h(v)
	}
	return nil
}
