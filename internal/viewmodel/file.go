package viewmodel

import (
	"fyne.io/fyne/v2"

	//"github.com/dubbersthehoser/mayble/internal/csv"
)

type FileVM struct {
	
}

func (c *FileVM) ImportCSV(r fyne.URIReadCloser, err error) {
	println(err)
}

func (c *FileVM) ExportCSV(w fyne.URIWriteCloser, err error) {
	println(err)
}

func (c *FileVM) OpenDatabase(r fyne.URIReadCloser, err error) {
	println(err)
}

func (c *FileVM) CreateDatabase(r fyne.URIWriteCloser, err error) {
	println(err)
}


