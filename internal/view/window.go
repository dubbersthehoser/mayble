package view

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/dubbersthehoser/mayble/doc"
	"github.com/dubbersthehoser/mayble/internal/viewmodel"
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

	vm.ShowError.AddListener(func() {
		dialog.ShowError(vm.ShowError.Err, w)
	})

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
	manual := newManual(vm)

	vm.Body.RegisterHandlers(viewmodel.BodyNoData, noData.Hide, noData.Show)
	vm.Body.RegisterHandlers(viewmodel.BodyTable, table.Hide, table.Show)
	vm.Body.RegisterHandlers(viewmodel.BodyBookCreate, create.Hide, create.Show)
	vm.Body.RegisterHandlers(viewmodel.BodyBookEdit, edit.Hide, edit.Show)
	vm.Body.RegisterHandlers(viewmodel.BodyManual, manual.Hide, manual.Show)

	// allow the handlers to be called.
	vm.Body.Set(vm.Body.Value())

	body := container.NewStack(
		noData,
		table,
		edit,
		create,
		manual,
	)

	return body
}

func newManual(w *viewmodel.Window) fyne.CanvasObject {
	rt := widget.NewRichTextFromMarkdown(doc.Manual)
	rt.Wrapping = fyne.TextWrapWord
	rt.Scroll = fyne.ScrollVerticalOnly
	closeBtn := widget.NewButton("Close", w.Body.Back)
	content := container.NewBorder(nil, closeBtn, nil, nil, closeBtn, container.NewVScroll(rt))
	return content
}

func newNoData(vm *viewmodel.Window) fyne.CanvasObject {
	view := widget.NewLabel("")
	view.Wrapping = fyne.TextWrapWord
	view.TextStyle = fyne.TextStyle{
		Bold: false,
	}
	view.SetText(vm.NoData.Message())
	return view
}

func newControls(vm *viewmodel.Window) fyne.CanvasObject {

	unselect := widget.NewButtonWithIcon("", theme.CancelIcon(), vm.Controls.OnUnselect)
	create := widget.NewButtonWithIcon("", theme.ContentAddIcon(), vm.Controls.OnCreate)
	edit := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), vm.Controls.OnEdit)

	selectedBind := binding.NewString()
	selectedLbl := widget.NewLabelWithData(selectedBind)

	vm.Table.Selected.AddListener(func() {
		point := vm.Table.Selected.Get()
		if vm.Table.Selected.Has() {
			data, err := vm.Table.Sheet.Get(point)
			if err != nil {
				log.Println("view display select:", err)
				return
			}
			if data != "" {
				data = " | " + data
			}
			format := fmt.Sprintf("%d:%d%s", point.Row, point.Col, data)
			selectedBind.Set(format)
		} else {
			selectedBind.Set("")
		}
	})

	var timer *time.Timer
	duration := time.Second * 2
	final := false

	deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
	deleteBtn.OnTapped = func() {
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
		if vm.Table.Selected.Has() {
			deleteBtn.Enable()
			edit.Enable()
			unselect.Enable()
		} else {
			deleteBtn.Disable()
			edit.Disable()
			unselect.Disable()
		}
	}

	vm.Table.Selected.AddListener(update)
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

	vm.DoOnClear = func() {
		fyne.Do(func() {
			label.SetText("")
		})
	}

	return label
}

