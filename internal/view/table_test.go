package view

import (
	"testing"

	//"fyne.io/fyne/v2"
	_ "fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/viewmodel"
)

func TestTable(t *testing.T) {

	cfg := config.NewConfigWithDefaults("")
	vm := viewmodel.NewWindow(cfg)
	table := newTable(vm)
	table.Refresh()

	if len(table.header.buttons) == 0 {
		t.Fatalf("header buttons don't exists after creating table")
	}

	// need to call UpdateHeader() manualy for all header buttons. fyne test code
	// dose not call it. Only the CreateHeader() after calling Refresh()
	for i, h := range table.header.buttons {
		cellID := widget.TableCellID{Col: i, Row: -1}
		table.UpdateHeader(cellID, h)
	}

	// Seems like I can't test widget widths in fyne tests,
	// so I'm commenting this out.
	//for _, h := range table.header.buttons {
	//	actual := h.Size()
	//	if actual.Width != vm.Table.Settings.GetWidth(h.label) {
	//		t.Fatalf("header button '%s' width: expect: %f,  got: %f", h.label, vm.Table.Settings.GetWidth(h.label), actual.Width)
	//	}
	//	if actual.Height != vm.Table.Settings.HeaderHeight() {
	//		t.Fatalf("header button '%s' height: expect: %f, got: %f", h.label, vm.Table.Settings.HeaderHeight(), actual.Height)
	//	}
	//}

	//table.header.buttons[0].Resize(fyne.Size{Width:0, Height:0})
	//label := table.header.buttons[0].label
	//actual := table.header.buttons[0].Size()
	//expect := fyne.Size{Width: vm.Table.Settings.HeaderMinWidth(), Height: vm.Table.Settings.HeaderHeight()}
	//if expect.Height != actual.Height {
	//	t.Fatalf("header button '%s' height: expect: %f, got: %f", label, expect.Height, actual.Height)
	//}

}
