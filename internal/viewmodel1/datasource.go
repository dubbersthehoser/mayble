package viewmodel 


import (
	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/app"
)


type dataSource struct {
	GetAllBooks func() ([]models.BookEntry, error)
	DeleteBook  func(int64) error
	UpdateBook  func(*models.BookEntry) error
	CreateBook  func(*models.BookEntry) (int64, error)
	AddListener func(fn func()) 
}

func newDataSourceFromService(a *app.Service) *dataSource {
	ds := &dataSource{
		GetAllBooks: a.GetAllBooks,
		DeleteBook: a.DeleteBook,
		UpdateBook: a.UpdateBook,
		CreateBook: a.CreateBook,
	}
}
