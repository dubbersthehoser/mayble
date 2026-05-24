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

	body := container.NewStack(
		menu,
		table,
		form,
	)

	bodyButtons := map[int]*widget.Button{
		viewmodel.BodyMenu: widget.NewButton("Menu", nil),
		viewmodel.BodyData: widget.NewButton("Table", nil),
		viewmodel.BodyForm: widget.NewButton("Submit", nil),
	}


	switcher := uiVM.GetBodySwitcher()
	switcher.SetBodies(
		viewmodel.BodyButton{
			ID: viewmodel.BodyMenu,
			OnLock: bodyButtons[viewmodel.BodyMenu].Disable,
			OnUnlock: bodyButtons[viewmodel.BodyMenu].Enable,

			Window: viewmodel.BodyWindow{
				OnHide: menu.Hide, // !I may need to refresh
				OnShow: menu.Show,
			},
		},
		viewmodel.BodyButton{
			ID: viewmodel.BodyData,
			OnLock: bodyButtons[viewmodel.BodyData].Disable,
			OnUnlock: bodyButtons[viewmodel.BodyData].Enable,

			Window: viewmodel.BodyWindow{
				OnHide: table.Hide,
				OnShow: table.Show,
			},
		},
		viewmodel.BodyButton{
			ID: viewmodel.BodyForm,
			OnLock: bodyButtons[viewmodel.BodyForm].Disable,
			OnUnlock: bodyButtons[viewmodel.BodyForm].Enable,

			Window: viewmodel.BodyWindow{
				OnHide: form.Hide,
				OnShow: form.Show,
			},
		},
	)

	bodySelect := container.NewHBox()
	for _, o := range bodyButtons {
		bodySelect.Objects = append(bodySelect.Objects, o)
	}

	header := container.NewHBox(
		bodySelect,
		statusLabel,
	)

	frame := container.NewBorder(header, nil, nil, nil, body)

	return frame
}
