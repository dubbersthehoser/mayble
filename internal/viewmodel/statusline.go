package viewmodel

import (
	"time"

	"fyne.io/fyne/v2/data/binding"
)

const (
	StatusInfo int = iota
	StatusError
	StatusSuccess
)

type StatusLine struct {
	Text      binding.String
	DoOnClear func()
	Type      int
	clrTimer  *time.Timer
	start     chan struct{}
}

func newStatusLine() *StatusLine {
	sl := &StatusLine{
		clrTimer: time.NewTimer(0),
		start:    make(chan struct{}),

		Text: binding.NewString(),
	}

	countDown := time.Duration(time.Minute / 10)

	go func() {
		for {
			select {
			case <-sl.start:
				_ = sl.clrTimer.Reset(countDown)
			case <-sl.clrTimer.C:
				sl.Clear()
			}
		}
	}()

	return sl
}

func (sl *StatusLine) Clear() {
	sl.Text.Set("")
	if sl.DoOnClear != nil {
		sl.DoOnClear()
	}
}

func (sl *StatusLine) startClearTimer() {
	sl.start <- struct{}{}
}

func (sl *StatusLine) sendError(msg string) {
	sl.Type = StatusError
	_ = sl.Text.Set(msg)
	sl.startClearTimer()
}

func (sl *StatusLine) sendInfo(msg string) {
	sl.Type = StatusInfo
	_ = sl.Text.Set(msg)
	sl.startClearTimer()
}

func (sl *StatusLine) sendSuccess(msg string) {
	sl.Type = StatusSuccess
	_ = sl.Text.Set(msg)
	sl.startClearTimer()
}
