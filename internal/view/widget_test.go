package view

import (
	"testing"

	"fyne.io/fyne/v2"
)

func TestEnterButton(t *testing.T) {

	count := 0
	expect := 2

	eb := NewEnterButton("Test", func() { count += 1 })

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

func Test_backspaceWord(t *testing.T) {
	tests := []struct{
		name   string
		text   string
		cursor int
		expect int
	}{
		{
			name: "first-test",
			text: "  remove  ",
			cursor: 10,
			expect: 8,
		},
		{
			name: "second-test",
			text: "remove  ",
			cursor: 8,
			expect: 8,
		},
		{
			name: "third-test",
			text: "  remove  ",
			cursor: 4,
			expect: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := backspaceWord(tt.text, tt.cursor)
			if actual != tt.expect {
				t.Fatalf("expect %d, got %d", tt.expect, actual)
			}
		})
	}


}
