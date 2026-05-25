package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/view"
	"github.com/dubbersthehoser/mayble/internal/viewmodel"
)

func main() {

	appName := "mayble"
	a := app.NewWithID("com.dubbersthehoser.mayble")
	window := a.NewWindow("Mayble")
	window.Resize(fyne.NewSize(900, 600))
	window.CenterOnScreen()

	cfgFile, err := config.GetDefaultConfigFile(appName)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.Load(cfgFile, appName)
	if err != nil {
		log.Fatal(err)
	}

	if cfg != nil {
		defer cfg.Save()
	}

	uiVM := viewmodel.NewMainUI(cfg)
	content := view.NewMainUI(window, uiVM)
	window.SetContent(content)
	window.ShowAndRun()
}
