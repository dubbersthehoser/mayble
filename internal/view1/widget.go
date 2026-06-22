package view

import (

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/dubbersthehoser/mayble/internal/viewmodel1"
)

type RatingSelect struct {
	widget.Select
}

func NewRatingSelect(fn func(s string)) *RatingSelect{
	rs := &RatingSelect{}
	rs.Options = viewmodel.Ratings()
	rs.OnChanged = fn
	rs.ExtendBaseWidget(rs)
	return rs
}

type GenreEntry struct {
	widget.SelectEntry
}
func NewGenreEntry(gs *viewmodel.UniqueGenres, fn func(s string)) *GenreEntry {
	ge := &GenreEntry{}
	ge.OnChanged = fn
	ge.ExtendBaseWidget(ge)

	ge.SetOptions(gs.Genres())	
	gs.AddListener(func() {
		ge.SetOptions(gs.Genres())
	})
	return ge
}

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

type SearchEntry struct {
	widget.Entry
	OnNext func()
	OnPrev func()
}
func NewSearchEntry(next func(), prev func()) *SearchEntry {
	se := &SearchEntry{
		OnNext: next,
		OnPrev: prev,
	}
	se.ExtendBaseWidget(se)
	se.OnSubmitted = func(_ string) {se.OnNext()}
	return se
}

func (eb *SearchEntry) TypedShortcut(cut fyne.Shortcut) {
	short, ok := cut.(*desktop.CustomShortcut)
	if !ok {
		eb.Entry.TypedShortcut(cut)
		return
	}
	switch short.Mod() {
	case fyne.KeyModifierControl:
		if short.Key() == fyne.KeyReturn {
			eb.OnPrev()
		}
	}
}
