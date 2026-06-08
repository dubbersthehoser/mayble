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

func fatalLaunch(err error) {
	// Todo: create window to display fatal error.
	log.Fatal(err)
}

func main() {

	appName := "mayble"

	cfgPath, err := config.GetDefaultConfigFile(appName)
	if err != nil {
		fatalLaunch(err)
		return
	}
	cfg, err := config.Load(cfgPath)
	switch {
	case errors.Is(err, config.ErrIsOldConfig):
		cfg, err = config.Migrate(cfgPath, appName)
		if err != nil {
			fatalLaunch(err)
			return
		}

	case errors.Is(err, os.ErrNotExist):
		cfg, err = config.NewConfigWithDefaults(appName)
		if err != nil {
			fatalLaunch(err)
			return
		}

	case err != nil:
		fatalLaunch(err)
		return
	}
	defer cfg.Save()

	a := app.NewWithID("com.dubbersthehoser." + appName)
	window := a.NewWindow(appName)

	// Todo: add window size to config.
	window.Resize(fyne.NewSize(900, 600))
	window.CenterOnScreen()

	uiVM := viewmodel.NewMainUI(cfg)
	content := view.NewMainUI(window, uiVM)

	window.SetContent(content)
	window.ShowAndRun()
}
