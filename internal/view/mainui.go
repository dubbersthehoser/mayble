package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/dubbersthehoser/mayble/internal/viewmodel"
)

func NewMainUI(w fyne.Window, uiVM *viewmodel.MainUI) *fyne.Container {

	// render unexpected errors. Force Stop.
	if uiVM.HasErrored() {
		o := container.NewVBox()
		for _, s := range uiVM.Errors() {
			label := widget.NewLabel(s)
			label.Importance = widget.WarningImportance
			o.Add(label)
		}
		return o
	}

	menuButton := widget.NewButton("Menu", func() {
		uiVM.SetBody(viewmodel.BodyMenu)
	})
	tableButton := widget.NewButton("Table", func() {
		uiVM.SetBody(viewmodel.BodyData)
	})
	createButton := widget.NewButton("Submit", func() {
		uiVM.SetBody(viewmodel.BodyForm)
	})

	bodySelect := container.NewHBox(
		menuButton,
		tableButton,
		createButton,
	)

	//toolbar := widget.NewToolbar(
	//	widget.NewToolbarAction(
	//		theme.SettingsIcon(),
	//		func() {
	//			uiVM.SetBody(viewmodel.BodyMenu)
	//		},
	//	),
	//	widget.NewToolbarAction(
	//		theme.ListIcon(),
	//		func() {
	//			uiVM.SetBody(viewmodel.BodyData)
	//		},
	//	),
	//	widget.NewToolbarAction(
	//		theme.DocumentIcon(),
	//		func() {
	//			uiVM.SetBody(viewmodel.BodyForm)
	//		},
	//	),
	//)

	//toolbar.Items[0].ToolbarObject().(*widget.Button).Enable()

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
		bodySelect,
		statusLabel,
	)

	menu := NewMenu(w, uiVM.GetMenuVM())
	form := NewBookSubmissionForm(uiVM.GetBookSubmissionForm())

	table := fullBookTable(
		uiVM.GetTableControllersVM(),
		uiVM.GetTableVM(),
	)

	body := container.NewStack(
		menu,
		table,
		form,
	)

	uiVM.OpenedBody.AddListener(binding.NewDataListener(func() {
		open, _ := uiVM.OpenedBody.Get()
		createButton.Enable()
		menuButton.Enable()
		tableButton.Enable()
		switch open {
		case viewmodel.BodyForm:
			createButton.Disable()
			menu.Hide()
			table.Hide()
			form.Show()
			statusLabel.SetText("")
			body.Refresh()
		case viewmodel.BodyMenu:
			menuButton.Disable()
			menu.Show()
			table.Hide()
			form.Hide()
			statusLabel.SetText("")
			body.Refresh()
		case viewmodel.BodyData:
			tableButton.Disable()
			menu.Hide()
			table.Show()
			form.Hide()
			statusLabel.SetText("")
			body.Refresh()
		default:
			panic("opened body was not found")
		}
	}))

	frame := container.NewBorder(header, nil, nil, nil, body)

	return frame
}
