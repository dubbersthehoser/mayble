package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"

	"fyne.io/fyne/v2/data/binding"
	"github.com/dubbersthehoser/mayble/internal/viewmodel"
)



func BookEditForm(vm *viewmodel.EditBookVM) fyne.CanvasObject {
	
	cancel := widget.NewButton("Cancel", vm.Close)
	update := widget.NewButton("Update", vm.Submit)
	
	bookEntry := newBookEntry(vm.Title, vm.Author, vm.Genre, vm.UniqueGenres)
	isRead := widget.NewCheckWithData("Is Read", vm.IsRead)
	isLoaned := widget.NewCheckWithData("On Loan", vm.IsLoaned)

	loanedEntry := newLoanEntry(vm.IsLoaned, vm.Borrower, vm.Date)
	readEntry := newReadEntry(vm.IsRead, vm.Rating, vm.Completed)

	c := container.NewVBox(
		bookEntry,
		isRead,
		readEntry,
		isLoaned,
		loanedEntry,
		container.NewHBox(update, cancel),
	)
	return c
}



func BookTables(vm *viewmodel.TablesVM) fyne.CanvasObject {

	tabTables := container.NewAppTabs()
	for _, table := range vm.TableNames() {
		tvm := vm.GetTable(table)
		tab := container.NewTabItem(table, BookTable(tvm))
		tabTables.Append(tab)
	}

	tabTables.OnSelected = func(tab *container.TabItem) {
		vm.SetTable(tab.Text)
	}
	
	edit := BookEditForm(vm.EditBookVM())

	view := container.NewStack(
		edit,
		tabTables,
	)

	edit.Hide()

	vm.EditIsOpen.AddListener(binding.NewDataListener(func() {
		isOpen, _ := vm.EditIsOpen.Get()
		if isOpen {
			edit.Show()
			tabTables.Hide()
		} else {
			edit.Hide()
			tabTables.Show()
		}
	}))


	return view
}

func BookTable(vm *viewmodel.TableVM) fyne.CanvasObject {

	//
	// Table
	//
	// Code sections (A) are to allow user to resize the last column with the header,
	// by adding an invisable header with an empty column, allows the user to move 
	// the last visable header / column to be resized with the mouse. Down side is
	// that there is an empty selectable item on the first entry of that last column.

	table := widget.NewTableWithHeaders(
		func() (rowLen, colLen int) {
			rowLen, colLen = vm.Size()
			colLen += 1 // (A) have an extra header.
			return 
		},
		func() fyne.CanvasObject {
			object := widget.NewLabel("")
			object.Truncation = fyne.TextTruncateEllipsis
			return object
		},
		func(cellID widget.TableCellID, object fyne.CanvasObject) {
			_, colLen := vm.Size()
			if cellID.Col < colLen {
				data := vm.Get(cellID.Row, cellID.Col)
				if vm.IsItemHidden(cellID.Row, cellID.Col) {
					object.(*widget.Label).SetText("")
					object.(*widget.Label).Hide()
				} else {
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
		return NewColumnButton("", vm)
	}

	for i := range vm.Headers() {
		width := vm.GetColumnWidth(i)
		table.SetColumnWidth(i, width)
	}

	table.UpdateHeader = func(cellID widget.TableCellID, object fyne.CanvasObject) {
		if cellID.Row != -1 {
			return
		}
		vm.StoreColumnWidth(cellID.Col, object.Size().Width)
		_, colLen := vm.Size()
		if cellID.Col < colLen {
			if vm.IsHeaderHidden(cellID.Col) {
				object.(*ColumnButton).Hide()
			} else {
				header := vm.Headers()[cellID.Col]
				object.(*ColumnButton).InitLabel(header)
				object.(*ColumnButton).Show()
			}
		} else { // (A) create hidden header.
			object.(*ColumnButton).Hide()
		}
	}

	// Selection
	table.OnSelected = func(id widget.TableCellID) {
	}

	column := widget.NewCheckGroup(vm.Headers(), func(list []string) {
		vm.SetHidden(list)
	})
	column.Horizontal = true

	search := widget.NewEntryWithData(vm.SearchText)

	// Listen for updates from table
	vm.AddListener(binding.NewDataListener(func() {
		table.Refresh()
	}))

	actions := container.NewHBox(
	)
	for _, a := range vm.Actions() {
		actions.Add(widget.NewButton(a.Label, a.Action))
	}

	top := container.NewAdaptiveGrid(
		1,
		search, 
		container.NewBorder(nil, nil, column, actions, column, actions),
	)
	return container.NewBorder(top, nil, nil, nil, table)
}


type ColumnButton struct {
	widget.Button
	minSize    fyne.Size
	label      string
	vm *viewmodel.TableVM
}

func NewColumnButton(label string, vm *viewmodel.TableVM) *ColumnButton {
	hb := &ColumnButton{
		label: label,
		vm: vm,
	}

	hb.OnTapped = func() {
		order, _ := vm.SortOrder.Get()
		vm.SortBy.Set(hb.label)
		if order == "ASC" && hb.ASCLabel() == hb.Text {
			_ = vm.SortOrder.Set("DESC")
		} else {
			_ = vm.SortOrder.Set("ASC")
		}
		vm.Sort()
	}

	by, _ := vm.SortBy.Get()
	if by == hb.label {
		hb.SetText(hb.ASCLabel())
	}
	hb.minSize = fyne.NewSize(80, 30)
	hb.ExtendBaseWidget(hb)
	return hb
}

func (hb *ColumnButton) InitLabel(s string) {
	if hb.label == "" {
		hb.label = s
	}	
	by, _ := hb.vm.SortBy.Get()
	ording, _ := hb.vm.SortOrder.Get()
	if by == hb.label {
		if ording == "ASC" {
			hb.SetText(hb.ASCLabel())
		} else {
			hb.SetText(hb.DESCLabel())
		}
	} else {
		hb.SetText(hb.NormalLabel())
	}
}


func (hb *ColumnButton) NormalLabel() string {
	return "- " + hb.label
}

func (hb *ColumnButton) ASCLabel() string {
	return "↑ " + hb.label
}

func (hb *ColumnButton) DESCLabel() string {
	return "↓ " + hb.label
}

func (hb *ColumnButton) MinSize() fyne.Size {
	return hb.minSize
}
