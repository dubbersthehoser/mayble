package viewmodel


import (
	"time"
	"testing"
)


func TestDataItem(t *testing.T) {
	item := newDataItem(123, "hello")

	v, err := item.AsView()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if v != "hello" {
		t.Fatalf("expect '%s', got '%s'", "hello", v)
	}

	if item.GetID() != 123 {
		t.Fatalf("expect %d, got %d", 123, item.GetID())
	}

	if item.nextItem != nil {
		t.Fatalf("expect %v, got %v", nil, item.nextItem)
	}

	item2 := newDataItem(32, 21)

	item.setNext(item2)

	if item.nextItem == nil {
		t.Fatal("expected spot to be filled")
	}

	v, err = item2.AsView()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if v != "21" { // Note: This will be changed to use stars.
		t.Fatalf("expect \"21\", got \"%s\"", v)
	}


	date := time.Date(2000, 02, 01, 0,0,0,0, time.UTC)
	item3 := newDataItem(43, &date)

	v, err = item3.AsView()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if v != "01/02/2000" {
		t.Fatalf("expect '01/02/2000', got %s", v)
	}

}


func TestDataHeader(t *testing.T) {

	dt := newDataHeader("Test")

	r := dt.size()
	if r != 0 {
		t.Fatalf("expect 0, got %d", r)
	}

	values := []string{
		"hello",
		"jack",
		"will",
	}

	for i := range values{
		item := newDataItem(int64(i), "hello")
		dt.append(item)
		if dt.size() != i+1 {
			t.Fatalf("expect %d, got %d", i+1, dt.size())
		}

		e, err := dt.get(i).AsView()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		v, err := item.AsView()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if v != e {
			t.Fatalf("expect '%s', got '%s'", v, e)
			
		}
	}
}


func TestDataTable(t *testing.T) {
	
}
