package bus

import (
	"testing"
)

func TestBus(t *testing.T) {

	callCount := 0
	call := func(e *Event) {
		callCount += 1
	}

	var bus *Bus = &Bus{}
	
	tests := []struct{
		input     Handler
		id        int
		callCount int
	}{
		{ input: Handler{ Handler: call, Name: "first",}, id: 1, callCount: 1},
		{ input: Handler{ Handler: call, Name: "second"}, id: 2, callCount: 3},
		{ input: Handler{ Handler: call, Name: "second"}, id: 3, callCount: 5},
		{ input: Handler{ Handler: call, Name: "third"},  id: 4, callCount: 6},
		{ input: Handler{ Handler: call, Name: "fourth"}, id: 5, callCount: 7},
	}

	for i, c := range tests {
		id := bus.Subscribe(c.input)
		if id != c.id {
			t.Fatalf("[%d] expect %d, got %d", i, c.id, id)
		}
	}

	for i, c := range tests {
		bus.Notify(Event{Name: c.input.Name})
		if c.callCount != callCount {
			t.Fatalf("[%d] expect %d, got %d", i, c.callCount, callCount)
		}
	}

	for i, c := range tests {
		bus.Unsubscribe(c.id)
		if bus.free != c.id {
			t.Fatalf("[%d] expect %d, got %d", i, c.id, bus.free)
		}
	}

	for k, v := range bus.live {
		if v != 0 {
			t.Fatalf("[%s] expect %d, got %d", k, 0, v)
		}
	}

	fid := bus.free
	id := bus.Subscribe(Handler{
		Name: "freelist",
		Handler: nop,
	})
	if fid != id {
		t.Fatalf("expect %d, got %d", fid, id)
	}
	first, ok := bus.live["freelist"]
	if !ok {
		t.Fatalf("expect key '%s' in map", "freelist")
	}
	if first != id {
		t.Fatalf("expect %d, got %d", id, first)
	}
}

