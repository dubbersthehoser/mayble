package viewmodel

import (
	repo "github.com/dubbersthehoser/mayble/internal/repository"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/bus"

	"fyne.io/fyne/v2/data/binding"
)

type appService struct {
	cfg  *config.Config
	bookRetriever  repo.BookRetriever
	genreRetriever repo.GenreRetriever
	bookCreator    repo.BookCreator
	bookUpdator    repo.BookUpdator
	bookDeletor    repo.BookDeletor
}

func newAppService(cfg *config.Config, a *app.Application) *appService {
	return &appService{
		cfg: cfg,
		bookRetriever:  a,
		genreRetriever: a,
		bookCreator:    a,
		bookUpdator:    a,
		bookDeletor:    a,
	}
}


type vmService struct {
	bus    *bus.Bus
	genres *UniqueGenres
	app    *appService
}

func newVMService(as *appService) *vmService {
	bus := &bus.Bus{}
	genres := NewUniqueGenres(bus, binding.NewStringList(), as.genreRetriever)
	vs := &vmService{
		bus: bus,
		genres: genres,
		app: as,
	}
	return vs
}



type UniqueGenres struct {
	list   binding.StringList
	genres repo.GenreRetriever
	l      *listener
}

func NewUniqueGenres(b *bus.Bus, l binding.StringList, g repo.GenreRetriever) *UniqueGenres {
	ug := &UniqueGenres{
		list:   l,
		genres: g,
		l:      &listener{},
	}
	b.Subscribe(bus.Handler{
		Name: msgDataChanged,
		Handler: func(e *bus.Event) {
			ug.Update()
		},
	})
	ug.Update()
	return ug
}

func (u *UniqueGenres) Get() []string {
	s, err := u.list.Get()
	if err != nil {
		return []string{
			"__STUB__",
			"__STUB__",
			"__STUB__",
			"__STUB__",
		}
	}
	return s
}

func (u *UniqueGenres) Update() {
	genres, err := u.genres.GetUniqueGenres()
	if err != nil {
		return 
	}
	for i := range u.list.Length() {
		v, _ := u.list.GetValue(i)
		_ = u.list.Remove(v)
	}
	for i := range genres {
		_ = u.list.Append(genres[i])
	}
	u.l.notify()
}

func (u *UniqueGenres) AddListener(l binding.DataListener) {
	u.l.AddListener(l)
}
