package main

import (
	"os"
	"log"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/view"
	"github.com/dubbersthehoser/mayble/internal/viewmodel"
)

func fatalLaunch(w fyne.Window, err error) {
	log.Fatal(err)
	body := view.NewFatal("Fatal", "Failed to launch application.", err.Error())

	w.SetContent(body)
	w.ShowAndRun()
	os.Exit(1)
}

func main() {

	appName := "mayble"
	a := app.NewWithID("com.dubbersthehoser." + appName)
	window := a.NewWindow(appName)

	// open config
	cfgPath, err := config.GetDefaultConfigFile(appName)
	if err != nil {
		fatalLaunch(window, err)
		return
	}
	cfg, err := config.Load(cfgPath)
	switch {
	case errors.Is(err, config.ErrIsOldConfig):
		cfg, err = config.Migrate(cfgPath)
		if err != nil {
			fatalLaunch(window, err)
			return
		}

	case errors.Is(err, os.ErrNotExist):
		configFile, err := config.GetDefaultConfigFile(appName)
		cfg = config.NewConfigWithDefaults(configFile)
		if err != nil {
			fatalLaunch(window, err)
			return
		}

	case err != nil:
		fatalLaunch(window, err)
		return
	}
	defer cfg.Save()


	// window set up
	window.Resize(fyne.NewSize(cfg.UI.WindowWidth, cfg.UI.WindowHeight))
	window.SetMaster()
	if cfg.UI.WindowCenterOnScreen {
		window.CenterOnScreen()
	}
	defer func() {
		size := window.Content().Size()
		cfg.UI.WindowWidth = size.Width
		cfg.UI.WindowHeight = size.Height
		cfg.UI.WindowFullScreen = window.FullScreen()
	}()

	// create and show content
	vm := viewmodel.NewWindow(cfg)
	f := view.NewFyne(a, window)
	content := view.NewWindow(f, vm)

	window.SetContent(content)
	window.ShowAndRun()
}
