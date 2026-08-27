package view

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
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

	cb := NewCommandBus()
	eb := NewEventBus()

	setupCommandBus(cb, vm, f)
	setupEventBus(eb, vm)

	w.SetMainMenu(newMainMenu(vm, eb, cb))

	status := newStatusLine(vm.StatusLine)
	controls := newControls(vm)

	vm.ShowError.AddListener(func() {
		dialog.ShowError(vm.ShowError.Err, w)
	})

	topBar := container.NewGridWithColumns(2, status, controls)
	body := newBody(vm)

	view := container.NewBorder(topBar, nil, nil, nil, body)

	vm.Table.NotifyOnColumnHidden()

	return view
}

func setupCommandBus(cb *CommandBus, vm *viewmodel.Window, f *Fyne) {
	cb.Regester("OpenDialog", func(v any) error {
		return restrictsOpenDialog(v.(OpenDialog), vm, f.w)
	})
	cb.Regester("ToggleHiddenColumn", func(v any) error {
		return commandToggleHiddenColumn(v.(ToggleHiddenColumn), vm)
	})
}

func setupEventBus(eb *EventBus, vm *viewmodel.Window) {
	vm.Body.AddListener(func() {
		eb.Notify(BodyChanged{Body: vm.Body.Value()})
	})

	vm.Table.Settings.AddOnHidden(func() {
		eb.Notify(ColumnHiddenChanged{
			ID:      vm.Table.Settings.IsIDHidden(),
			LoanSet: vm.Table.Settings.IsLoanHidden(),
			ReadSet: vm.Table.Settings.IsReadHidden(),
		})
	})
}

func commandToggleHiddenColumn(args ToggleHiddenColumn, vm *viewmodel.Window) error {
	switch args.Column {
	case "LoanSet":
		_ = vm.Table.Settings.ToggleHiddenLoan()
	case "ReadSet":
		_ = vm.Table.Settings.ToggleHiddenRead()
	case "ID":
		_ = vm.Table.Settings.ToggleHiddenID()
	default:
		log.Printf("toggle_hidden_column %s: invalid column", args.Column)
	}
	return nil
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
			data, err := vm.Table.GetDataPoint(point)
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


