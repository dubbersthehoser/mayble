package view

import (
	"log"
	"fmt"
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

	menu := NewMenu(w, uiVM.GetMenu())
	form := NewBookSubmissionForm(uiVM.GetBookSubmissionForm())
	table := fullBookTable(uiVM.GetTable())
	setBodyShow := func(bodyID int) {
		switch bodyID {
		case viewmodel.BodyMenu:
			menu.Show()
			form.Hide()
			table.Hide()
		case viewmodel.BodyData:
			menu.Hide()
			form.Hide()
			table.Show()
		case viewmodel.BodyForm:
			menu.Hide()
			form.Show()
			table.Hide()
		default:
			log.Println(fmt.Sprintf("ERROR: invalid body id %d", bodyID))
		}
	}

	uiVM.HasDatabase.AddListener(binding.NewDataListener(func() {
		ok, _ := uiVM.HasDatabase.Get()
		if ok {
			form.Show()
			table.Show()
		}
	}))

	body := container.NewStack(
		menu,
		table,
		form,
	)

	switcher := uiVM.GetBodySwitcher()
	switcher.SetBodies(
		*viewmodel.NewBodyButton(
			"Menu",
			viewmodel.BodyMenu,
			&viewmodel.BodyWindow{
				OnHide: func() {
					menu.Hide()
					body.Refresh()
				},
				OnShow: func() {
					menu.Show()
					body.Refresh()
				},
			},
		),
		*viewmodel.NewBodyButton(
			"Table",
			viewmodel.BodyData,
			&viewmodel.BodyWindow{
				OnHide: func() {
					table.Hide()
					body.Refresh()
				},
				OnShow: func() {
					table.Show()
					body.Refresh()
				},
			},
		),
		*viewmodel.NewBodyButton(
			"Submit",
			viewmodel.BodyForm,
			&viewmodel.BodyWindow{
				OnHide: func() {
					table.Hide()
					body.Refresh()
				},
				OnShow: func() {
					table.Show()
					body.Refresh()
				},
			},
		),
	)

	bodyButtons := map[int]fyne.CanvasObject{
		viewmodel.BodyMenu: widget.NewButton("Menu", func() {
			switcher.Switch(viewmodel.BodyMenu)
		}),
		viewmodel.BodyData: widget.NewButton("Table", func() {
			switcher.Switch(viewmodel.BodyData)
		}),
		viewmodel.BodyForm: widget.NewButton("Submit", func() {
			switcher.Switch(viewmodel.BodyForm)
		}),
	}

	switcher.Buttons[viewmodel.BodyMenu].Locked.AddListener(binding.NewDataListener(func(){
		
	}))

	for _, btn := range switcher.Buttons() {
		w := widget.NewButton(btn.Label(), func() {
			switcher.Switch(btn.ID())
		})
	

		btn.Active.AddListener(binding.NewDataListener(func() {
			active, _ := btn.Active.Get()
			if active {
				w.Enable()
				setBodyShow(btn.ID())
				body.Refresh()
			} else {
				w.Disable()
				body.Refresh()
			}
			
		}))
		bodyButtons = append(bodyButtons, w)
	}


	bodySelect := container.NewHBox(bodyButtons...)

	header := container.NewHBox(
		bodySelect,
		statusLabel,
	)

	frame := container.NewBorder(header, nil, nil, nil, body)

	return frame
}
