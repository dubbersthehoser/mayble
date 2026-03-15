package viewmodel

import (
	"fmt"
	"slices"
	"testing"

	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/database"
)

func TestMainUI(t *testing.T) {
	db, err := database.OpenMem()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Conn.Close()
	cfg := &config.Config{}
	err = db.Conn.Ping()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	mainUI := NewMainUI(cfg, db, []error{})

	if mainUI.OpenedBody == nil {
		t.Fatal("unexpected nil")
	}
	if mainUI.DBFile == nil {
		t.Fatal("unexpected nil")
	}
	if mainUI.Error == nil {
		t.Fatal("unexpected nil")
	}
	if mainUI.Info == nil {
		t.Fatal("unexpected nil")
	}
	if mainUI.Clear == nil {
		t.Fatal("unexpected nil")
	}

	msgTests := []struct {
		name   string
		expect string
		getter binding.String
	}{
		{
			name:   msgUserInfo,
			expect: "Info",
			getter: mainUI.Info,
		},
		{
			name:   msgUserSuccess,
			expect: "Success",
			getter: mainUI.Success,
		},
		{
			name:   msgUserError,
			expect: "Error",
			getter: mainUI.Error,
		},
	}

	for _, c := range msgTests {
		t.Run(c.expect, func(t *testing.T) {
			mainUI.bus.Notify(bus.Event{
				Name: c.name,
				Data: c.expect,
			})
			actual, _ := c.getter.Get()
			if actual != c.expect {
				t.Fatalf("expect '%s', got '%s'", c.expect, actual)
			}
		})

	}
}

func Test_formatRating(t *testing.T) {
	tests := []struct {
		input  int
		expect string
	}{
		{
			input:  0,
			expect: "",
		},
		{
			input:  1,
			expect: "⭐",
		},
		{
			input:  2,
			expect: "⭐⭐",
		},
		{
			input:  3,
			expect: "⭐⭐⭐",
		},
		{
			input:  4,
			expect: "⭐⭐⭐⭐",
		},
		{
			input:  5,
			expect: "⭐⭐⭐⭐⭐",
		},
		{
			input:  6,
			expect: "ERROR",
		},
		{
			input:  7,
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
