package launcher

import (
	"fmt"
	
	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	
	myapp "github.com/dubbersthehoser/mayble/internal/app"
	appStub "github.com/dubbersthehoser/mayble/internal/app/stub"
	storeDriver "github.com/dubbersthehoser/mayble/internal/storage/driver"

	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/gui/controller"
	"github.com/dubbersthehoser/mayble/internal/gui/view"
	"github.com/dubbersthehoser/mayble/internal/config"
)

func loadConfig(s *settings.Settings) (*config.Config, error) {
	return config.Load(s.ConfigPath)
}

func SetupTextGrid() *widget.TextGrid {
	textgrid := widget.NewTextGrid()
	textgrid.Scroll = fyne.ScrollBoth
	return textgrid
}



func Run() error {

	guiApp := fyneApp.NewWithID("app.dubbersthehoser.mayble")
	
	splash := guiApp.NewWindow("Loading...")


	cfgFile := "config.json"

	dir, err := config.GetDefaultDir()
	if err != nil {
		return err
	}

	cfg, err := config.Load(dir, "config.json")
	if err != nil {
		return err
	}
	
	
	window := guiApp.NewWindow("Mayble")
	window.Resize(fyne.NewSize(800, 500))
	window.CenterOnScreen()

	logGrid := SetupTextGrid()

	logGrid.Append("Hello, World!")
	logGrid.Append("--- PATHS ---")
	logGrid.Append(fmt.Sprintf("config: '%s'", cfg.ConfigPath))
	logGrid.Append(fmt.Sprintf("storage: '%s'", cfg.DBPath))
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

	coreApp := myapp.API(&appStub.App{})
	coreApp, err = myapp.New(storage)
	if err != nil {
		logGrid.Append(fmt.Sprintf("- app: failed: %s", err))
	}

	control := controller.New(coreApp, s)
	funkView, err := view.NewFunkView(control, window)
	if err != nil {
		logGrid.Append(fmt.Sprintf("- view: failed: %s", err))
		Errored = true
	}

	mainView := funkView.View

	if Errored  {
		mainView = logGrid
	}

	window.SetContent(mainView)
	window.ShowAndRun()
	return nil
}





