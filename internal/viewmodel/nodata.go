package viewmodel

import (
	"fmt"
)

type NoDataState int
const (
	DataNoDB NoDataState = iota
	DataErr
)

type NoDataBody struct {
	s NoDataState
	m string
	l []func()
}

func (nb *NoDataBody) State() NoDataState {
	return nb.s
}

func (nb *NoDataBody) SetDataErr(path string, err error) {
	nb.s = DataErr 
	nb.m = fmt.Sprintf("Something when wrong when opening database: \"%s\"\nError: %s", path, err)	
	nb.notify()
}

func (nb *NoDataBody) SetNoDB() {
	nb.s = DataNoDB
	nb.m = "Create or open a new database to work on."
	nb.notify()
}

func (nb *NoDataBody) AddListener(fn func()) {
	if nb.l == nil {
		nb.l = make([]func(), 0)
	}

	nb.l = append(nb.l, fn)
}

func (nb *NoDataBody) notify() {
	for _, fn := range nb.l {
		fn()
	}
}
