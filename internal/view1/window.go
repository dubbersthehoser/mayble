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
	"fyne.io/fyne/v2/theme"

	"github.com/dubbersthehoser/mayble/internal/viewmodel1"
)

type Fyne struct {
	w fyne.Window
	a fyne.App
}

func NewFyne(a fyne.App, w fyne.Window) *Fyne {
	f := &Fyne{
		w: w,
		a: a,
	}
	return f
}

func NewWindow(f *Fyne, vm *viewmodel.Window) *fyne.Container {
	
	w := f.w

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

	unselect := widget.NewButtonWithIcon("", theme.CancelIcon(), vm.Controls.OnUnselect)
	create := widget.NewButtonWithIcon("", theme.ContentAddIcon(), vm.Controls.OnCreate)
	edit := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), vm.Controls.OnEdit)

	selectedBind := binding.NewString()
	selectedLbl := widget.NewLabelWithData(selectedBind)

	vm.Selected.AddListener(func() {
		row, col := vm.Selected.Get()
		if vm.Selected.Has() {
			data := vm.DataTable.Get(row, col)
			if data != "" {
				data = " | " + data
			}
			format := fmt.Sprintf("%d:%d%s", row, col, data)
			selectedBind.Set(format)
		} else {
			selectedBind.Set("")
		}
	})

	var timer *time.Timer
	duration := time.Second * 2
	final := false

	deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
	deleteBtn.OnTapped = func(){
		if !final {
			final = true
			deleteBtn.SetText("?")
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(duration, func() {
				fyne.Do(func() {
					final = false
					deleteBtn.SetText("")
				})
			})
			return
		}

		timer.Stop()
		timer = nil

		deleteBtn.SetText("")
		final = false
		vm.Controls.OnDelete()
	}

	
	view := container.NewHBox(
		create,
		edit,
		deleteBtn,
		unselect,
		selectedLbl,
	)

	update := func() {
		if vm.Body.Value() != viewmodel.BodyTable {
			view.Hide()
		} else {
			view.Show()
		}
		if vm.Selected.Has() {
			deleteBtn.Enable()
			edit.Enable()
			unselect.Enable()
		} else {
			deleteBtn.Disable()
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
		fyne.Do(func() {
			label.SetText("")
		})
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

	bodyUpdate := func() {
		if vm.Body.Value() != viewmodel.BodyTable {
			table.Items[loanIdx].Disabled = true
			table.Items[readIdx].Disabled = true
		} else {
			table.Items[loanIdx].Disabled = false
			table.Items[readIdx].Disabled = false
		}
		menu.Refresh()
	}

	updateCheck := func() {
		if vm.ColumnSettings.IsLoanHidden() {
			table.Items[loanIdx].Checked = false
		} else {
			table.Items[loanIdx].Checked = true
		}
		if vm.ColumnSettings.IsReadHidden() {
			table.Items[readIdx].Checked = false
		} else {
			table.Items[readIdx].Checked = true
		}
		menu.Refresh()
	}

	table.Items[loanIdx].Action = func() {
		if !vm.ColumnSettings.IsLoanHidden() {
			vm.ColumnSettings.SetLoanHidden(true)
		} else {
			vm.ColumnSettings.SetLoanHidden(false)
		}
		updateCheck()
	}

	table.Items[readIdx].Action = func() {
		if !vm.ColumnSettings.IsReadHidden() {
			vm.ColumnSettings.SetReadHidden(true)
		} else {
			vm.ColumnSettings.SetReadHidden(false)
		}
		updateCheck()
	}

	vm.Body.AddListener(func() {
		bodyUpdate()
	})


	bodyUpdate()
	updateCheck()
	return menu
}