func newMainMenu(vm *viewmodel.Window, w fyne.Window) *fyne.MainMenu {

	// Labels for paticular MenuItems for disabling, and adding check marks.
	const (
		ShowLoanLabel string = "Show Loaned"
		ShowReadLabel string = "Show Read"
		ShowIDLabel   string = "Show ID"

		ImportLabel string = "Import"
		ExportLabel string = "Export"
	)

	// File Menu
	file := fyne.NewMenu("File",
		fyne.NewMenuItem("Current Database", func() {
			db := fmt.Sprintf("\"%s\"", vm.DBPath.Get())
			// quick and straigh forward way to check database path in app.
			dialog.ShowInformation("Current Database", db, w)
		}),
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

		fyne.NewMenuItem(ImportLabel, func() {
			d := dialog.NewFileOpen(
				viewmodel.WrapFyneFileOpen(vm.FileManage.ImportFile),
				w,
			)
			d.Resize(w.Canvas().Size())
			d.SetTitleText("Import CSV")
			d.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
			d.Show()
		}),
		fyne.NewMenuItem(ExportLabel, func() {
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

	// Defined latter.
	var updateCheck func()

	// Table Menu
	table := fyne.NewMenu("Table",
		fyne.NewMenuItem(ShowLoanLabel, func() {
			if !vm.Table.Settings.IsLoanHidden() {
				vm.Table.Settings.SetLoanHidden(true)
			} else {
				vm.Table.Settings.SetLoanHidden(false)
			}
			updateCheck()
		}),
		fyne.NewMenuItem(ShowReadLabel, func() {
			if !vm.Table.Settings.IsReadHidden() {
				vm.Table.Settings.SetReadHidden(true)
			} else {
				vm.Table.Settings.SetReadHidden(false)
			}
			updateCheck()
		}),

		fyne.NewMenuItem(ShowIDLabel, func() {
			if !vm.Table.Settings.IsIDHidden() {
				vm.Table.Settings.SetIDHidden(true)
			} else {
				vm.Table.Settings.SetIDHidden(false)
			}
			updateCheck()
		}),
	)

	// Help Menu
	help := fyne.NewMenu("Help",
		fyne.NewMenuItem("Manual", func() {
			vm.Body.Set(viewmodel.BodyManual)
		}),
	)

	// Create MainMenu
	menu := fyne.NewMainMenu(file, table, help)

	// Setting up the MenuItem locking system for when the Body value changes, and
	// maintains as active if its in Table view.
	menuItemsToLock := make([]*fyne.MenuItem, 0)
	for _, mi := range table.Items {
		switch mi.Label {
		case ShowLoanLabel, ShowReadLabel, ShowIDLabel:
			menuItemsToLock = append(menuItemsToLock, mi)
		}
	}
	for _, mi := range file.Items {
		switch mi.Label {
		case ImportLabel, ExportLabel:
			menuItemsToLock = append(menuItemsToLock, mi)
		}
	}
	bodyUpdate := func() {
		for _, mi := range menuItemsToLock {
			mi.Disabled = vm.Body.Value() != viewmodel.BodyTable
		}
		menu.Refresh()
	}

	// Setting up the Checking system. Mainly used for table Menu Items for when a
	// set of columns are set to hidden by that Item action.
	type check struct {
		mi        *fyne.MenuItem
		predicate func() bool
	}
	menuItemsToCheck := make([]check, 0)
	for _, mi := range table.Items {
		switch mi.Label {
		case ShowLoanLabel:
			ch := check{mi: mi, predicate: func() bool {
				return !vm.Table.Settings.IsLoanHidden()
			}}
			menuItemsToCheck = append(menuItemsToCheck, ch)

		case ShowReadLabel:
			ch := check{mi: mi, predicate: func() bool {
				return !vm.Table.Settings.IsReadHidden()
			}}
			menuItemsToCheck = append(menuItemsToCheck, ch)

		case ShowIDLabel:
			ch := check{mi: mi, predicate: func() bool {
				return !vm.Table.Settings.IsIDHidden()
			}}
			menuItemsToCheck = append(menuItemsToCheck, ch)
		}
	}
	updateCheck = func() {
		for _, ch := range menuItemsToCheck {
			ch.mi.Checked = ch.predicate()
		}
		menu.Refresh()
	}

	vm.Body.AddListener(func() {
		bodyUpdate()
	})

	bodyUpdate()
	updateCheck()

	return menu
}
