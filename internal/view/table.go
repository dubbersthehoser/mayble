package view

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"fyne.io/fyne/v2/data/binding"
	"github.com/dubbersthehoser/mayble/internal/viewmodel"
)

func fullBookTable(vmt *viewmodel.Table) fyne.CanvasObject {

	headers := viewmodel.NewTableHeaders(vmt)
	selector := viewmodel.NewTableSelect(vmt)
	editor := viewmodel.NewTableEdit(vmt)

	editBtn := widget.NewButton("Edit", func() {
		editor.Open()
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
		editor.Delete()
		deleteFinal.Hide()
		deleteInitial.Show()
	}

	deleteFinal.Hide()

	deleteBtn := container.NewStack(
		deleteInitial,
		deleteFinal,
	)

	// Search 
	search := widget.NewEntry()
	search.OnChanged = selector.Search
	searchOptions := widget.NewSelect(selector.SearchOptions(), selector.SetSearchBy)
	searchOptions.SetSelectedIndex(0)



	table := bookTable(vmt, headers, selector)

	// Hide Column Options
	hideLabel := widget.NewLabel("Hidden:")
	hideOptions := widget.NewCheckGroup(headers.HideOptions(), func(list []string) {
		headers.SetHidden(list)
	})
	hideOptions.Horizontal = true
	hideOptions.SetSelected(headers.GetHidden())

	// Hidden Column Bar
	hidden := container.NewBorder(nil, nil, hideLabel, nil, hideLabel, hideOptions)
	controllers := container.NewVBox(
		container.NewBorder(
			nil, nil,
			searchOptions,
			nil,
			search,
		),
		container.NewBorder(
			nil, nil,
			hidden,
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

	edit := bookEditForm(editor)

	view := container.NewStack(
		edit,
		fullTable,
	)
	edit.Hide()

	editBtn.Hide()
	deleteBtn.Hide()
	selector.AddListener(binding.NewDataListener(func() {
		if selector.HasSelected() {
			editBtn.Show()
			deleteBtn.Show()
		} else {
			editBtn.Hide()
			deleteBtn.Hide()
		}
	}))

	editor.IsOpen.AddListener(binding.NewDataListener(func() {
		ok, _ := editor.IsOpen.Get()
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

func bookEditForm(vm *viewmodel.TableEdit) fyne.CanvasObject {

	cancel := widget.NewButton("Cancel", vm.Close)
	update := widget.NewButton("Update", vm.Submit)

	bookEntry := newBookEntry(vm.Title, vm.Author, vm.Genre, vm.Genres())
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


func bookTable(vm *viewmodel.Table, headers *viewmodel.TableHeaders, selector *viewmodel.TableSelect) fyne.CanvasObject {

	//
	// Table Note
	//
	// Code sections (A) are to allow user to resize the last column with the header.
	// By adding an invisable header with an empty column, allows the user to move
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
				if vm.IsHidden(cellID.Row, cellID.Col) {
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
		return NewHeaderButton("", headers)
	}

	table.UpdateHeader = func(cellID widget.TableCellID, object fyne.CanvasObject) {
		if cellID.Row != -1 {
			return
		}
		//vm.StoreColumnWidth(cellID.Col, object.Size().Width)
		headers.StoreWidth(cellID.Col, object.Size().Width) // (2) this has to be relitive to column position
		_, colLen := vm.Size()
		if cellID.Col < colLen {
			if headers.IsHidden(cellID.Col) {
				object.(*HeaderButton).Hide()
			} else {
				header := headers.Headers()[cellID.Col]
				object.(*HeaderButton).InitLabel(header)
				object.(*HeaderButton).Show()
			}
		} else { // (A) create hidden header.
			object.(*HeaderButton).Hide()
		}
	}

	// Set the width of the columns.
	for i, label := range headers.Headers() {
		// (2) this has to be label per label. ignore relitive column position.
		width := headers.GetWidth(label)
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

func NewHeaderButton(label string, h *viewmodel.TableHeaders) *HeaderButton {
	hb := &HeaderButton{
		label: label,
		headers: h,
	}

	hb.OnTapped = func() {
		if hb.headers.GetAscending() && hb.ASCLabel() == hb.Text {
			hb.headers.SetAscending(false)
		} else {
			hb.headers.SetAscending(true)
		}
		hb.headers.Sort()
	}

	if hb.headers.GetSortBy() == hb.label {
		hb.SetText(hb.ASCLabel())
	}
	hb.minSize = fyne.NewSize(80, 30)
	hb.ExtendBaseWidget(hb)
	return hb
}

func (hb *HeaderButton) InitLabel(s string) {
	hb.label = s
	if hb.headers.GetSortBy() == hb.label {
		if hb.headers.GetAscending() {
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
