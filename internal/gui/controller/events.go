package controller

type EventType string

const (
	openBookFormCreate EventType = "openBookFormCreate"
	openBookFormUpdate           = "openBookFormUpdate"
	closedBookForm               = "closedBookForm"
	bookListAdd                  = "bookListAdd"
	bookListRemove               = "bookListRemove"
)
