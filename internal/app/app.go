package app

import (
	"log"
	"os"
)

type State struct {
	DB     *Database
	Logs   *Logs
	Paths
}

func Init() (*State, error) {
	
	state := &State{}

	state.Logs = NewLogs()
	log.SetOutput(state.Logs)
	log.Printf("app: Hello, logs")

	var err error
	paths, err := NewPaths()
	if err != nil {
		return nil, err
	}
	state.Paths = *paths

	err = state.InitStoreage()
	if err != nil {
		return nil, err
	}
	log.Printf("app: Hello, storage: '%s'", state.StoragePath)

	file, err := os.OpenFile(state.LogsPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	state.Logs.Register(file)

	log.Printf("app: Hello, '%s'", state.LogsPath)

	state.DB, err = OpenDatabase(state.DBPath)
	if err != nil {
		return nil, err
	}
	log.Println("app: Hello, database!")
	
	err = DatabaseMigrateUp(state.DB)
	if err != nil {
		return nil, err
	}
	log.Println("app: Hello, state!")
	return state, nil
}
