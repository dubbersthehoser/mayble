package view

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	_"fyne.io/fyne/v2/data/binding"
	_"fyne.io/fyne/v2/theme"
	_ "fyne.io/fyne/v2/canvas"

	"github.com/dubbersthehoser/mayble/internal/listing"
	"github.com/dubbersthehoser/mayble/internal/emiter"
)

type HeaderButton struct {
	label string
	button *widget.Button
	emiter *emiter.Emiter
}

func NewHeaderButton(e *emiter.Emiter, label string) *HeaderButton {
	hb := &HeaderButton{
		label: label,
		emiter: e,
	}
	hb.button = widget.NewButton("", hb.OnTapped)
	hb.Reset()
	return hb
}

func (hb *HeaderButton) OnTapped() {
	var (
		IsNormal bool = hb.button.Text == hb.LabelNormal()
		IsASC    bool = hb.button.Text == hb.LabelASC()
		IsDESC   bool = hb.button.Text == hb.LabelDESC() 
	)
	switch {
	case IsNormal:
		hb.button.SetText(hb.LabelASC())
		hb.emiter.Emit(OnSetOrderBy, hb.label)
		hb.emiter.Emit(OnSetOrdering, listing.ASC)

	case IsASC:
		hb.button.SetText(hb.LabelDESC())
		hb.emiter.Emit(OnSetOrdering, listing.DESC)

	case IsDESC:
		hb.button.SetText(hb.LabelASC())
		hb.emiter.Emit(OnSetOrdering, listing.ASC)
	}
}

func (hb *HeaderButton) SetOrdering(o listing.Ordering) {
	switch o {
	case listing.ASC:
		hb.button.SetText(hb.LabelASC())
	case listing.DESC:
		hb.button.SetText(hb.LabelDESC())
	}
}

func (hb *HeaderButton) LabelNormal() string {
	return "- " + hb.label
}
func (hb *HeaderButton) LabelASC() string {
	return "↑ " + hb.label
}
func (hb *HeaderButton) LabelDESC() string {
	return "↓ " + hb.label
}
func (hb *HeaderButton) Reset() {
	hb.button.SetText(hb.LabelNormal())
}

type Header struct {
	buttons []*HeaderButton
	view    fyne.CanvasObject
}

func (h *Header) OnSetOrderBy(data any) {
	var orderBy listing.OrderBy
	switch v := data.(type) {
	case string:
		orderBy = listing.MustStringToOrderBy(v)
	case listing.OrderBy:
		orderBy = v
	default:
		panic("invalid data OnOrderBy event")
	}
	for _, btn := range h.buttons {
		if btn.label != string(orderBy) {
			btn.Reset()
		}
	}
}

func NewHeader(e *emiter.Emiter, by listing.OrderBy, o listing.Ordering) *Header {
	h := &Header{
		buttons: make([]*HeaderButton, 0),
	}
	labels := listing.SortByList()
	fields := make([]fyne.CanvasObject, 0)
	for _, label := range labels {
		btn := NewHeaderButton(e, label)
		if by == listing.MustStringToOrderBy(label) {// synced to booklist state
			btn.SetOrdering(o)
		}
		h.buttons = append(h.buttons, btn)
		fields = append(fields, btn.button)
	}
	e.OnEvent(OnSetOrderBy, h.OnSetOrderBy)


	h.view = container.New(layout.NewGridLayout(len(fields)), fields...)
	return h
}


func (f *FunkView) Table() fyne.CanvasObject {

	heading := NewHeader(f.emiter, f.controller.List.OrderBy(), f.controller.List.Ordering())

	/***************************
		List's Functions
	****************************/
	OnListLength := func() int {
		n := f.controller.List.Len()
		return n
	}

	OnCanvasCreation := func() fyne.CanvasObject {

		titleLabel := widget.NewLabel("")
		authorLabel := widget.NewLabel("")
		genreLabel := widget.NewLabel("")
		rattingLabel := widget.NewLabel("")
		onLoanName := widget.NewLabel("")
		onLoanDate := widget.NewLabel("")

		titleLabel.Wrapping = fyne.TextWrapWord
		authorLabel.Wrapping = fyne.TextWrapWord
		genreLabel.Wrapping = fyne.TextWrapWord
		rattingLabel.Wrapping = fyne.TextWrapWord
		onLoanName.Wrapping = fyne.TextWrapWord
		onLoanDate.Wrapping = fyne.TextWrapWord
		

		titleLabel.Truncation = fyne.TextTruncateEllipsis
		authorLabel.Truncation = fyne.TextTruncateEllipsis
		genreLabel.Truncation = fyne.TextTruncateEllipsis
		rattingLabel.Truncation = fyne.TextTruncateEllipsis
		onLoanName.Truncation = fyne.TextTruncateEllipsis
		onLoanDate.Truncation = fyne.TextTruncateEllipsis

		titleLabel.Selectable = false
		authorLabel.Selectable = false
		genreLabel.Selectable = false
		rattingLabel.Selectable = false
		onLoanName.Selectable = false
		onLoanDate.Selectable = false

		fields := []fyne.CanvasObject{
			titleLabel,
			authorLabel,
			genreLabel,
			rattingLabel,
			onLoanName,
			onLoanDate,
		}
		entry := container.New(layout.NewGridLayout(len(fields)), fields...)
		return entry
	}   
	OnCanvasInit := func(index int, o fyne.CanvasObject) {
		book, err := f.controller.List.Get(index)
		if err != nil {
			f.displayError(err)
			return
		}
		println(book.Title)
		o.(*fyne.Container).Objects[0].(*widget.Label).SetText(book.Title)
		o.(*fyne.Container).Objects[1].(*widget.Label).SetText(book.Author)
		o.(*fyne.Container).Objects[2].(*widget.Label).SetText(book.Genre)
		o.(*fyne.Container).Objects[3].(*widget.Label).SetText(book.Ratting)
		o.(*fyne.Container).Objects[4].(*widget.Label).SetText(book.Borrower)
		o.(*fyne.Container).Objects[5].(*widget.Label).SetText(book.Date)
	} 

	OnSelect := func(index int) {
		f.emiter.Emit(OnSelected, index)
	}
	OnUnselect := func(index int) {
		f.emiter.Emit(OnUnselected, nil)
	}

	/*******************
		List
	********************/
	List := widget.NewList(OnListLength, OnCanvasCreation, OnCanvasInit)
	List.HideSeparators = false
	List.OnSelected = OnSelect
	List.OnUnselected = OnUnselect

	listOnModification := func(_ any) {
		List.UnselectAll()
	}

	listOnSearch := func(_ any) {
		fmt.Println("list on search")
		idx := f.controller.List.SelectedIndex
		List.Select(widget.ListItemID(idx))
	}
	listOnSort := func(_ any) {
		println("table sort")
		List.Refresh()
	}

	listOnSelectNext := listOnSearch
	listOnSelectPrev := listOnSearch

	f.emiter.OnEvent(OnModification, listOnModification)
	f.emiter.OnEvent(OnSearch, listOnSearch)
	f.emiter.OnEvent(OnSelectNext, listOnSelectNext)
	f.emiter.OnEvent(OnSelectPrev, listOnSelectPrev)
	f.emiter.OnEvent(OnSort, listOnSort)

	// Table
	table := container.New(layout.NewBorderLayout(heading.view, nil, nil, nil), heading.view, List)
	return table
}
