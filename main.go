package main

import (
	"os"
	"log"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/view1"
	"github.com/dubbersthehoser/mayble/internal/viewmodel1"
)

func fatalLaunch(w fyne.Window, err error) {
	// Todo: create window to display fatal error.
	log.Fatal(err)
	body := view.NewFatal("Fatal", "Failed to launch application.", err.Error())

	w.SetContent(body)
}

func main() {

	appName := "mayble"
	a := app.NewWithID("com.dubbersthehoser." + appName)
	window := a.NewWindow(appName)

	cfgPath, err := config.GetDefaultConfigFile(appName)
	if err != nil {
		fatalLaunch(window, err)
		return
	}
	cfg, err := config.Load(cfgPath)
	switch {
	case errors.Is(err, config.ErrIsOldConfig):
		cfg, err = config.Migrate(cfgPath, appName)
		if err != nil {
			fatalLaunch(window, err)
			return
		}

	case errors.Is(err, os.ErrNotExist):
		cfg, err = config.NewConfigWithDefaults(appName)
		if err != nil {
			fatalLaunch(window, err)
			return
		}

	case err != nil:
		fatalLaunch(window, err)
		return
	}
	defer cfg.Save()


	// Todo: add window size to config.
	window.Resize(fyne.NewSize(900, 600))
	window.CenterOnScreen()
	window.SetMaster()

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			// This code maybe bad for fyne. Needs testing.
			// And the err.(string) most likely an bad idea.
			body := view.NewFatal("Crash", "Unexpected crash.", err.(string))
			window.SetContent(body)
		}
	}()

	vm := viewmodel.NewWindow(cfg)
	f := view.NewFyne(a, window)
	content := view.NewWindow(f, vm)

	window.SetContent(content)
	window.ShowAndRun()
}
