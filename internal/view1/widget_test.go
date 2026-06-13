package view

import (
	"testing"

	"fyne.io/fyne/v2"
)

func TestEnterButton(t *testing.T) {

	count := 0
	expect := 2
	
	eb := NewEnterButton("Test", func(){ count += 1})

	evReturn := &fyne.KeyEvent{
		Name: fyne.KeyReturn,
	}

	evEnter := &fyne.KeyEvent{
		Name: fyne.KeyEnter,
	}

	eb.TypedKey(evReturn)
	eb.TypedKey(evEnter)

	if count != expect {
		t.Fatalf("expect %d, got %d", expect, count)
	}
}
