package viewmodel

import (
	"log"

	repo "github.com/dubbersthehoser/mayble/internal/repository"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/database"

	"fyne.io/fyne/v2/data/binding"
)

type appService struct {
	cfg  *config.Config
	dbs  *database.Service
	
	bookRetriever  repo.BookRetriever
	genreRetriever repo.GenreRetriever
	bookCreator    repo.BookCreator
	bookUpdator    repo.BookUpdator
	bookDeletor    repo.BookDeletor
}

func (as *appService) setDB(db *database.Database) error {
	if as.dbs == nil {
		as.dbs = database.NewService(db)
	} else {
		err := as.dbs.SetDB(db)
		if err != nil {
			return err
		}
	}
	as.bookRetriever = db
	as.genreRetriever = db
	as.bookCreator = db
	as.bookUpdator = db
	as.bookDeletor = db
	return nil
}

func newAppService(cfg *config.Config, dbs *database.Service) *appService {
	as := &appService{
		cfg: cfg,
		dbs: dbs,
	}
	as.setDB(dbs.DB())
	return as
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
		log.Println("unique.genres.get: ", err.Error())
		return []string{
		}
	}
	return s
}

func (u *UniqueGenres) Update() {
	genres, err := u.genres.GetUniqueGenres()
	if err != nil {
		log.Println("unique.genres.update: ", err.Error())
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
