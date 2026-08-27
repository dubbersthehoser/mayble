package view

import (
	"log"

	"fyne.io/fyne/v2"
	_ "fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/dubbersthehoser/mayble/internal/viewmodel"
	"github.com/dubbersthehoser/mayble/internal/viewmodel/table"
)

func newBodyTable(vm *viewmodel.Window) fyne.CanvasObject {

	search := NewSearchEntry(
		func() {
			vm.Table.NextSearchedItem()
		},
		func() {
			vm.Table.PrevSearchedItem()
		},
	)
	search.OnChanged = vm.Table.Search
	searchBy := widget.NewSelect(
		vm.Table.SearchableOptions(),
		vm.Table.SearchableSetBy,
	)
	searchBy.SetSelected(vm.Table.SearchableOptions()[0])

	top := container.NewGridWithColumns(2, search, searchBy)
	// Wrapped table view with stack layout, so table can be changed with out needing to know its exact location in body container.
	tbl := container.NewStack(newTable(vm))
	body := container.NewBorder(top, nil, nil, nil, tbl)

	// Create a new table when column is hidden.
	vm.Table.AddListenerOnHidden(func() {
		var (
			searchByIdx int = 1
		)
		tbl.RemoveAll()
		go func(){
			vm.Table.NewSheet()
			fyne.DoAndWait(func() {
				ntbl := newTable(vm)
				tbl.Add(ntbl)
				tbl.Refresh()
			})
		}()
		top.Objects[searchByIdx].(*widget.Select).SetOptions(
			vm.Table.SearchableOptions(),
		)
		top.Objects[searchByIdx].(*widget.Select).SetSelectedIndex(0)
	})

	return body
}

type Table struct {
	widget.Table
	vm     *viewmodel.Window
	header *Header
}

func newTable(vm *viewmodel.Window) *Table {

	// Code sections (A) are to allow user to resize the last column with the header.
	// By adding an invisible header with an empty column, allows the user to move
	// the last visible header / column to be resized with the mouse. Down side is
	// that there is an empty selectable item on the first entry of that last column.
	// The selection system will ignore this selection of those cells.

	Length := func() (rowLen, colLen int) {
		rowLen, colLen = vm.Table.Size()
		colLen += 1 // (A) have an extra header.
		return
	}
	CreateCell := func() fyne.CanvasObject {
		object := widget.NewLabel("")
		object.Truncation = fyne.TextTruncateEllipsis
		return object
	}
	UpdateCell := func(cellID widget.TableCellID, object fyne.CanvasObject) {
		_, colLen := vm.Table.Size()
		if cellID.Col < colLen {
			point := table.Point{Row: cellID.Row, Col: cellID.Col}
			data, err := vm.Table.GetDataPoint(point)
			if err != nil {
				log.Println("view table:", err)
				data = "ERROR"
			}
			object.(*widget.Label).Show()
			object.(*widget.Label).SetText(data)

		} else { // (A) create empty cell.
			object.(*widget.Label).SetText("")
		}
	}

	tbl := &Table{
		Table: widget.Table{
			Length:           Length,
			CreateCell:       CreateCell,
			UpdateCell:       UpdateCell,
			ShowHeaderColumn: false,
			ShowHeaderRow:    true,
		},
	}
	tbl.ExtendBaseWidget(tbl)

	// Header buttons
	tbl.header = NewHeader(vm)

	tbl.CreateHeader = func() fyne.CanvasObject {
		return tbl.header.NewHeaderButton()
	}

	tbl.UpdateHeader = func(cellID widget.TableCellID, object fyne.CanvasObject) {
		if cellID.Row != -1 {
			return
		}

		_, colLen := vm.Table.Size()
		if cellID.Col < colLen {
			// set size
			label := vm.Table.Header()[cellID.Col]
			btnWidth := object.Size().Width
			cfgWidth := vm.Table.GetColumnWidth(label)
			if btnWidth != cfgWidth {
				if object.(*HeaderButton).label != label {
					tbl.SetColumnWidth(cellID.Col, cfgWidth)
				} else {
					vm.Table.SetColumnWidth(label, btnWidth)
				}
			}
			// update header button
			by := vm.Table.GetSortingOrderBy()

			asc := vm.Table.GetSortingAscending()
			object.(*HeaderButton).Update(label, by, asc)
			object.(*HeaderButton).Show()
		} else { // (A) create hidden header.
			object.(*HeaderButton).Hide()
		}
	}

	// Set the width of the columns from settings.
	for i, label := range vm.Table.Header() {
		width := vm.Table.GetColumnWidth(label)
		tbl.SetColumnWidth(i, width)
	}

	// Selection Events
	tbl.OnSelected = func(id widget.TableCellID) {
		point := table.Point{Row: id.Row, Col: id.Col}
		vm.Table.SelectPoint(point)
	}

	tbl.OnUnselected = func(id widget.TableCellID) {
		vm.Table.UnselectPoint()
		tbl.UnselectAll()
	}

	// Listen for select events, then select, or unselect.
	vm.Table.AddSelectedListener(func(has bool, point table.Point) {
		if has {
			maxRow, maxCol := vm.Table.Size()
			if point.Row >= maxRow || point.Col >= maxCol { // (A) unselect the hidden cell if selected.
				id := widget.TableCellID{Row: point.Row, Col: point.Col}
				tbl.Unselect(id)
				return
			}
			tbl.Select(widget.TableCellID{Row: point.Row, Col: point.Col})

		} else {
			tbl.UnselectAll()
		}
	})

	// Listen for updates from table
	vm.Table.AddListenerOnDataChanged(func() {
		tbl.UnselectAll()
		tbl.Refresh()
	})
	return tbl
}

type Header struct {
	vm      *viewmodel.Window
	buttons []*HeaderButton
}

func NewHeader(vm *viewmodel.Window) *Header {
	h := &Header{
		vm:      vm,
		buttons: make([]*HeaderButton, 0),
	}
	return h
}

func (h *Header) NewHeaderButton() *HeaderButton {
	hb := NewHeaderButton(h)
	h.buttons = append(h.buttons, hb)
	return hb
}

func (h *Header) Pressed(label string) {

	by := h.vm.Table.GetSortingOrderBy()
	asc := h.vm.Table.GetSortingAscending()

	if by == label {
		asc = !asc
	} else {
		by = label
		asc = false
	}

	h.vm.Table.SetSortingOrderBy(by)
	h.vm.Table.SetSortingAscending(asc)

	for _, btn := range h.buttons {
		btn.Update(btn.label, by, asc)
	}

	h.vm.Table.Sort()
}

type HeaderButton struct {
	widget.Button
	header *Header
	label  string
}

func NewHeaderButton(h *Header) *HeaderButton {

	hb := &HeaderButton{
		header: h,
	}

	hb.OnTapped = func() {
		hb.header.Pressed(hb.label)
	}

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

	hb.Refresh()
}

func (hb *HeaderButton) MinSize() fyne.Size {
	minWidth := hb.header.vm.Table.HeaderMinWidth()
	height := hb.header.vm.Table.HeaderHeight()
	minSize := fyne.Size{
		Width:  minWidth,
		Height: height,
	}
	return minSize
}
