package stub

import (
	"log"

	"github.com/dubbersthehoser/mayble/internal/data"
)

type App struct {}

func (a *App) CreateBookLoan(bl *data.BookLoan) error {
	log.Printf("appstub: create: %#v", bl)
	return nil
}

func (a *App) DeleteBookLoan(bl *data.BookLoan) error {
	log.Printf("appstub: delete: %#v", bl)
	return nil
}

func (a *App) UpdateBookLoan(bl *data.BookLoan) error {
	log.Printf("appstub: update: %#v", bl)
	return nil
}
func (a *App) GetBookLoans() ([]data.BookLoan, error) {
	stubBookLoans := make([]data.BookLoan, 0)
	log.Printf("appstub: get book loans")
	return stubBookLoans, nil
}

func (a *App) ImportBookLoans(bl []data.BookLoan) error {
	log.Printf("appstub: importing book loans: %#v", bl)
	return nil
} 

func (a *App) Save() error {
	log.Printf("appstub: saving book loans")
	return nil
}

func (a *App) Undo() error {
	log.Printf("appstub: undo action")
	return nil
}
func (a *App) UndoIsEmpty() bool {
	log.Printf("appstub: is undo stack empty?")
	return true
}

func (a *App) Redo() error {
	log.Printf("appstub: redo action")
	return nil
}
func (a *App) RedoIsEmpty() bool {
	log.Printf("appstub: is redo stack empty?")
	return true
}

