package viewmodel

import (
	"slices"
	"testing"
	"fmt"

	"fyne.io/fyne/v2/data/binding"
)

func Test_formatRating(t *testing.T) {
	tests := []struct{
		input  int
		expect string
	}{
		{
			input: 0,
			expect: "",
		},
		{
			input: 1,
			expect: "⭐",
		},
		{
			input: 2,
			expect: "⭐⭐",
		},
		{
			input: 3,
			expect: "⭐⭐⭐",
		},
		{
			input: 4,
			expect: "⭐⭐⭐⭐",
		},
		{
			input: 5,
			expect: "⭐⭐⭐⭐⭐",
		},
		{
			input: 6,
			expect: "ERROR",
		},
		{
			input: 7,
			expect: "ERROR",
		},
	}

	for i, c := range tests {
		t.Run(fmt.Sprintf("case#%d", i), func(t *testing.T) {
			actual := formatRating(c.input)
			if actual != c.expect {
				t.Fatalf("expect '%s', got '%s'", c.expect, actual)
			}
		})
	}
}

func TestRattings(t *testing.T) {
	expect := []string{
		"",
		"⭐",
		"⭐⭐",
		"⭐⭐⭐",
		"⭐⭐⭐⭐",
		"⭐⭐⭐⭐⭐",
	}
	actual := Ratings()
	if !slices.Equal(expect, actual) {
		t.Fatalf("expect\n%#v\n  got\n%#v", expect, actual)
	}
}

func Test_listener(t *testing.T) {
	l := &listener{}
	var callCount int

	listeners := []binding.DataListener{
		binding.NewDataListener(func() {
			callCount += 1
		}),
		binding.NewDataListener(func() {
			callCount += 1
		}),
		binding.NewDataListener(func() {
			callCount += 1
		}),
		binding.NewDataListener(func() {
			callCount += 1
		}),
	}

	l.AddListener(listeners[0])
	l.notify()

	if callCount != 1 {
		t.Fatalf("expect %d, got %d", 1, callCount)
	}

	if len(l.listeners) != 1 {
		t.Fatalf("expect %d, got %d", 1, len(l.listeners))
	}

	l.AddListener(listeners[1])
	l.AddListener(listeners[2])
	l.AddListener(listeners[3])

	if len(l.listeners) != 4 {
		t.Fatalf("expect %d, got %d", 4, len(l.listeners))
	}

	for i := range listeners {
		l.RemoveListener(listeners[i])
	}

	if len(l.listeners) != 0 {
		t.Fatalf("expect %d, got %d", 0, len(l.listeners))
	}

	
}
