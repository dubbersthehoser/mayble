package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	"github.com/dubbersthehoser/mayble/internal/listing"
	"github.com/dubbersthehoser/mayble/internal/emiter"
	"github.com/dubbersthehoser/mayble/internal/gui/controller"
)

func (f *FunkView) Table() fyne.CanvasObject {

	heading := NewHeader(f.emiter, f.controller.List.OrderBy(), f.controller.List.Ordering())
	listTable := NewTableList(f.emiter, f.controller.List)

	// Table
	table := container.New(layout.NewBorderLayout(heading.view, nil, nil, nil), heading.view, listTable.List)
	return table
}

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
		hb.button.Importance = widget.MediumImportance
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
	hb.button.Importance = widget.LowImportance
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
			btn.button.Importance = widget.MediumImportance
		}
		h.buttons = append(h.buttons, btn)
		fields = append(fields, btn.button)
	}
	e.OnEvent(OnSetOrderBy, h.OnSetOrderBy)

	h.view = container.New(layout.NewGridLayout(len(fields)), fields...)
	return h
}



type TableList struct {
	Emiter     *emiter.Emiter
	List       *widget.List
	Controller *controller.BookList
}

func (tl *TableList) OnSelect(index int) {
	tl.Emiter.Emit(OnSelected, index)
}

func NewTableList(e *emiter.Emiter, controller *controller.BookList) *TableList {
	tl := &TableList{
		Emiter: e,
		Controller: controller,
	}
	tl.List = widget.NewList(tl.OnListLength, tl.OnCanvasCreation, tl.OnCanvasInit)
	tl.List.HideSeparators = false

	tl.List.OnSelected = tl.OnSelect

	listOnModification := func(_ any) {
		tl.List.UnselectAll()
	}

	listOnSearch := func(_ any) {
		idx := tl.Controller.SelectedIndex
		tl.List.Select(widget.ListItemID(idx))
	}

	listOnSelectNext := listOnSearch
	listOnSelectPrev := listOnSearch

	tl.Emiter.OnEvent(OnModification, listOnModification)
	tl.Emiter.OnEvent(OnSearch, listOnSearch)
	tl.Emiter.OnEvent(OnSelectNext, listOnSelectNext)
	tl.Emiter.OnEvent(OnSelectPrev, listOnSelectPrev)
	tl.Emiter.OnEvent(OnUnselected, func(_ any) {
		tl.List.UnselectAll()
	})
	return tl
}

func (tl *TableList) OnCanvasCreation() fyne.CanvasObject {
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
func (tl *TableList) OnCanvasInit(index int, o fyne.CanvasObject) {
	book, err := tl.Controller.Get(index)
	if err != nil {
		tl.Emiter.Emit(OnShowError, err)
		return
	}
	o.(*fyne.Container).Objects[0].(*widget.Label).SetText(book.Title)
	o.(*fyne.Container).Objects[1].(*widget.Label).SetText(book.Author)
	o.(*fyne.Container).Objects[2].(*widget.Label).SetText(book.Genre)
	o.(*fyne.Container).Objects[3].(*widget.Label).SetText(book.Ratting)
	o.(*fyne.Container).Objects[4].(*widget.Label).SetText(book.Borrower)
	o.(*fyne.Container).Objects[5].(*widget.Label).SetText(book.Date)
}
func (tl *TableList) OnListLength() int {
	return tl.Controller.Len()
}

