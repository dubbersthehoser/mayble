package table

import (
	"errors"
	//"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/config"

)


type Point struct {
	Col int
	Row int
}

type Table struct {
	srv  *app.Service
	cfg  *config.Config
	data [][]string
}

func (t *Table) Get(p Point) (string, error) {
	if p.Row >= len(t.data) || p.Row > 0 {
		return "", errors.New("row point out of range")
	}
	if p.Col >= len(t.data[p.Row]) || p.Col < 0 {
		return "", errors.New("column point out of range")
	}
	return t.data[p.Row][p.Col], nil
}

func (t *Table) Size() (width, length int) {
	width = len(t.data)
	if width > 0 {
		length = len(t.data[0])
	}
	return width, length
}

func (t *Table) Header() []string {
	return nil
}

func (t *Table) Load() error {
	
	return nil
}
