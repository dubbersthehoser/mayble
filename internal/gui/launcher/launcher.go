package launcher

import (
	"fmt"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	
	myapp "github.com/dubbersthehoser/mayble/internal/app"
	appStub "github.com/dubbersthehoser/mayble/internal/app/stub"
	storeDriver "github.com/dubbersthehoser/mayble/internal/storage/driver"

	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/gui/controller"
	"github.com/dubbersthehoser/mayble/internal/gui/view"
	"github.com/dubbersthehoser/mayble/internal/settings"
)

func loadConfig(s *settings.Settings) (*config.Config, error) {
	return config.Load(s.ConfigPath)
}

func SetupTextGrid() *widget.TextGrid {
	textgrid := widget.NewTextGrid()
	textgrid.Scroll = fyne.ScrollBoth
	return textgrid
}

func Run(options ...Option) {

	s := settings.Settings{}
	for _, option := range options {
		option(&s)
	}

	defaultConfigDir(&s)
	defaultConfigPath(&s)
	defaultDBDriver(&s)
	defaultDBPath(&s)
	
	App := app.NewWithID("app.dubbersthehoser.mayble")
	window := App.NewWindow("Mayble Launcher")
	window.Resize(fyne.NewSize(800, 500))
	window.CenterOnScreen()

	logGrid := SetupTextGrid()

	logGrid.Append("Hello, World!")
	logGrid.Append("--- PATHS ---")
	logGrid.Append(fmt.Sprintf("config: '%s'", s.ConfigPath))
	logGrid.Append(fmt.Sprintf("storage: '%s'", s.DBPath))
	logGrid.Append("\n--- LOADING ---")

	var Errored bool = false

	_, err := loadConfig(&s)
	if err != nil {
		logGrid.Append(fmt.Sprintf("- config: failed: %s", err.Error()))
		Errored = true
	} else {
		logGrid.Append("- config: success")
	}

	storage, err := storeDriver.Load("memory", s.DBPath)
	if err != nil && !Errored{
		logGrid.Append(fmt.Sprintf("- storage: failed: %s", err.Error()))
		Errored = true
	} else {
		logGrid.Append("- storage: success")
	}

	 _ = storage

	app := myapp.Mayble(&appStub.App{})
	app, err = myapp.New(storage)
	//core, err := core.New(storage)
	//if err != nil && !Errored {
	//	logGrid.Append(fmt.Sprintf("- core failed: %s", err.Error()))
	//	Errored = true
	//} else {
	//	logGrid.Append("- core: success")
	//}

	control := controller.New(app)
	funkView, err := view.NewFunkView(control, window)
	if err != nil {
		logGrid.Append(fmt.Sprintf("- view: failed: %s", err.Error()))
		Errored = true
	}

	mainView := funkView.View

	if Errored  {
		mainView = logGrid
	}

	window.SetContent(mainView)
	window.ShowAndRun()
}




