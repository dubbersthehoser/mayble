package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	"github.com/dubbersthehoser/mayble/internal/listing"
	"github.com/dubbersthehoser/mayble/internal/emiter"
	"github.com/dubbersthehoser/mayble/internal/gui"
	"github.com/dubbersthehoser/mayble/internal/gui/controller"
)

func (f *FunkView) Table() fyne.CanvasObject {

	heading := NewHeader(f.broker)
	listTable := NewTableList(f.broker, f.controller.List)

	table := container.New(layout.NewBorderLayout(heading.view, nil, nil, nil), heading.view, listTable.List)
	return table
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

func NewHeader(b *emiter.Broker) *Header {
	h := &Header{
		buttons: make([]*HeaderButton, 0),
	}
	labels := listing.SortByList()
	fields := make([]fyne.CanvasObject, 0)
	for _, label := range labels {
		btn := NewHeaderButton(b, label)
		h.buttons = append(h.buttons, btn)
		fields = append(fields, btn.button)
	}
	h.view = container.New(layout.NewGridLayout(len(fields)), fields...)
	return h
}




type HeaderButton struct {
	label string
	button *widget.Button
	broker *emiter.Broker
}

func NewHeaderButton(b *emiter.Broker, label string) *HeaderButton {
	hb := &HeaderButton{
		label: label,
		broker: b,
	}
	hb.button = widget.NewButton("", hb.OnTapped)
	hb.Reset()

	b.Subscribe(&emiter.Listener{
		Handler: func(e *emiter.Event) {
			by := e.Data.(listing.OrderBy)
			if by != listing.MustStringToOrderBy(hb.label) {
				hb.Reset()
				return
			}
			hb.button.Importance = widget.MediumImportance
			var (
				IsNormal bool = hb.button.Text == hb.LabelNormal()
				IsASC    bool = hb.button.Text == hb.LabelASC()
				IsDESC   bool = hb.button.Text == hb.LabelDESC() 
			)
			switch {
			case IsNormal:
				hb.SetOrdering(listing.ASC)
			case IsASC:
				hb.SetOrdering(listing.DESC)
			case IsDESC:
				hb.SetOrdering(listing.ASC)
			}
		},
	},
		gui.EventListOrderBy,
	)
	return hb
}

func (hb *HeaderButton) OnTapped() {
	hb.broker.Notify(emiter.Event{
		Name: gui.EventListOrderBy,
		Data: listing.MustStringToOrderBy(hb.label),
	})
}	

func (hb *HeaderButton) SetOrdering(o listing.Ordering) {
	switch o {
	case listing.ASC:
		hb.button.SetText(hb.LabelASC())
	case listing.DESC:
		hb.button.SetText(hb.LabelDESC())
	}
	hb.broker.Notify(emiter.Event{
		Name: gui.EventListOrdering,
		Data: o,
	})
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






type TableList struct {
	broker     *emiter.Broker
	List       *widget.List
	control *controller.BookLoanList
}

func (tl *TableList) OnSelect(index int) {
	tl.broker.Notify(emiter.Event{
		Name: gui.EventEntrySelected,
		Data: index,
	})
}

func NewTableList(b *emiter.Broker, c *controller.BookLoanList) *TableList {
	tl := &TableList{
		broker: b,
		control: c,
	}
	tl.List = widget.NewList(tl.OnListLength, tl.OnCanvasCreation, tl.OnCanvasInit)
	tl.List.HideSeparators = false
	tl.List.OnSelected = func(id widget.ListItemID) {
		tl.broker.Notify(emiter.Event{
			Name: gui.EventEntrySelected,
			Data: int(id),
		})
	}
	tl.List.OnUnselected = func(id widget.ListItemID) {
		tl.broker.Notify(emiter.Event{
			Name: gui.EventEntryUnselected,
		})
	}

	tl.broker.Subscribe(&emiter.Listener{
			Handler: func(e *emiter.Event) {
				switch  e.Name {
				case gui.EventListOrdered:
					tl.List.UnselectAll()
					tl.List.ScrollTo(0)
					tl.List.Refresh()
				case gui.EventEntrySelected:
					id := e.Data.(int)
					tl.List.Select(id)
				case gui.EventEntryUnselected:
					tl.List.UnselectAll()
				case gui.EventSelectionNone:
					tl.List.UnselectAll()
				}
			},
	},
		gui.EventListOrdered,
		gui.EventEntrySelected,
		gui.EventEntryUnselected,
		gui.EventSelectionNone,
	)

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

	genreLabel.Alignment = fyne.TextAlignCenter
	rattingLabel.Alignment = fyne.TextAlignCenter
	onLoanName.Alignment = fyne.TextAlignCenter
	onLoanDate.Alignment = fyne.TextAlignCenter

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
	book := tl.control.Get(index)
	o.(*fyne.Container).Objects[0].(*widget.Label).SetText(book.Title)
	o.(*fyne.Container).Objects[1].(*widget.Label).SetText(book.Author)
	o.(*fyne.Container).Objects[2].(*widget.Label).SetText(book.Genre)
	o.(*fyne.Container).Objects[3].(*widget.Label).SetText(book.Ratting)
	o.(*fyne.Container).Objects[4].(*widget.Label).SetText(book.Borrower)
	o.(*fyne.Container).Objects[5].(*widget.Label).SetText(book.Date)
}
func (tl *TableList) OnListLength() int {
	return tl.control.Len()
}

