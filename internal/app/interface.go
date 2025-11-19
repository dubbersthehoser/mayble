package app

import (
	"time"
)

type BookLoan struct {
	ID       int64
	Title    string
	Author   string
	Genre    string
	Ratting  string
	OnLoan   bool
	Borrower string
	Date     time.Time
}

type BookLoaning interface {
	CreateBookLoan(*BookLoan) error
	UpdateBookLoan(*BookLoan) error
	DeleteBookLoan(*BookLoan) error
	GetBookLoans() ([]BookLoan, error)
}

type Importable interface {
	ImportBookLoans([]data.BookLoan) error
}

type Mayble interface {
	BookLoaning
	Importable
	Redoable
	Undoable
	Savable
}


type Redoable interface {
	Redo() error
	RedoIsEmpty()  bool
}
type Undoable interface {
	Undo() error
	UndoIsEmpty() bool
}

type Savable interface {
	Save() error
}

