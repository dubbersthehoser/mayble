package view

import (

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"

	//"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/gui/viewmodel"
)

func NewMainUI(w fyne.Window) *fyne.Container {

	uiVM := viewmodel.NewMainUI()
	formVM := viewmodel.NewBookForm(uiVM.Error, uiVM.Success)

	addButton := widget.NewButton("SUBMIT", func() {
		_ = uiVM.OpenedBody.Set(viewmodel.BodyForm)
	})
	menuButton := widget.NewButton("MENU", func() {
		_ = uiVM.OpenedBody.Set(viewmodel.BodyMenu)
	})
	tablesButton := widget.NewButton("DATA", func() {
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
	tableVM := viewmodel.NewTable()
	table := BookTable(w, tableVM)

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

func BookTables() *fyne.Container {
	return nil
}

func BookTable(w fyne.Window, vm *viewmodel.Table) fyne.CanvasObject {

	table := widget.NewTableWithHeaders(
		func() (row, col int) {
			row = vm.Length()
			col = len(vm.Header)
			return row, col
		},
		func() fyne.CanvasObject {
			object := widget.NewLabel("placeholder")
			return object
		},
		func(cellID widget.TableCellID, object fyne.CanvasObject) {
			bind := vm.Items[cellID.Row][cellID.Col]
			object.(*widget.Label).SetText(bind)
		},
	)
	table.ShowHeaderColumn = false

	table.CreateHeader = func() fyne.CanvasObject {
		return widget.NewButton("", nil)
	}
	table.UpdateHeader = func(cellID widget.TableCellID, object fyne.CanvasObject) {
		if cellID.Row != -1 {
			return
		}
		header := vm.Header[cellID.Col]
		object.(*widget.Button).SetText(header)
	}
	
	for i, size := range vm.Sizes {
		table.SetColumnWidth(i, float32(size * 10))
	}
	
	return container.NewBorder(nil, nil, nil, nil, table)
}




