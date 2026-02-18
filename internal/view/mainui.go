package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"

	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/viewmodel"
)

func NewMainUI(w fyne.Window, uiVM *viewmodel.MainUI) *fyne.Container {

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(
			theme.SettingsIcon(),
			func() {
				_ = uiVM.OpenedBody.Set(viewmodel.BodyMenu)
			},
		),
		widget.NewToolbarAction(
			theme.ListIcon(),
			func() {
				_ = uiVM.OpenedBody.Set(viewmodel.BodyData)
			},
		),
		widget.NewToolbarAction(
			theme.DocumentIcon(),
			func() {
				_ = uiVM.OpenedBody.Set(viewmodel.BodyForm)
			},
		),
	)

	toolbar.Items[0].ToolbarObject().(*widget.Button).Enable()
	menuButton := toolbar.Items[0].ToolbarObject().(*widget.Button)
	tablesButton := toolbar.Items[1].ToolbarObject().(*widget.Button)
	addButton := toolbar.Items[2].ToolbarObject().(*widget.Button)

	// Status Line 
	// Displays input form .Error, .Info, .Success string bindings with proper colors.
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
	uiVM.Clear.AddListener(binding.NewDataListener(func() {
		ok, _ := uiVM.Clear.Get()
		if ok {
			statusLabel.SetText("")
		}
	}))

	header := container.NewHBox(
		toolbar,
		statusLabel,
	)

	form := NewCreateBookForm(uiVM.GetCreateBookFormVM())
	table := BookTables(uiVM.GetTablesVM())
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
			statusLabel.SetText("")
			body.Refresh()
		case viewmodel.BodyMenu:
			menuButton.Disable()
			body.Objects[0] = widget.NewLabel("not implemented")
			statusLabel.SetText("")
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


