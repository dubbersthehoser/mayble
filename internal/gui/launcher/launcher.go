package launcher

import (
	
	"fmt"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	
	"github.com/dubbersthehoser/mayble/internal/core"
	"github.com/dubbersthehoser/mayble/internal/storage"
	//_"github.com/dubbersthehoser/mayble/internal/sqlitedb"
	"github.com/dubbersthehoser/mayble/internal/memdb"
	"github.com/dubbersthehoser/mayble/pkg/config"
	"github.com/dubbersthehoser/mayble/internal/gui/controller"
	"github.com/dubbersthehoser/mayble/internal/gui/view"
)



//func loadSqliteDatabase(s *Settings) (*sqlitedb.Database, error) {
//	database := sqlitedb.NewDatabase()
//	if err := database.Open(s.dbPath); err != nil {
//		return nil, err
//	}
//	if err := database.MigrateUp(); err != nil {
//		return nil, err
//	}
//	return database, nil
//}

func loadMemDatabase(s *Settings) (*memdb.MemStorage, error) {
	store := memdb.NewMemStorage()
	return store, nil
} 

func loadDatabase(s *Settings, driver string) (storage.Storage, error) {
	switch driver {
	case "memory":
		return loadMemDatabase(s)
	//case "sqlite":
	//	return loadSqliteDatabase(s)
	default:
		return nil, fmt.Errorf("driver '%s', not found", driver)
	}
}

func loadConfig(s *Settings) (*config.Config, error) {
	return config.Load(s.configPath)
}


func SetupTextGrid() *widget.TextGrid {
	textgrid := widget.NewTextGrid()
	textgrid.Scroll = fyne.ScrollBoth
	return textgrid
}


func Run(options ...Option) {

	s := Settings{}
	
	for _, option := range options {
		option(&s)
	}

	defaultConfigDir(&s)
	defaultConfigPath(&s)
	defaultDBPath(&s)
	
	App := app.NewWithID("app.dubbersthehoser.mayble")
	window := App.NewWindow("Mayble Launcher")
	window.Resize(fyne.NewSize(800, 500))
	window.CenterOnScreen()

	logGrid := SetupTextGrid()

	logGrid.Append("Hello, World!")
	logGrid.Append("--- PATHS ---")
	logGrid.Append(fmt.Sprintf("config: '%s'", s.configPath))
	logGrid.Append(fmt.Sprintf("storage: '%s'", s.dbPath))
	logGrid.Append("\n--- LOADING ---")

	var Errored bool = false

	_, err := loadConfig(&s)
	if err != nil {
		logGrid.Append(fmt.Sprintf("- config: failed: %s", err.Error()))
		Errored = true
	} else {
		logGrid.Append("- config: success")
	}

	storage, err := loadDatabase(&s, "memory")
	if err != nil && !Errored{
		logGrid.Append(fmt.Sprintf("- storage: failed: %s", err.Error()))
		Errored = true
	} else {
		logGrid.Append("- storage: success")
	}

	core, err := core.New(storage)
	if err != nil && !Errored {
		logGrid.Append(fmt.Sprintf("- core: failed: %s", err.Error()))
		Errored = true
	} else {
		logGrid.Append("- core: success")
	}

	master := controller.New(core)
	funkView, err := view.NewFunkView(master, window)
	if err != nil {
		logGrid.Append(fmt.Sprintf("- view: failed: %s", err.Error()))
		Errored = true
	}

	mainView := funkView.View

	if Errored  {
		mainView = logGrid
	}

	window.SetContent(mainView)
	window.ShowAndRun()
}




