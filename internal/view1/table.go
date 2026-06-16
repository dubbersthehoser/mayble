package view

import (

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"

	"github.com/dubbersthehoser/mayble/internal/viewmodel1"
	"github.com/dubbersthehoser/mayble/internal/models"
)

func newBodyTable(vm *viewmodel.Window) fyne.CanvasObject {

	search := widget.NewEntry()
	searchBy := widget.NewSelect(
		vm.Searching.GetOptions(),
		vm.Searching.SetBy,
	)

	searchBy.SetSelected(vm.Searching.GetOptions()[0])

	top := container.NewGridWithColumns(2, search, searchBy)
	table := container.NewStack(newTable(vm))
	body := container.NewBorder(top, nil,nil,nil, table)

	vm.ColumnSettings.AddListener(func() {
		table.Objects[0] = newTable(vm)
		table.Refresh()
	})

	return body
}

//func newTable(vm *viewmodel.Table, headers *viewmodel.TableHeaders, selector *viewmodel.TableSelect) fyne.CanvasObject {
func newTable(vm *viewmodel.Window) fyne.CanvasObject {

	//
	// Table Note
	//
	// Code sections (A) are to allow user to resize the last column with the header.
	// By adding an invisable header with an empty column, allows the user to move
	// the last visable header / column to be resized with the mouse. Down side is
	// that there is an empty selectable item on the first entry of that last column.

	table := widget.NewTableWithHeaders(
		func() (rowLen, colLen int) {
			rowLen, colLen = vm.DataTable.Size()
			colLen += 1 // (A) have an extra header.
			return
		},
		func() fyne.CanvasObject {
			object := widget.NewLabel("")
			object.Truncation = fyne.TextTruncateEllipsis
			return object
		},
		func(cellID widget.TableCellID, object fyne.CanvasObject) {
			_, colLen := vm.DataTable.Size()
			if cellID.Col < colLen {
				data := vm.DataTable.Get(cellID.Row, cellID.Col)
				object.(*widget.Label).Show()
				object.(*widget.Label).SetText(data)

			} else { // (A) create empty column.
				object.(*widget.Label).SetText("")
			}
		},
	)

	table.ShowHeaderColumn = false
	header := NewHeader(vm)

	table.CreateHeader = func() fyne.CanvasObject {
		return header.NewHeaderButton()
	}

	table.UpdateHeader = func(cellID widget.TableCellID, object fyne.CanvasObject) {
		if cellID.Row != -1 {
			return
		}

		_, colLen := vm.DataTable.Size()
		if cellID.Col < colLen {
			vm.ColumnSettings.SetWidth(vm.ColumnSettings.Headers()[cellID.Col], object.Size().Width)
			label := models.BookEntryFields()[cellID.Col]
			by := vm.Sorting.GetOrderBy()
			asc := vm.Sorting.GetAscending()
			object.(*HeaderButton).Update(label, by, asc)
			object.(*HeaderButton).Show()
		} else { // (A) create hidden header.
			object.(*HeaderButton).Hide()
		}
	}

	// Set the width of the columns.
	for i, label := range models.BookEntryFields() {
		width := vm.ColumnSettings.GetWidth(label)
		table.SetColumnWidth(i, width)
	}

	// Selection
	table.OnSelected = func(id widget.TableCellID) {
		vm.Selected.Select(id.Row, id.Col)
	}
	table.OnUnselected = func(id widget.TableCellID) {
		vm.Selected.Unselect()
		table.UnselectAll()
	}

	vm.Selected.AddListener(func() {
		if vm.Selected.Has() {
			row, col := vm.Selected.Get()
			table.Select(widget.TableCellID{Row: row, Col: col})
		} else {
			table.UnselectAll()
		}
	})

	// Listen for updates from table
	vm.DataTable.AddListener(func() {
		table.Refresh()
	})

	return table
}

type Header struct {
	vm      *viewmodel.Window
	buttons []*HeaderButton
	minSize fyne.Size

}

func NewHeader(vm *viewmodel.Window) *Header {
	h := &Header{
		vm: vm,
		buttons: make([]*HeaderButton, 0),
	}
	return h
}

func (h *Header) NewHeaderButton() *HeaderButton {
	minSize := fyne.NewSize(
		h.vm.ColumnSettings.MinWidth(),
		25.0,
	)
	hb := NewHeaderButton(h, minSize)
	h.buttons = append(h.buttons, hb)
	return hb
}

func (h *Header) Pressed(label string) {
	by := h.vm.Sorting.GetOrderBy()
	asc := h.vm.Sorting.GetAscending()

	if by == label {
		asc = !asc
	} else {
		by = label
		asc = false
	}

	h.vm.Sorting.SetOrderBy(by)
	h.vm.Sorting.SetAscending(asc)

	for _, btn := range h.buttons {
		btn.Update(btn.label, by, asc)
	}

	h.vm.Sorting.Sort()
}

type HeaderButton struct {
	widget.Button
	header *Header
	minSize fyne.Size
	label   string
}

func NewHeaderButton(h *Header, minSize fyne.Size) *HeaderButton {
	hb := &HeaderButton{
		header: h,
		minSize: minSize,
	}

	hb.OnTapped = func() {
		hb.header.Pressed(hb.label)
	}

	hb.minSize = fyne.NewSize(80, 30)
	hb.ExtendBaseWidget(hb)

	return hb
}

func (hb *HeaderButton) Update(label string, by string, asc bool) {
	hb.label = label
	if label == by {
		if asc {
			hb.SetText("↑ " + label)
		} else {
			hb.SetText("↓ " + label)
		}
	} else {
		hb.SetText("- " + label)
	}
}

func (hb *HeaderButton) MinSize() fyne.Size {
	return hb.minSize
}
