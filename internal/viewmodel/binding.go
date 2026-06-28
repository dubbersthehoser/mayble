package viewmodel

import (
	"log"

	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/config"
)

type UniqueGenres struct {
	s *app.Service

	l []func()
}

func newUniqueGenres(s *app.Service) *UniqueGenres {
	ug := &UniqueGenres{
		s: s,
	}
	ug.s.AddListener(func() {
		ug.notify()
	})
	return ug
}

func (ug *UniqueGenres) Genres() []string {
	g, err := ug.s.GetUniqueGenres()
	if err != nil {
		log.Println("Error:", err)
	}
	return g
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
	l []func()
}
func newDBPath(cfg *config.Config) *DBPath{
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

type Body struct {
	value int
	l []func()
}

func (b *Body) Value() int {
	return b.value
}

func (b *Body) Set(v int) {
	b.value = v
	b.notify()
}

func (b *Body) AddListener(fn func()) {
	if b.l == nil {
		b.l = make([]func(), 0)
	}
	b.l = append(b.l, fn)
}

func (b *Body) notify() {
	for _, fn := range b.l {
		fn()
	}
}
