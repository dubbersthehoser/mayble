package viewmodel

import (
	"github.com/dubbersthehoser/mayble/internal/event"
)

const (
	BodyNoData int = iota
	BodyTable
	BodyBookEdit
	BodyBookCreate
	BodyManual
)

type Body struct {
	prev  int
	value int
	l     []func()

	hideHandlers map[int]func()
	showHandlers map[int]func()
}

func newBody(eb *event.EventBus) *Body {
	b := &Body{}
	setupBodyToEvents(b, eb)
	return b
}

func (b *Body) RegisterHandlers(body int, hide, show func()) {
	if b.hideHandlers == nil {
		b.hideHandlers = make(map[int]func())
	}

	if b.showHandlers == nil {
		b.showHandlers = make(map[int]func())
	}

	b.hideHandlers[body] = hide
	b.showHandlers[body] = show
}

func (b *Body) Value() int {
	return b.value
}

func (b *Body) Set(v int) {
	for _, h := range b.hideHandlers {
		h()
	}
	if b.showHandlers != nil {
		b.showHandlers[v]()
	}

	if BodyManual != b.value {
		b.prev = b.value
	}

	b.value = v
	b.notify()
}

func (b *Body) Back() {
	b.Set(b.prev)
}

func (b *Body) AddListener(fn func()) {
	if b.l == nil {
		b.l = make([]func(), 0)
	}
	b.l = append(b.l, fn)
}

func (b *Body) notify() {
	for _, fn := range b.l {
		fn()
	}
}

func setupBodyToEvents(b *Body, eb *event.EventBus) {
	eb.Subscribe(event.UpdatedBookEntry{}, func(v event.Event) {
		e := v.(event.UpdatedBookEntry)
		// when updating an entry go back to table when completed successfully.
		if !e.Failed {
			b.Set(BodyTable)
		}
	})
	eb.Subscribe(event.OpenedDatabase{}, func(v event.Event) {
		e := v.(event.OpenedDatabase)
		if e.Failed && b.Value() != BodyTable {
			b.Set(BodyNoData)
		} else {
			b.Set(BodyTable)
		}
	})
	eb.Subscribe(event.CreatedDatabase{}, func(v event.Event) {
		e := v.(event.CreatedDatabase)
		if e.Failed && b.Value() != BodyTable {
			b.Set(BodyNoData)
		} else {
			b.Set(BodyTable)
		}
	})
}
