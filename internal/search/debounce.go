package search

import (
	"time"
	"context"
)

type Debouncer struct {
	delay  time.Duration
	timer  *time.Timer
	cancel context.CancelFunc
}

func (db *Debouncer) Search(pattern string, fn func(context.Context, string)) {
	if db.cancel != nil {
		db.cancel()
	}

	if db.timer != nil {
		db.timer.Stop()
	}

	db.timer = time.AfterFunc(db.delay, func(){
		ctx, cancel := context.WithCancel(context.Background())
		db.cancel = cancel

		fn(ctx, pattern)
	})

}

