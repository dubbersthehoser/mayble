package table

import (
	"errors"
	//"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/app"

)

type Point struct {
	Col int
	Row int
}

type Table struct {
	service *app.Service
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

