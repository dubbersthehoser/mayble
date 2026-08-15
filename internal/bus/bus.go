package bus

import (
	"log"
)

type Event struct {
	Name string
	Data any
}

type Handler struct {
	Handle func(e *Event)
}

type Bus struct {
	bus map[string][]Handler
}

func (b *Bus) Register(Name string, h Handler) {
	if b.bus == nil {
		b.bus = make(map[string][]Handler)
	}

	s, ok := b.bus[Name]
	if !ok {
		s = make([]Handler, 0)
	}
	s = append(s, h)
	b.bus[Name] = s
}

func (b *Bus) Emit(e Event) {
	if b == nil {
		log.Println("called emit without bus handlers")
		return
	}

	s, ok := b.bus[e.Name]
	if !ok {
		log.Printf("event '%s' not found\n", e.Name)
		return
	}

	for _, h := range s {
		h.Handle(&e)
	}
}
