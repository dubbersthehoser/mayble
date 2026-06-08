package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type EnterButton struct {
	widget.Button
}

func NewEnterButton(text string, onTapped func()) *EnterButton {
	eb := &EnterButton{}
	eb.Text = text
	eb.OnTapped = onTapped
	eb.ExtendBaseWidget(eb)
	return eb
}

func (eb *EnterButton) TypedKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeyReturn, fyne.KeyEnter:
		eb.OnTapped()
	}
}
