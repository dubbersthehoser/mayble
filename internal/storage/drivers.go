package storage

import (
	//_"github.com/dubbersthehoser/mayble/internal/sqlite"
	"github.com/dubbersthehoser/mayble/internal/memory"
)

//func loadSQLite(path string) (*sqlite.Database, error) {
//	database := sqlitedb.NewDatabase()
//	if err := database.Open(path); err != nil {
//		return nil, err
//	}
//	if err := database.MigrateUp(); err != nil {
//		return nil, err
//	}
//	return database, nil
//}

func loadMemory(path string) (*memstore.MemStore, error) {
	store := memory.NewMemory()
	return store, nil
} 

func Load(driver, path string) (Storage, error){
	switch dirver {
	case "memory":
		return loadMemory(path)
	//case "sqlite":
	//	return loadSQLite(path)
	default:
		return nil, fmt.Errorf("storage: driver '%s', not found", driver)
	}
}
