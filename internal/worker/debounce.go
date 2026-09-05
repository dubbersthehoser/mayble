package worker

import (
	"sync"
	"time"
)

type Debouncer struct {
	mu    sync.Mutex
	timer *time.Timer
	delay time.Duration
	fn     func()
}

func NewDebouncer(delay time.Duration, fn func()) *Debouncer {
	d := &Debouncer{
		delay: delay,
	}
	return d
}

func (d *Debouncer) Debounce() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer == nil {
		d.timer = time.AfterFunc(d.delay, d.fn)
		return
	}
	d.timer.Reset(d.delay)
}

func Debounce(d time.Duration) func(func()) {
	var timer *time.Timer
	return func(fn func()) {
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(d, fn)
	}
}
