package view

import (
	"fyne.io/fyne/v2"

	"github.com/dubbersthehoser/mayble/internal/viewmodel"
)

func newMainMenu(vm *viewmodel.Window, eb *EventBus, cb *CommandBus) *fyne.MainMenu {

	file := newFileMenu(eb, cb)

	// Table Menu
	table := newTableMenu(eb, cb)

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

func newFileMenu(eb *EventBus, cb *CommandBus) *fyne.Menu {

	const (
		ImportLabel           string = "Import"
		ExportLabel           string = "Export"
		ShowDatabasePathLabel string = "Current Database"
		CreateDatabaseLabel   string = "Create"
		OpenDatabaseLabel     string = "Open"
	)

	labelToDialogName := func(label string) string {
		var name string
		switch label {
		case OpenDatabaseLabel:
			name = "OpenDatabase"
		case CreateDatabaseLabel:
			name = "CreateDatabase"
		case ExportLabel:
			name = "ExportCSV"
		case ImportLabel:
			name = "ImportCSV"
		case ShowDatabasePathLabel:
			name = "ShowDatabasePath"
		}
		return name
	}

	file := fyne.NewMenu("File",
		fyne.NewMenuItem(ShowDatabasePathLabel, func() {
			cb.Dispatch(OpenDialog{Name: labelToDialogName(ShowDatabasePathLabel)})
		}),
		fyne.NewMenuItem(OpenDatabaseLabel, func() {
			cb.Dispatch(OpenDialog{Name: labelToDialogName(OpenDatabaseLabel)})
		}),
		fyne.NewMenuItem(CreateDatabaseLabel, func() {
			cb.Dispatch(OpenDialog{Name: labelToDialogName(CreateDatabaseLabel)})
		}),
		fyne.NewMenuItem(ImportLabel, func() {
			cb.Dispatch(OpenDialog{Name: labelToDialogName(ImportLabel)})
		}),
		fyne.NewMenuItem(ExportLabel, func() {
			cb.Dispatch(OpenDialog{Name: labelToDialogName(ExportLabel)})
		}),
	)

	eb.Subscribe("BodyChanged", func(v any) {
		args := v.(BodyChanged)
		for _, item := range file.Items {
			item.Disabled = !canOpenDialogWithBody(labelToDialogName(item.Label), args.Body)
		}
		file.Refresh()
	})

	return file
}

func newTableMenu(ev *EventBus, cb *CommandBus) *fyne.Menu {

	var (
		ShowLoanLabel string = "Show Loaned"
		ShowReadLabel string = "Show Read"
		ShowIDLabel   string = "Show ID"
	)

	labelToColumnName := func(s string) string {
		var name string
		switch s {
		case ShowLoanLabel:
			name = "LoanSet"
		case ShowReadLabel:
			name = "ReadSet"
		case ShowIDLabel:
			name = "ID"
		}
		return name
	}

	table := fyne.NewMenu("Table",
		fyne.NewMenuItem(ShowLoanLabel, func() {
			cb.Dispatch(ToggleHiddenColumn{Column: labelToColumnName(ShowLoanLabel)})
		}),
		fyne.NewMenuItem(ShowReadLabel, func() {
			cb.Dispatch(ToggleHiddenColumn{Column: labelToColumnName(ShowReadLabel)})
		}),
		fyne.NewMenuItem(ShowIDLabel, func() {
			cb.Dispatch(ToggleHiddenColumn{Column: labelToColumnName(ShowIDLabel)})
		}),
	)

	ev.Subscribe("ColumnHiddenChanged", func(v any) {
		args := v.(ColumnHiddenChanged)
		for _, item := range table.Items {
			if item.Label == ShowLoanLabel {
				item.Checked = !args.LoanSet
			}
			if item.Label == ShowReadLabel {
				item.Checked = !args.ReadSet
			}
			if item.Label == ShowIDLabel {
				item.Checked = !args.ID
			}
		}
		table.Refresh()
	})

	ev.Subscribe("BodyChanged", func(v any) {
		args := v.(BodyChanged)
		for _, item := range table.Items {
			item.Disabled = args.Body != viewmodel.BodyTable
		}
		table.Refresh()
	})
	return table
}
