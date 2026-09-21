package viewmodel

import (
	"fmt"
	"time"

	"github.com/dubbersthehoser/mayble/internal/event"
)

const (
	StatusInfo int = iota
	StatusError
	StatusSuccess
)

type StatusLine struct {
	DoOnClear func()
	OnChanged func(string, int)
	Type      int
	clrTimer  *time.Timer
}

func newStatusLine(eb *event.EventBus) *StatusLine {
	sl := &StatusLine{
		clrTimer:  time.NewTimer(0),
		DoOnClear: func() {},
		OnChanged: func(_ string, _ int) {},
	}

	setupStatusLineToEvents(sl, eb)
	return sl
}

func (sl *StatusLine) startClearTimer() {
	countDown := time.Duration(time.Minute / 10)
	if sl.clrTimer != nil {
		sl.clrTimer.Stop()
	}
	sl.clrTimer = time.AfterFunc(countDown, sl.DoOnClear)
}

func (sl *StatusLine) sendError(msg string) {
	sl.OnChanged(msg, StatusError)
	sl.startClearTimer()
}

func (sl *StatusLine) sendInfo(msg string) {
	sl.OnChanged(msg, StatusInfo)
	sl.startClearTimer()
}

func (sl *StatusLine) sendSuccess(msg string) {
	sl.OnChanged(msg, StatusSuccess)
	sl.startClearTimer()
}

func setupStatusLineToEvents(sl *StatusLine, eb *event.EventBus) {
	eb.Subscribe(event.CreatedBookEntry{}, func(v event.Event) {
		e := v.(event.CreatedBookEntry)
		if e.Failed {
			sl.sendError(e.Message)
		} else {
			sl.sendSuccess("Entry Added!")
		}
	})

	eb.Subscribe(event.UpdatedBookEntry{}, func(v event.Event) {
		e := v.(event.UpdatedBookEntry)
		if e.Failed {
			sl.sendError(e.Message)
		} else {
			sl.sendSuccess("Entry Updated!")
		}
	})
	eb.Subscribe(event.DeletedBookEntry{}, func(v event.Event) {
		e := v.(event.DeletedBookEntry)
		if e.Failed {
			sl.sendError(e.Message)
		} else {
			sl.sendSuccess("Entry Removed!")
		}
	})
	eb.Subscribe(event.ImportedFile{}, func(v event.Event) {
		e := v.(event.ImportedFile)
		if e.Failed {
			sl.sendError(e.Message)
		} else {
			sl.sendInfo(fmt.Sprintf("Imported %s", e.Path))
		}
	})
	eb.Subscribe(event.ExportedFile{}, func(v event.Event) {
		e := v.(event.ExportedFile)
		if e.Failed {
			sl.sendError(e.Message)
		} else {
			sl.sendInfo(fmt.Sprintf("Exported to %s", e.Path))
		}
	})
	eb.Subscribe(event.OpenedDatabase{}, func(v event.Event) {
		e := v.(event.OpenedDatabase)
		if e.Failed {
			sl.sendError(e.Message)
		} else {
			sl.sendInfo(fmt.Sprintf("Opened %s", e.Path))
		}
	})
	eb.Subscribe(event.CreatedDatabase{}, func(v event.Event) {
		e := v.(event.CreatedDatabase)
		if e.Failed {
			sl.sendError(e.Message)
		} else {
			sl.sendInfo(fmt.Sprintf("Created %s", e.Path))
		}
	})
}
