package viewmodel

import (
	"fyne.io/fyne/v2/data/binding"
)

type MenuVM struct {
	DBFile binding.String
}

func NewMenuVM(dbFile binding.String) *MenuVM {
	m := &MenuVM{
		DBFile: dbFile,
	}
	return m
} 

func (m *MenuVM) NewFileVM() *FileVM {
	c := &FileVM{}
	return c
}

