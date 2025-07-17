package gui

import (
	//"log"
	//"image/color"
	//"log"
	"fmt"
	//"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	_"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	_"fyne.io/fyne/v2/data/binding"
	_"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/driver/desktop"
	_"fyne.io/fyne/v2/dialog"
	_"fyne.io/fyne/v2/canvas"

	myapp "github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/event"
)

func GetRattingStrings() []string {
	return []string{"TBR", "⭐", "⭐⭐", "⭐⭐⭐", "⭐⭐⭐⭐", "⭐⭐⭐⭐⭐"}
}



type BookLibary struct {
	BookSelected  int
	BookOrderedBy string
	BookList      []string // TODO add book data
	UniqueGenres  []string
	UniqueAuthors []string
}

type UIState struct {
	BookLibary
	Window         fyne.Window
	DataHasChanged bool
	Emiter         *event.EventEmiter
}
func NewUIState(window fyne.Window) *UIState {
	u := &UIState{
		Emiter: event.NewEventEmiter(),
		Window: window,
	}
	u.InitEvents()
	return u
}


const (
	SaveButtonClicked    string = "SaveButtonClicked"
	ChangedOrderBy              = "ChangedOrderBy"
	ChangedOrderByAsc           = "ChangedOrderByAcs"
	ChangedOrderByDesc          = "ChangedOrderByDesc"
	ChangedSearchBy             = "ChangedSearchBy"
	ChangedSearch               = "ChangedSearch"
	NewBookEvent                = "NewBookEvent"
	NewOnLoanEvent              = "NewOnLoanEvent"
	UpdateBookEvent             = "UpdateBookEvent"
)

func Run() {

	a := app.New()
	window := a.NewWindow(myapp.AppName)

	window.Resize(fyne.NewSize(800, 500))

	ctrlN := &desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: fyne.KeyModifierControl}
	//ctrlS := &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}
	//ctrlD := &desktop.CustomShortcut{KeyName: fyne.KeyD, Modifier: fyne.KeyModifierControl}
	//ctrlE := &desktop.CustomShortcut{KeyName: fyne.KeyE, Modifier: fyne.KeyModifierControl}

	window.Canvas().AddShortcut(ctrlN, func (shortcut fyne.Shortcut){
		fmt.Println("Control+N: Create New book")
	})

	UI := NewUIState(window)


	header := UI.NewHeaderComp()
	body := UI.NewBookTableComp()
	mainComp := container.New(layout.NewBorderLayout(header, nil, nil, nil), header, body)


	window.SetContent(mainComp)


	window.ShowAndRun()
}



















