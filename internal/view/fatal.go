package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewFatal(banner, message, report string) *fyne.Container {
	b := widget.NewLabel(banner)
	m := widget.NewLabel(message)

	textBox := widget.NewRichTextWithText(report)
	body := container.NewVBox(b, m, textBox)
	return body
}
