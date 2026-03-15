package viewmodel

import (
	"log"

	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/database"
	repo "github.com/dubbersthehoser/mayble/internal/repository"

	"fyne.io/fyne/v2/data/binding"
)

type appService struct {
	cfg *config.Config
	dbs *database.Service

	bookRetriever  repo.BookRetriever
	genreRetriever repo.GenreRetriever
	bookCreator    repo.BookCreator
	bookUpdator    repo.BookUpdator
	bookDeletor    repo.BookDeletor

	uniqueGenres *UniqueGenres
}

func (as *appService) changeDB(db *database.Database) error {
	err := as.dbs.SetDB(db)
	if err != nil {
		return err
	}
	as.setRepos(db)
	return nil
}

func (as *appService) setRepos(db *database.Database) {
	as.bookRetriever = db
	as.genreRetriever = db
	as.bookCreator = db
	as.bookUpdator = db
	as.bookDeletor = db
}

func newAppService(bus *bus.Bus, cfg *config.Config, db *database.Database) *appService {
	as := &appService{
		cfg: cfg,
		dbs: database.NewService(db),
	}
	as.setRepos(db)
	as.uniqueGenres = NewUniqueGenres(bus, as.genreRetriever)
	return as
}

type UniqueGenres struct {
	list   binding.StringList
	genres repo.GenreRetriever
	l      *listener
}

func NewUniqueGenres(b *bus.Bus, g repo.GenreRetriever) *UniqueGenres {
	ug := &UniqueGenres{
		list:   binding.NewStringList(),
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
		return []string{}
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
