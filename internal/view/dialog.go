package view

import (
	"fmt"
	"log"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"

	"github.com/dubbersthehoser/mayble/internal/viewmodel"
)

func commandOpenDialog(args OpenDialog, vm *viewmodel.Window, w fyne.Window) error {
	switch args.Name {
	case "OpenDatabase":
		dialogOpenDatabase(vm, w)
	case "CreateDatabase":
		dialogCreateDatabase(vm, w)
	case "ShowDatabasePath":
		dialogShowDatabasePath(vm, w)
	case "ExportCSV":
		dialogExportCSV(vm, w)
	case "ImportCSV":
		dialogImportCSV(vm, w)
	default:
		log.Printf("open_dialog %s: invalid command", args.Name)
	}
	return nil
}

func restrictsOpenDialog(args OpenDialog, vm *viewmodel.Window, w fyne.Window) error {
	if canOpenDialogWithBody(args.Name, vm.Body.Value()) {
		return commandOpenDialog(args, vm, w)
	} else {
		return nil
	}
}

func canOpenDialogWithBody(name string, body int) bool {
	if body == viewmodel.BodyTable {
		return true
	}
	switch name {
	case "OpenDatabase", "CreateDatabase":
		return slices.Contains(
			[]int{viewmodel.BodyTable, viewmodel.BodyNoData},
			body,
		)
	case "ExportCSV", "ImportCSV":
		return slices.Contains(
			[]int{viewmodel.BodyTable},
			body,
		)
	default:
		return true
	}
}

func dialogShowDatabasePath(vm *viewmodel.Window, w fyne.Window) {
	db := fmt.Sprintf("\"%s\"", vm.DBPath.Get())
	dialog.ShowInformation("Current Database", db, w)
}

func dialogOpenDatabase(vm *viewmodel.Window, w fyne.Window) {
	d := dialog.NewFileOpen(
		viewmodel.WrapFyneFileOpen(vm.FileManage.OpenDatabase),
		w,
	)
	d.Resize(w.Canvas().Size())
	d.SetTitleText("Open Database")
	d.SetFilter(storage.NewExtensionFileFilter([]string{".db", ".sqlite", ".sqlite3"}))
	d.Show()
}

func dialogCreateDatabase(vm *viewmodel.Window, w fyne.Window) {
	d := dialog.NewFileSave(
		viewmodel.WrapFyneFileCreate(vm.FileManage.CreateDatabase),
		w,
	)
	d.Resize(w.Canvas().Size())
	d.SetTitleText("Create Database")
	d.SetFilter(storage.NewExtensionFileFilter([]string{".db", ".sqlite", ".sqlite3"}))
	d.Show()
}

func dialogImportCSV(vm *viewmodel.Window, w fyne.Window) {
	d := dialog.NewFileOpen(
		viewmodel.WrapFyneFileOpen(vm.FileManage.ImportFile),
		w,
	)
	d.Resize(w.Canvas().Size())
	d.SetTitleText("Import CSV")
	d.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
	d.Show()
}

func dialogExportCSV(vm *viewmodel.Window, w fyne.Window) {
	d := dialog.NewFileSave(
		viewmodel.WrapFyneFileCreate(vm.FileManage.ExportFile),
		w,
	)
	d.Resize(w.Canvas().Size())
	d.SetTitleText("Export CSV")
	d.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
	d.Show()
}
