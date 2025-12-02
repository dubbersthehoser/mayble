package emiter

import (
	"testing"
)


func TestBroker(t *testing.T) {

	GetTestHandler := func(count *int) func(*Event) {
		return func(e *Event) {
			t.Log(e.Name)
			*count++
		}
	}
	
	broker := Broker{}
	events := []string{
		"create",
		"delete",
		"update",
		"show",
	}

	var calls int
	
	handler := GetTestHandler(&calls)

	listener := Listener{
		Handler: handler,
	}

	id := broker.Subscribe(&listener, events...)

	broker.Notify(Event{Name: "create"})

	if calls != 1{
		t.Fatalf("expect call count %d, got %d", 1, calls)
	}
	broker.Notify(Event{Name: "delete"})
	broker.Notify(Event{Name: "update"})
	broker.Notify(Event{Name: "show"})
	broker.Notify(Event{Name: "nothing"})
	if calls != 4 {
		t.Fatalf("expect call count %d, got %d", 4, calls)
	}

	broker.Unsubscribe(id, "delete", "nothing")
	broker.Notify(Event{Name: "delete"})
	if calls != 4 {
		t.Fatalf("expect call count %d, got %d", 4, calls)
	}

	broker.Unsubscribe(id, "update", "show", "create")
	broker.Notify(Event{Name: "update"})
	broker.Notify(Event{Name: "show"})
	broker.Notify(Event{Name: "nothing"})
	if calls != 4 {
		t.Fatalf("expect call count %d, got %d", 4, calls)
	}
}
