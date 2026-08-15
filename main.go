package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/view"
	"github.com/dubbersthehoser/mayble/internal/viewmodel"
)

func fatalLaunch(w fyne.Window, err error) {
	log.Println(err)
	body := view.NewFatal("Something Went Wrong!", "Failed to launch application.", err.Error())

	w.SetContent(body)
	w.ShowAndRun()
	os.Exit(1)
}

func main() {

	appName := "mayble"
	appID := "com.dubbersthehoser"
	a := app.NewWithID(appID + "." + appName)
	window := a.NewWindow(appName)

	// locate config directory
	cfgRoot, err := config.GetRootDir()
	if err != nil {
		fatalLaunch(window, err)
		return
	}
	configDir := filepath.Join(cfgRoot, appID, appName)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		fatalLaunch(window, err)
		return
	}

	// open config
	cfgFile := filepath.Join(configDir, config.Filename)
	cfg, err := config.Load(cfgFile)
	switch {
	case errors.Is(err, os.ErrNotExist):
		cfg = config.NewConfigWithDefaults(cfgFile)

	case errors.Is(err, config.ErrIsOldConfig):
		cfg, err = config.Migrate(cfgFile)
		if err != nil {
			fatalLaunch(window, err)
			return
		}

	case err != nil:
		fatalLaunch(window, err)
		return
	}
	defer cfg.Save()

	//
	// Setup Window
	//

	if cfg.UI.WindowFullScreen {
		window.FullScreen()
	}
	window.SetMaster()
	if cfg.UI.WindowCenterOnScreen {
		window.CenterOnScreen()
	}

	// Save window size before closing.
	window.SetCloseIntercept(func() {
		size := window.Canvas().Size()
		cfg.UI.WindowWidth = size.Width
		cfg.UI.WindowHeight = size.Height
		cfg.UI.WindowFullScreen = window.FullScreen()
		window.Close()
	})

	// create and show content
	vm := viewmodel.NewWindow(cfg)
	f := view.NewFyne(a, window)
	content := view.NewWindow(f, vm)

	window.SetContent(content)
	window.Resize(fyne.NewSize(cfg.UI.WindowWidth, cfg.UI.WindowHeight))
	window.ShowAndRun()
}
