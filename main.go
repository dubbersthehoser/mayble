package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/dubbersthehoser/mayble/internal/view"
	"github.com/dubbersthehoser/mayble/internal/viewmodel"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/database"
)


func main() {
	appName := "mayble"
	a := app.NewWithID("com.dubbersthehoser.mayble")
	window := a.NewWindow("Mayble")
	window.Resize(fyne.NewSize(900, 600))
	window.CenterOnScreen()

	errList := make([]error, 0)

	configDir, err := config.GetDefaultDir(appName)
	if err != nil {
		log.Println(err)
	}

	cfg, err := config.Load(configDir)
	if err != nil {
		log.Println(err)
		errList = append(errList, err)
	}

	if cfg != nil {
		defer cfg.Save()
	}

	db, err := openDatabase(cfg)

	if cfg.DBFile == "" {
		db, err = database.OpenMem()
		if err != nil {
			log.Println(err)
			errList = append(errList, err)
		}
	} else {
		db, err = database.Open(cfg.DBFile)
		if err != nil {
			errList = append(errList, err)
		}
	}
	if db != nil {
		defer db.Conn.Close()
	}

	uiVM := viewmodel.NewMainUI(cfg, db, errList)
	content := view.NewMainUI(window, uiVM)
	window.SetContent(content)
	window.ShowAndRun()
}



func openDatabase(cfg *config.Config ) (*database.Database, error) {
	if cfg.DBFile == "" {
		cfg.DBFile = ":memory:"
	}
	return database.Open(cfg.DBFile)
}
