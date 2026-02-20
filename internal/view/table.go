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
	
	bookEntry := newBookEntry(vm.Title, vm.Author, vm.Genre, vm.Genres)
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



//func BookTables(vm *viewmodel.TablesVM) fyne.CanvasObject {
//
//	tabTables := container.NewAppTabs()
//	for _, table := range vm.TableNames() {
//		tvm := vm.GetTable(table)
//		tab := container.NewTabItem(table, BookTable(tvm))
//		tabTables.Append(tab)
//	}
//
//	tabTables.OnSelected = func(tab *container.TabItem) {
//		vm.SetTable(tab.Text)
//	}
//	
//	edit := BookEditForm(vm.EditBookVM())
//
//	view := container.NewStack(
//		edit,
//		tabTables,
//	)
//
//	edit.Hide()
//
//	vm.EditIsOpen.AddListener(binding.NewDataListener(func() {
//		isOpen, _ := vm.EditIsOpen.Get()
//		if isOpen {
//			edit.Show()
//			tabTables.Hide()
//		} else {
//			edit.Hide()
//			tabTables.Show()
//		}
//	}))
//
//
//	return view
//}


func fullBookTable(vmc *viewmodel.TableControllersVM, vmt *viewmodel.TableVM) fyne.CanvasObject {

	editBtn := widget.NewButton("Edit", func() {
		vmc.Edit()
	})
	deleteBtn := widget.NewButton("Delete", func() {
		vmc.Delete()
	})
	search := widget.NewEntryWithData(vmt.Search.Text)
	searchOptions := widget.NewSelect(vmt.SearchOptions(), func(s string) {
		_ = vmt.Search.Option.Set(s)
	})
	searchOptions.SetSelectedIndex(0)

	vmt.SetSelector(vmc.Selector())
	table := bookTable(vmt)
	
	hide := widget.NewCheckGroup(vmt.HideOptions(), func(list []string) {
		vmt.SetHidden(list)
	})
	hide.Horizontal = true

	controllers := container.NewVBox(
		container.NewBorder(
			nil, nil,
			searchOptions,
			nil,
			search,
		),
		container.NewBorder(
			nil, nil,
			hide,
			container.NewHBox(
				editBtn,
				deleteBtn,
			),
		),
	)

	fullTable := container.NewBorder(
		controllers,
		nil, nil, nil,
		table,
	)

	edit := BookEditForm(vmc.GetEditBook())

	view := container.NewStack(
		edit,
		fullTable,
	)
	edit.Hide()

	editBtn.Disable()
	deleteBtn.Disable()
	vmc.Selector().AddListener(binding.NewDataListener(func() {
		if vmc.Selector().HasSelected() {
			editBtn.Enable()
			deleteBtn.Enable()
		} else {
			editBtn.Disable()
			deleteBtn.Disable()
		}
	}))

	vmc.EditIsOpen.AddListener(binding.NewDataListener(func() {
		ok, _ := vmc.EditIsOpen.Get()
		if ok {
			edit.Show()
			fullTable.Hide()
		} else {
			edit.Hide()
			fullTable.Show()
		}
	}))

	return view
}

func bookTable(vm *viewmodel.TableVM) fyne.CanvasObject {

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
		vm.Select(id.Row, id.Col)
	}
	table.OnUnselected = func(id widget.TableCellID) {
		vm.Unselect(id.Row, id.Col)
		table.UnselectAll()
	}

	vm.Selector().AddListener(binding.NewDataListener(func() {
		if vm.Selector().HasSelected() {
			row, col := vm.Selector().SelectedCell()
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
	hb.label = s
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
