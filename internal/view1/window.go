package view

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"

	"github.com/dubbersthehoser/mayble/internal/viewmodel1"
)

func NewWindow(w fyne.Window, vm *viewmodel.Window) *fyne.Container {

	w.SetMainMenu(newMainMenu(vm, w))

	status := newStatusLine(vm.StatusLine)
	controls := newControls(vm)

	topBar := container.NewGridWithColumns(2, status, controls)
	body := newBody(vm)

	view := container.NewBorder(topBar, nil, nil, nil, body)

	return view
}

func newBody(vm *viewmodel.Window) fyne.CanvasObject {

	noData := newNoData(vm)
	table := newBodyTable(vm)
	edit := newEdit(vm)
	create := newCreate(vm)

	update := func() {
		switch vm.Body.Value() {
		case viewmodel.BodyNoData:
			table.Hide()
			edit.Hide()
			create.Hide()
			noData.Show()
		case viewmodel.BodyTable:
			table.Show()
			edit.Hide()
			create.Hide()
			noData.Hide()
		case viewmodel.BodyBookCreate:
			table.Hide()
			edit.Hide()
			create.Show()
			noData.Hide()
		case viewmodel.BodyBookEdit:
			table.Hide()
			edit.Show()
			create.Hide()
			noData.Hide()
		default:
			log.Printf("Error: unexpected body type %d", vm.Body.Value())
		}
	}

	vm.Body.AddListener(update)

	update()

	body := container.NewStack(
		noData,
		table,
		edit,
		create,
	)

	return body
}


func newNoData(vm *viewmodel.Window) fyne.CanvasObject {
	// todo: needs more work.
	view := widget.NewLabel(fmt.Sprintf("database not found '%s'", vm.DBPath.Get()))
	return view
}

func newControls(vm *viewmodel.Window) fyne.CanvasObject {

	unselect := widget.NewButton("Unselect", vm.Controls.OnUnselect)
	create := widget.NewButton("Create", vm.Controls.OnCreate)
	edit := widget.NewButton("Edit", vm.Controls.OnEdit)
	
	deleteFinal := widget.NewButton("Are You Sheer?", nil)
	deleteFinal.Hide()
	deleteInitial := widget.NewButton("Delete", nil)

	deleteInitial.OnTapped = func() {
		deleteFinal.Show()
		deleteInitial.Hide()
		fyne.Do(func() {
			timer := time.NewTimer(time.Second * 2)
			<-timer.C
			fyne.Do(func() {
				deleteFinal.Hide()
				deleteInitial.Show()
			})
		})

	}

	deleteFinal.OnTapped = func() {
		vm.Controls.OnDelete()
		deleteFinal.Hide()
		deleteInitial.Show()
	}

	deleteView := container.NewStack(
		deleteInitial,
		deleteFinal,
	)

	view := container.NewHBox(
		create,
		unselect,
		edit,
		deleteView,
	)

	update := func() {
		if vm.Body.Value() != viewmodel.BodyTable {
			view.Hide()
		} else {
			view.Show()
		}
		if vm.Selected.Has() {
			println("yest")
			deleteInitial.Enable()
			deleteFinal.Enable()
			edit.Enable()
			unselect.Enable()
		} else {
			deleteInitial.Disable()
			deleteFinal.Disable()
			edit.Disable()
			unselect.Disable()
		}
	}

	vm.Selected.AddListener(update)
	vm.Body.AddListener(update)
	update()
	return view
}


func newStatusLine(vm *viewmodel.StatusLine) fyne.CanvasObject {
	label := widget.NewLabel("")
	
	vm.Text.AddListener(binding.NewDataListener(func() {
		text, _ := vm.Text.Get()
		var importance widget.Importance
		switch vm.Type {
		case viewmodel.StatusInfo:
			importance = widget.MediumImportance
		case viewmodel.StatusSuccess:
			importance = widget.SuccessImportance
		case viewmodel.StatusError:
			importance = widget.DangerImportance
		default:
			log.Printf("Error: invalid status line type %d", vm.Type)
		}
		label.Importance = importance
		label.SetText(text)
	}))
	
	vm.SetOnClear(func() {
		label.SetText("")
	})

	return label
}

func newMainMenu(vm *viewmodel.Window, w fyne.Window) *fyne.MainMenu {
	file := fyne.NewMenu("File",
		fyne.NewMenuItem("Open", func() {
			d := dialog.NewFileOpen(
				viewmodel.WrapFyneFileOpen(vm.FileManage.OpenDatabase),
				w,
			)
			d.Resize(w.Canvas().Size())
			d.SetTitleText("Open Database")
			d.SetFilter(storage.NewExtensionFileFilter([]string{".db", ".sqlite", ".sqlite3"}))
			d.Show()
		}), 

		fyne.NewMenuItem("Create", func() {
			d := dialog.NewFileSave(
				viewmodel.WrapFyneFileCreate(vm.FileManage.CreateDatabase),
				w,
			)
			d.Resize(w.Canvas().Size())
			d.SetTitleText("Create Database")
			d.SetFilter(storage.NewExtensionFileFilter([]string{".db", ".sqlite", ".sqlite3"}))
			d.Show()
			
		}),

		fyne.NewMenuItem("Import", func() {
			d := dialog.NewFileOpen(
				viewmodel.WrapFyneFileOpen(vm.FileManage.ImportFile),
				w,
			)
			d.Resize(w.Canvas().Size())
			d.SetTitleText("Import CSV")
			d.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
			d.Show()
		}),
		fyne.NewMenuItem("Export", func() {
			d := dialog.NewFileSave(
				viewmodel.WrapFyneFileCreate(vm.FileManage.ExportFile),
				w,
			)
			d.Resize(w.Canvas().Size())
			d.SetTitleText("Export CSV")
			d.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
			d.Show()
		}),
	)
	table := fyne.NewMenu("Table", 
		fyne.NewMenuItem("Show Loaned", nil),
		fyne.NewMenuItem("Show Read", nil),
	)

	menu := fyne.NewMainMenu(file, table)

	const (
		loanIdx int = iota
		readIdx
	)

	table.Items[loanIdx].Action = func() {
		if vm.ColumnSettings.IsLoanHidden() {
			vm.ColumnSettings.SetLoanHidden(true)
			table.Items[loanIdx].Checked = false
		} else {
			vm.ColumnSettings.SetLoanHidden(false)
			table.Items[loanIdx].Checked = true
		}
		menu.Refresh()
		fmt.Println("show/hide loaned columns")
	}
	table.Items[readIdx].Action = func() {
		if vm.ColumnSettings.IsReadHidden() {
			vm.ColumnSettings.SetReadHidden(true)
			table.Items[readIdx].Checked = false
		} else {
			vm.ColumnSettings.SetReadHidden(false)
			table.Items[readIdx].Checked = true
		}
		menu.Refresh()
		fmt.Println("show/hide loaned columns")
	}

	update := func() {
		if vm.Body.Value() != viewmodel.BodyTable {
			table.Items[loanIdx].Disabled = true
			table.Items[readIdx].Disabled = true
		} else {
			table.Items[loanIdx].Disabled = false
			table.Items[readIdx].Disabled = false
		}
		menu.Refresh()
	}
	update()
	return menu
}
