package viewmodel

import (
	"fmt"
	"github.com/dubbersthehoser/mayble/internal/event"
)

type NoDataState int

const (
	DataNoDB NoDataState = iota
	DataErr
)

type NoDataBody struct {
	s         NoDataState
	m         string
	OnChanged func()
}

func newNoDataBody(eb *event.EventBus) *NoDataBody {
	nb := &NoDataBody{
		s:         DataNoDB,
		m:         "Error: message not set",
		OnChanged: func() {},
	}

	eb.Subscribe(event.OpenedDatabase{}, func(v event.Event) {
		e := v.(event.OpenedDatabase)
		if !e.Failed {
			return
		}
		if e.Path == "" {
			nb.SetNoDB()
		} else {
			nb.SetDataErr(e.Path, e.Err)
		}
	})

	return nb
}

func (nb *NoDataBody) State() NoDataState {
	return nb.s
}

func (nb *NoDataBody) SetDataErr(path string, err error) {
	nb.s = DataErr
	nb.m = fmt.Sprintf("Something when wrong when opening database: \"%s\"\nError: %s", path, err)
	nb.OnChanged()
}

func (nb *NoDataBody) SetNoDB() {
	nb.s = DataNoDB
	nb.m = "Create, or Open a new database from the File drop-down."
	nb.OnChanged()
}

func (nb *NoDataBody) Message() string {
	return nb.m
}
