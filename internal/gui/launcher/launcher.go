package launcher

import (
	
	"fmt"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	//"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	
	//"github.com/dubbersthehoser/mayble/internal/core"
	"github.com/dubbersthehoser/mayble/internal/sqlitedb"
	"github.com/dubbersthehoser/mayble/pkg/config"
)

func loadConfig(s *Settings) (*config.Config, error) {
	return config.Load(s.configPath)
}

func loadDatabase(s *Settings) (*sqlitedb.Database, error) {
	database := sqlitedb.NewDatabase()
	if err := database.Open(s.dbPath); err != nil {
		return nil, err
	}
	return database, nil
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
	
	App := app.New()
	window := App.NewWindow("Mayble Launcher")
	window.Resize(fyne.NewSize(800, 500))
	logGrid := SetupTextGrid()

	logGrid.Append("Hello, World!")
	logGrid.Append(fmt.Sprintf("config: %s\ndatabase: %s", s.configPath, s.dbPath))

	_, err := loadConfig(&s)
	if err != nil {
		logGrid.Append(err.Error())
	}
	_, err = loadDatabase(&s)
	if err != nil {
		logGrid.Append(err.Error())
	}

	window.SetContent(logGrid)
	window.ShowAndRun()


}
