package viewmodel

import (
	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/models"
	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

type storeUserMessaging struct {
	store repo.BookStore
	bus   *bus.Bus
}

func newStoreUserMessaging(s repo.BookStore, b *bus.Bus) *storeUserMessaging {
	return &storeUserMessaging{
		store: s,
		bus: b,
	}
}

func (u *storeUserMessaging) CreateBook(b *models.BookEntry) (int64, error) {
	var (
		eventName string = msgUserSuccess
		eventData string = "Book Added!"
	)
	id, err := u.store.CreateBook(b)
	if err != nil {
		eventName = msgUserError
		eventData = err.Error()
	}
	u.bus.Notify(bus.Event{
		Name: eventName,
		Data: eventData,
	})
	return id, err
}

func (u *storeUserMessaging) UpdateBook(b *models.BookEntry) (error) {
	var (
		eventName string = msgUserSuccess
		eventData string = "Book Updated!"
	)
	err := u.store.UpdateBook(b)
	if err != nil {
		eventName = msgUserError
		eventData = err.Error()
	}
	u.bus.Notify(bus.Event{
		Name: eventName,
		Data: eventData,
	})
	return err
}

func (u *storeUserMessaging) DeleteBook(id int64) (error) {
	var (
		eventName string = msgUserSuccess
		eventData string = "Book Deleted!"
	)
	err := u.store.DeleteBook(id)
	if err != nil {
		eventName = msgUserError
		eventData = err.Error()
	}
	u.bus.Notify(bus.Event{
		Name: eventName,
		Data: eventData,
	})
	return err
}
