package view

import (
	"time"

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

	loanedEntry := newLoanEntry(vm.IsLoaned, vm.Date, vm.Borrower)
	readEntry := newReadEntry(vm.IsRead, vm.Completed, vm.Rating)

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


func fullBookTable(vmc *viewmodel.TableControllersVM, vmt *viewmodel.TableVM) fyne.CanvasObject {

	editBtn := widget.NewButton("Edit", func() {
		vmc.Edit()
	})

	deleteFinal := widget.NewButton("Are You Sheer?", nil)
	deleteInitial := widget.NewButton("Delete", nil)

	deleteInitial.OnTapped = func() {
		deleteFinal.Show()
		deleteInitial.Hide()

		go func() {
			timer := time.NewTimer(time.Second * 2)
			<-timer.C
			fyne.Do(func() {
				deleteFinal.Hide()
				deleteInitial.Show()
			})
		}()

	}
	deleteFinal.OnTapped = func() {
		vmc.Delete()
		deleteFinal.Hide()
		deleteInitial.Show()
	}

	deleteFinal.Hide()

	deleteBtn := container.NewStack(
		deleteInitial,
		deleteFinal,
	)


	search := widget.NewEntryWithData(vmt.Search.Text)
	searchOptions := widget.NewSelect(vmt.SearchOptions(), func(s string) {
		_ = vmt.Search.Option.Set(s)
	})
	searchOptions.SetSelectedIndex(0)

	vmt.SetSelector(vmc.Selector())
	table := bookTable(vmt)
	
	hide := widget.NewCheckGroup(vmt.HiddenOptions(), func(list []string) {
		vmt.SetHidden(list)
	})
	hide.Horizontal = true

	hide.SetSelected(vmt.Hidden())

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

	editBtn.Hide()
	deleteBtn.Hide()
	vmc.Selector().AddListener(binding.NewDataListener(func() {
		if vmc.Selector().HasSelected() {
			editBtn.Show()
			deleteBtn.Show()
		} else {
			editBtn.Hide()
			deleteBtn.Hide()
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
		return NewHeaderButton("", vm)
	}

	table.UpdateHeader = func(cellID widget.TableCellID, object fyne.CanvasObject) {
		if cellID.Row != -1 {
			return
		}
		vm.StoreColumnWidth(cellID.Col, object.Size().Width)
		_, colLen := vm.Size()
		if cellID.Col < colLen {
			if vm.IsHeaderHidden(cellID.Col) {
				object.(*HeaderButton).Hide()
			} else {
				header := vm.Headers()[cellID.Col]
				object.(*HeaderButton).InitLabel(header)
				object.(*HeaderButton).Show()
			}
		} else { // (A) create hidden header.
			object.(*HeaderButton).Hide()
		}
	}

	// Set the width of the columns.
	for i := range vm.Headers() {
		width := vm.GetColumnWidth(i)
		table.SetColumnWidth(i, width)
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


type HeaderButton struct {
	widget.Button
	minSize    fyne.Size
	label      string
	vm *viewmodel.TableVM
}

func NewHeaderButton(label string, vm *viewmodel.TableVM) *HeaderButton {
	hb := &HeaderButton{
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

func (hb *HeaderButton) InitLabel(s string) {
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


func (hb *HeaderButton) NormalLabel() string {
	return "- " + hb.label
}

func (hb *HeaderButton) ASCLabel() string {
	return "↑ " + hb.label
}

func (hb *HeaderButton) DESCLabel() string {
	return "↓ " + hb.label
}

func (hb *HeaderButton) MinSize() fyne.Size {
	return hb.minSize
}
