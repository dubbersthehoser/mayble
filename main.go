package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/database"
	"github.com/dubbersthehoser/mayble/internal/view"
	"github.com/dubbersthehoser/mayble/internal/viewmodel"
	service "github.com/dubbersthehoser/mayble/internal/app"
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

	as := service.NewService(cfg)

	var db *database.Database

	if cfg.DBFile == "" {
		db, err = database.OpenMem()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		db, err = database.Open(cfg.DBFile)
		if err != nil {
			log.Fatal(err)
		}
	}
	if db != nil {
		defer db.Conn.Close()
	}


	uiVM := viewmodel.NewMainUI(cfg, db)
	content := view.NewMainUI(window, uiVM)
	window.SetContent(content)
	window.ShowAndRun()
}
