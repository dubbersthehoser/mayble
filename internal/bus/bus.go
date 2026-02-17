package bus

type Event struct {
	Name string
	Data any
}

type Handler struct {
	Name    string
	Handler func(e *Event)
	id      int
	next    int
}

func nop(e *Event) {}

type Bus struct {
	items []Handler
	live  map[string]int
	free  int
}

func (b *Bus) Subscribe(h Handler) int {
	if b.items == nil {
		b.items = make([]Handler, 1)
		b.items[0].Handler = nop
	}
	if b.live == nil {
		b.live = make(map[string]int)
	}

	if b.free == 0 {
		b.items = append(b.items, h)
		id := len(b.items) - 1
		h.id = id
	} else {
		id := b.free
		b.free = b.items[id].next
		h.id = id
	}

	b.items[h.id] = h
	first, ok := b.live[h.Name]
	if !ok {
		b.live[h.Name] = h.id
	} else {
		b.items[h.id].next = first
		b.live[h.Name] = h.id
	}
	return h.id
}

func (b *Bus) Unsubscribe(id int) {
	if b.items == nil || b.live == nil  {
		return
	}
	if id > len(b.items) || id <= 0 {
		return
	}
	h := b.items[id]
	first := b.live[h.Name]
	curr := first
	prev := 0
	for {
		if curr == id {
			next := b.items[curr].next
			if curr == first {
				b.live[h.Name] = next
			} else {
				b.items[prev].next = next
			}
			b.items[curr].next = b.items[b.free].next
			b.free = curr 
			return
		}
		prev = curr
		curr = b.items[curr].next
		if curr == 0 {
			break
		}
	}
}

func (b *Bus) Notify(e Event) {
	name := e.Name
	if b.items == nil {
		return
	}
	first, ok := b.live[name]
	if !ok {
		return
	}
	curr := first
	for {
		b.items[curr].Handler(&e)
		curr = b.items[curr].next
		if curr == 0 {
			break
		}
	}
}
