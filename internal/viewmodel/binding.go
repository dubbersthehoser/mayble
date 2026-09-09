package viewmodel

import (
	"github.com/dubbersthehoser/mayble/internal/snapshot"
	"github.com/dubbersthehoser/mayble/internal/event"
	"github.com/dubbersthehoser/mayble/internal/config"
)

type UniqueGenres struct {
	l  []func()
	ss *snapshot.Snapshot
}

func newUniqueGenres(eb *event.EventBus) *UniqueGenres {
	ug := &UniqueGenres{
		ss: snapshot.Current.Load(),
	}
	eb.Subscribe(event.StoredSnapshot{}, func(_ event.Event) {
		ug.ss = snapshot.Current.Load()
		ug.notify()
	})
	return ug
}

func (ug *UniqueGenres) Genres() []string {
	return ug.ss.UniqueGenres()
}

func (ug *UniqueGenres) AddListener(fn func()) {

	if ug.l == nil {
		ug.l = make([]func(), 0)
	}

	ug.l = append(ug.l, fn)
}

func (ug *UniqueGenres) notify() {
	for _, fn := range ug.l {
		fn()
	}
}

type DBPath struct {
	cfg *config.Config
	l   []func()
}

func newDBPath(cfg *config.Config) *DBPath {
	dbp := &DBPath{
		cfg: cfg,
	}
	return dbp
}
func (p *DBPath) Get() string {
	return p.cfg.DBFile
}
func (p *DBPath) Set(s string) {
	p.cfg.DBFile = s
	p.notify()
}
func (p *DBPath) AddListener(fn func()) {
	if p.l == nil {
		p.l = make([]func(), 0)
	}
	p.l = append(p.l, fn)
}
func (p *DBPath) notify() {
	for _, fn := range p.l {
		fn()
	}
}
