package launcher

import (
	"errors"
	"os"
	"path/filepath"
	
	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/gui/controller"
	"github.com/dubbersthehoser/mayble/internal/gui/view"
)

const AppName string = "mayble"

func Run() error {

	guiApp := fyneApp.NewWithID(AppName)

	window := guiApp.NewWindow("Mayble")
	window.Resize(fyne.NewSize(900, 600))
	window.CenterOnScreen()
	window.Show()

	err := SetContent(window, guiApp.Driver())
	if err != nil {
		return err
	}
	guiApp.Run()
	return nil
}


func SetContent(window fyne.Window, drv fyne.Driver) error {

	fatalView := func(err error) fyne.CanvasObject {
		v := widget.NewLabel(err.Error())
		return v
	}

	cfgDir, err := config.GetDefaultDir(AppName)
	if err != nil {
		window.SetContent(fatalView(err))
		return nil
	}



	cfg, err := config.Load(cfgDir, "config.json")
	if err != nil {
		window.SetContent(fatalView(err))
		return nil
	}

	if cfg.DBDriver == "" {
		err := cfg.SetDBDriver("sqlite")
		if err != nil {
			window.SetContent(fatalView(err))
			return nil
		}
	}

	if err := SetDBFile(cfg); err != nil {
		window.SetContent(fatalView(err))
		return nil
	}

	control, err := controller.New(cfg)
	if err != nil {
		window.SetContent(fatalView(err))
		return nil
	}
	funkView, err := view.NewFunkView(control, window)
	if err != nil {
		window.SetContent(fatalView(err))
		return nil
	}

	window.SetContent(funkView.View)
	return nil
}


func SetDBFile(cfg *config.Config) error {

	var err error

	fileExist := func(path string) bool {
		_, err = os.Stat(path)
		return !errors.Is(err, os.ErrNotExist)
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	if cfg.DBFile != "" || fileExist(cfg.DBFile) {
		return nil
	}

	
	// When there is no config file path, or the file path does not exists, then we'll set a default path.
	// The default path is set in a priorityList. The top of item in the list has the highest priority.
	// When checking the items in the priority list, if there is a database file at that path it will be
	// selected. If all database paths do not exists, then select the first one with an existing directory.
	//
	var (
		DBFileInConfig    string = filepath.Join(cfg.ConfigDir, "mayble.db")
		DBFileInHome      string = filepath.Join(userHome, "mayble.db")
		DBFileInDocuments string = filepath.Join(userHome, "Documents", "mayble.db")
	)

	prioityList := []string{
		DBFileInDocuments,
		DBFileInHome,
		DBFileInConfig,
	}

	for i := range prioityList {
		path := prioityList[i]
		if fileExist(path) {
			return cfg.SetDBFile(path)
		}
	}
	for i := range prioityList {
		path := prioityList[i]
		dir := filepath.Dir(path)
		if fileExist(dir) {
			return cfg.SetDBFile(path)
		}
	}
	return nil
}






