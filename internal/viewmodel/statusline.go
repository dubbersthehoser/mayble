package viewmodel

import (
	"log"
	"time"

	"fyne.io/fyne/v2/data/binding"
)

const (
	StatusInfo int = iota
	StatusError
	StatusSuccess
)

type StatusLine struct {
	Text  binding.String
	Type    int
	clear   func() 
	clrTimer *time.Timer
	start chan struct{}
}

func newStatusLine() *StatusLine {
	sl := &StatusLine{
		clear: func() {
			log.Println("Warning: on clear not set")
		},
		clrTimer: time.NewTimer(0),
		start: make(chan struct{}),

		Text: binding.NewString(),
	}

	sl.clear = func() {
		sl.Text.Set("")
	}

	countDown := time.Duration(time.Minute / 10)

	go func() {
		for {
			select {
			case _ = <- sl.start:
				_ = sl.clrTimer.Reset(countDown)
			case _ = <- sl.clrTimer.C:
				sl.clear()
			}
		}
	}()

	return sl
}

func (sl *StatusLine) SetOnClear(fn func()) {
	sl.clear = fn
}

func (sl *StatusLine) startClear() {
	sl.start <- struct{}{}
}

func (sl *StatusLine) sendError(msg string) {
	sl.Type = StatusError
	_ = sl.Text.Set(msg)
	sl.startClear()
}

func (sl *StatusLine) sendInfo(msg string) {
	sl.Type = StatusInfo
	_ = sl.Text.Set(msg)
	sl.startClear()
}

func (sl *StatusLine) sendSuccess(msg string) {
	sl.Type = StatusSuccess
	_ = sl.Text.Set(msg)
	sl.startClear()
}
