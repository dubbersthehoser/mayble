package launcher

import (
	"errors"
	"path/filepath"
	
	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	_"fyne.io/fyne/v2/widget"
	
	myapp "github.com/dubbersthehoser/mayble/internal/app"
	storeDriver "github.com/dubbersthehoser/mayble/internal/storage/driver"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/gui/controller"
	"github.com/dubbersthehoser/mayble/internal/gui/view"
)


func Run() error {

	AppName := "mayble"

	// Load Config

	dir, err := config.GetDefaultDir(AppName)
	if err != nil {
		return err
	}

	cfg, err := config.Load(dir, "config.json")
	if err != nil {
		return err
	}

	guiApp := fyneApp.NewWithID(AppName)

	if guiApp.Driver().Device().IsMobile() {
		return errors.New("unsupported platform")
	}

	if cfg.DBDriver == "" {
		err := cfg.SetDBDriver("sqlite")
		if err != nil {
			return err
		}
	}

	if cfg.DBFile == "" {
		err := cfg.SetDBFile(filepath.Join(dir, "mayble.db"))
		if err != nil {
			return err
		}
	}

	storage, err := storeDriver.Load(cfg.DBDriver, cfg.DBFile)
	if err != nil {
		return err	
	}

	coreApp, err := myapp.New(storage)
	if err != nil {
		return err
	}

	window := guiApp.NewWindow("Mayble")
	window.Resize(fyne.NewSize(800, 500))
	window.CenterOnScreen()

	control := controller.New(coreApp, cfg)
	funkView, err := view.NewFunkView(control, window)
	if err != nil {
		return err
	}

	mainView := funkView.View

	window.SetContent(mainView)
	window.ShowAndRun()
	return nil
}

func ShowApp() {
}






