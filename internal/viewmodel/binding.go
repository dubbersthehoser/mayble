package viewmodel

import (
	"log"

	"github.com/dubbersthehoser/mayble/internal/bus"
	repo "github.com/dubbersthehoser/mayble/internal/repository"

	"fyne.io/fyne/v2/data/binding"
)

type UniqueGenres struct {
	genres repo.GenreRetriever
	l      *listener
}

func NewUniqueGenres(b *bus.Bus, subject sourceSubject, r repo.GenreRetriever) *UniqueGenres {
	ug := &UniqueGenres{
		genres: r,
		l:      &listener{},
	}
	subject.AddListener(ug.Update)
	ug.Update()
	return ug
}

func (u *UniqueGenres) Get() []string {
	s, err := u.genres.GetUniqueGenres()
	if err != nil {
		log.Println("unique.genres.get: ", err.Error())
		return []string{}
	}
	return s
}

func (u *UniqueGenres) Update() {
	u.l.notify()
}

func (u *UniqueGenres) AddListener(l binding.DataListener) {
	u.l.AddListener(l)
}
