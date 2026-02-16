package view

import (
	

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"

	//"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/viewmodel"
)

func NewMainUI(w fyne.Window) *fyne.Container {

	uiVM := viewmodel.NewMainUI()
	formVM := viewmodel.NewBookForm(uiVM.Error, uiVM.Success)

	addButton := widget.NewButton("Create", func() {
		_ = uiVM.OpenedBody.Set(viewmodel.BodyForm)
	})
	menuButton := widget.NewButton("Menu", func() {
		_ = uiVM.OpenedBody.Set(viewmodel.BodyMenu)
	})
	tablesButton := widget.NewButton("Tables", func() {
		_ = uiVM.OpenedBody.Set(viewmodel.BodyData)
	})


	// Status Line 
	// Displays input form .Error, .Info, .Success bindings with proper colors.
	//
	statusLabel := widget.NewLabel("")
	statusLabel.Alignment = fyne.TextAlignLeading
	display := func(input binding.String, importance widget.Importance) {
		msg, _ := input.Get()
		if msg == "" {
			return
		}
		statusLabel.Importance = importance
		_ = input.Set("")
		statusLabel.SetText(msg)
	}
	uiVM.Info.AddListener(binding.NewDataListener(func() {
		display(uiVM.Info, widget.MediumImportance)
	}))
	uiVM.Error.AddListener(binding.NewDataListener(func() {
		display(uiVM.Error, widget.DangerImportance)
	}))
	uiVM.Success.AddListener(binding.NewDataListener(func() {
		display(uiVM.Success, widget.SuccessImportance)
	}))

	header := container.NewHBox(
		menuButton,
		tablesButton,
		addButton,
		statusLabel,
	)

	form := BookForm(formVM)

	tableVM := viewmodel.NewTablesVM(uiVM.Repo)
	table := BookTables(tableVM)

	body := container.NewStack(form)

	uiVM.OpenedBody.AddListener(binding.NewDataListener(func() {
		open, _ := uiVM.OpenedBody.Get()
		addButton.Enable()
		menuButton.Enable()
		tablesButton.Enable()
		switch open {
		case viewmodel.BodyForm:
			addButton.Disable()
			body.Objects[0] = form
			body.Refresh()
		case viewmodel.BodyMenu:
			menuButton.Disable()
			body.Objects[0] = widget.NewLabel("not implemented")
			body.Refresh()
		case viewmodel.BodyData:
			tablesButton.Disable()
			body.Objects[0] = table
			statusLabel.SetText("")
			body.Refresh()
		default:
			panic("opened body was not found")
		}
	}))

	frame := container.NewBorder(header, nil, nil, nil, body)

	return frame
}


func BookTables(vm *viewmodel.TablesVM) fyne.CanvasObject {

	editBtn := widget.NewButton("EDIT", nil)
	_ = editBtn

	// switch 
	tabTables := container.NewAppTabs()
	for _, table := range vm.TableNames() {
		tvm := vm.GetTable(table)
		tab := container.NewTabItem(table, BookTable(tvm))
		tabTables.Append(tab)
	}

	tabTables.OnSelected = func(tab *container.TabItem) {
		vm.SetTable(tab.Text)
	}

	return tabTables

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
		return NewColumnButton("placeholder", vm.SortBy, vm.SortOrder)
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
				object.(*ColumnButton).SetLabel(header)
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

	// Listen for updates
	vm.AddListener(binding.NewDataListener(func() {
		table.Refresh()
	}))

	top := container.NewAdaptiveGrid(1, search, column)
	return container.NewBorder(top, nil, nil, nil, table)
}


type ColumnButton struct {
	widget.Button
	minSize    fyne.Size
	label      string
	sortOrder  binding.String
	orderField binding.String
}

func NewColumnButton(label string, orderField binding.String, sortOrder binding.String) *ColumnButton {
	hb := &ColumnButton{
		label: label,
		sortOrder: sortOrder,
		orderField: orderField,
	}
	orderField.AddListener(binding.NewDataListener(func() {
		field, _ := orderField.Get()
		if field == hb.label {
			_ = sortOrder.Set("ASC")
			hb.SetText(hb.ASCLabel())
		}
		if field != hb.label {
			hb.SetText(hb.NormalLabel())
		}
	}))

	hb.OnTapped = func() {
		switch hb.Text {
		case hb.NormalLabel():
			orderField.Set(hb.label)

		case hb.ASCLabel():
			_ = sortOrder.Set("DESC")
			hb.SetText(hb.DESCLabel())

		case hb.DESCLabel():
			_ = sortOrder.Set("ASC")
			hb.SetText(hb.ASCLabel())
		}
	}
	hb.minSize = fyne.NewSize(80, 30)
	hb.ExtendBaseWidget(hb)
	return hb
}

func (hb *ColumnButton) SetLabel(s string) {
	hb.label = s
	field, _ := hb.orderField.Get()
	if field == hb.label {
		_ = hb.sortOrder.Set("ASC")
		hb.SetText(hb.ASCLabel())
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
