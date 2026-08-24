package view

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"

	"github.com/dubbersthehoser/mayble/internal/viewmodel"
)

func newMainMenu(vm *viewmodel.Window, w fyne.Window) *fyne.MainMenu {

	file := newFileMenu(vm, w)

	// Table Menu
	table := newTableMenu(vm)

	// Help Menu
	help := fyne.NewMenu("Help",
		fyne.NewMenuItem("Manual", func() {
			vm.Body.Set(viewmodel.BodyManual)
		}),
	)

	// Create MainMenu
	menu := fyne.NewMainMenu(file, table, help)
	return menu
}

func newFileMenu(vm *viewmodel.Window, w fyne.Window) *fyne.Menu {

	const (
		ImportLabel    string = "Import"
		ExportLabel    string = "Export"
		CurrentDBLabel string = "Current Database"
	)

	file := fyne.NewMenu("File",
		fyne.NewMenuItem(CurrentDBLabel, func() {
			db := fmt.Sprintf("\"%s\"", vm.DBPath.Get())
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

	updateToBody := func() {
		for _, item := range file.Items {
			if item.Label != CurrentDBLabel {
				item.Disabled = vm.Body.Value() != viewmodel.BodyTable
			}
		}
		file.Refresh()
	}

	vm.Body.AddListener(updateToBody)
	return  file 
	
}

func newTableMenu(vm *viewmodel.Window) *fyne.Menu {
	var (
		ShowLoanLabel string = "Show Loaned"
		ShowReadLabel string = "Show Read"
		ShowIDLabel   string = "Show ID"
	)
	// Defined latter. Updates the checkmarks for hidden columns.
	var updateCheck func()

	table := fyne.NewMenu("Table",
		fyne.NewMenuItem(ShowLoanLabel, func() {
			// cmdb.dispatch(tableShowLoaned)
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
		table.Refresh()
	}

	updateToBody := func() {
		for _, item := range table.Items {
			item.Disabled = vm.Body.Value() != viewmodel.BodyTable
		}
		table.Refresh()
	}

	vm.Body.AddListener(updateToBody)

	updateCheck()
	return table
}

