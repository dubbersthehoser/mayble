package main

import (
	"errors"
	"log"
	"os"
	"fmt"

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

	if cfg.UI.WindowFullScreen {
		window.FullScreen()
	}

	window.SetMaster()
	if cfg.UI.WindowCenterOnScreen {
		window.CenterOnScreen()
	}

	window.SetCloseIntercept(func() {
		size := window.Canvas().Size()
		fmt.Printf("debug: get: width=%f, height=%f\n", size.Width, size.Height)
		cfg.UI.WindowWidth = size.Width
		cfg.UI.WindowHeight = size.Height
		cfg.UI.WindowFullScreen = window.FullScreen()
		window.Close()
	})

	// create and show content
	vm := viewmodel.NewWindow(cfg)
	f := view.NewFyne(a, window)
	content := view.NewWindow(f, vm)

	fmt.Printf("debug: set: width=%f, height=%f\n", cfg.UI.WindowWidth, cfg.UI.WindowHeight)
	window.SetContent(content)
	window.Resize(fyne.NewSize(cfg.UI.WindowWidth, cfg.UI.WindowHeight))
	fmt.Printf("debug: get first: width=%f, height=%f\n", window.Content().Size().Width, window.Content().Size().Height)
	window.ShowAndRun()
}

