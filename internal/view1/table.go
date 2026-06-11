package view

import (

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"github.com/dubbersthehoser/mayble/internal/viewmodel1"
)


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
				data := vm.Get(cellID.Row, cellID.Col)
				object.(*widget.Label).Show()
				object.(*widget.Label).SetText(data)
				}
			} else { // (A) create empty column.
				object.(*widget.Label).SetText("")
			}
		},
	)

	table.ShowHeaderColumn = false

	table.CreateHeader = func() fyne.CanvasObject {
		return NewHeaderButton(headers)
	}

	table.UpdateHeader = func(cellID widget.TableCellID, object fyne.CanvasObject) {
		if cellID.Row != -1 {
			return
		}
		//vm.StoreColumnWidth(cellID.Col, object.Size().Width)
		headers.SetWidthWithColumn(cellID.Col, object.Size().Width) // (2) this has to be relitive to column position
		_, colLen := vm.Size()
		if cellID.Col < colLen {
			label := headers.Headers()[cellID.Col]
			object.(*HeaderButton).Update(label)
			object.(*HeaderButton).Show()
		} else { // (A) create hidden header.
			object.(*HeaderButton).Hide()
		}
	}

	// Set the width of the columns.
	for i, label := range headers.Headers() {
		// (2) this has to be label per label. ignore relitive column position.
		width := headers.GetWidthWithLabel(label)
		table.SetColumnWidth(i, width)
	}

	// Selection
	table.OnSelected = func(id widget.TableCellID) {
		selector.Select(id.Row, id.Col)
	}
	table.OnUnselected = func(id widget.TableCellID) {
		selector.Unselect()
		table.UnselectAll()
	}

	selector.AddListener(binding.NewDataListener(func() {
		if selector.HasSelected() {
			row, col := selector.Selected()
			table.Select(widget.TableCellID{Row: row, Col: col})
		} else {
			table.UnselectAll()
		}
	}))

	// Listen for updates from table
	vm.AddListener(binding.NewDataListener(func() {
		table.Refresh()
	}))

	return table
}

type HeaderButton struct {
	widget.Button
	minSize fyne.Size
	label   string
	headers *viewmodel.TableHeaders
}

func NewHeaderButton( h *viewmodel.TableHeaders) *HeaderButton {
	hb := &HeaderButton{
		headers: h,
	}

	hb.OnTapped = func() {
		hb.headers.Sort(hb.label)
	}

	hb.minSize = fyne.NewSize(80, 30)
	hb.ExtendBaseWidget(hb)
	return hb
}

func (hb *HeaderButton) Update(l string) {
	hb.label = l
	text, _ := hb.headers.Labels[l].Get()
	hb.SetText(text)
	hb.headers.Labels[hb.label].AddListener(binding.NewDataListener(func() {
		text, _ := hb.headers.Labels[hb.label].Get()
		hb.SetText(text)
	}))
}

