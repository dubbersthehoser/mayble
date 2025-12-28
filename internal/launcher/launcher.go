package launcher

import (
	"errors"
	"os"
	"path/filepath"
	
	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	_"fyne.io/fyne/v2/widget"
	
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/gui/controller"
	"github.com/dubbersthehoser/mayble/internal/gui/view"
)

const AppName string = "mayble"

func Run() error {
	guiApp := fyneApp.NewWithID(AppName)

	if guiApp.Driver().Device().IsMobile() {
		return errors.New("unsupported platform")
	}
	
	window := guiApp.NewWindow("Mayble")
	window.Resize(fyne.NewSize(900, 600))
	window.CenterOnScreen()
	window.Show()

	err := SetContent(window)
	if err != nil {
		return err
	}
	guiApp.Run()
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

	var (
		DBFileInConfig    string = filepath.Join(cfg.ConfigDir, "mayble.db")
		//DBFileInHome      string = filepath.Join(userHome, "mayble.db")
		DBFileInDocuments string = filepath.Join(userHome, "Documents", "mayble.db")
	)

	// Select the path for default database path from a prioity list, when not 
	// set, or current path is not found.
	//
	prioityList := []string{
		DBFileInDocuments,
		DBFileInConfig,
	}
	if cfg.DBFile == "" ||  !fileExist(cfg.DBFile) {
		for i := range prioityList {
			path := prioityList[i]
			if fileExist(path) {
				return cfg.SetDBFile(path)
			}
		}
	}
	return cfg.SetDBFile(prioityList[0])
}

func SetContent(window fyne.Window) error {
	cfgDir, err := config.GetDefaultDir(AppName)
	if err != nil {
		return err
	}

	cfg, err := config.Load(cfgDir, "config.json")
	if err != nil {
		return err
	}

	if cfg.DBDriver == "" {
		err := cfg.SetDBDriver("sqlite")
		if err != nil {
			return err
		}
	}

	SetDBFile(cfg)

	control, err := controller.New(cfg)
	if err != nil {
		return err
	}
	funkView, err := view.NewFunkView(control, window)
	if err != nil {
		return err
	}

	window.SetContent(funkView.View)
	return nil
}






