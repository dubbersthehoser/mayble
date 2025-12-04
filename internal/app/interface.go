package app

import (
	"time"
)

type BookLoan struct {
	ID       int64
	Title    string
	Author   string
	Genre    string
	Ratting  int
	IsOnLoan   bool
	Borrower string
	Date     time.Time
}

type API interface {
	BookLoaning
	Importable
	Redoable
	Undoable
	Savable
}

type BookLoaning interface {
	CreateBookLoan(*BookLoan) error
	UpdateBookLoan(*BookLoan) error
	DeleteBookLoan(*BookLoan) error
	GetBookLoans() ([]BookLoan, error)
}

type Importable interface {
	ImportBookLoans([]BookLoan) error
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

// AppListerner defines avaliable event to listen to from application.
type AppListerner interface {
	SubscribeRedoUndoables
}

type SubscribeRedoUndoables {
	// SubscribeToRedos listen to DocumentRedo, and DocumentRedoEmpty from application document.
	SubscribeToUndos(func(*emiter.Event)

	// SubscribeToUndos listen to DocumentUndo, and DocumentUndoEmpty from application document.
	SubscribeToRedos(func(*emiter.Event))
}


