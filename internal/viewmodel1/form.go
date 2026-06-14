package viewmodel

import (
	"errors"
	"slices"
	"time"
	
	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/app"
)

type BookForm struct{

	s *app.Service

	entry struct{
		title, author, genre string
		isCompleted, isLoaned bool
		completedAt *time.Time
		rating string
		loanedAt  *time.Time
		borrower string
	}

	SubmitLabel string
	OnUpdate func()
	OnCreate func()

	l []func()
}

func (bf *BookForm) Reset() {
	bf.entry.title = ""
	bf.entry.author = ""
	bf.entry.genre = ""
	bf.entry.completedAt = nil
	bf.entry.rating = ""
	bf.entry.loanedAt = nil
	bf.entry.borrower = ""
	bf.entry.isCompleted = false
	bf.entry.isLoaned = false
}

func (bf *BookForm) SetTitle(s string) {
	bf.entry.title = s
	bf.notify()
}
func (bf *BookForm) GetTitle() string {
	return bf.entry.title
}

func (bf *BookForm) SetAuthor(s string) {
	bf.entry.author = s
	bf.notify()
}
func (bf *BookForm) GetAuthor() string {
	return bf.entry.author
}

func (bf *BookForm) SetGenre(s string) {
	bf.entry.genre = s
	bf.notify()
}
func (bf *BookForm) GetGenre() string {
	return bf.entry.genre
}

func (bf *BookForm) SetRating(s string) {
	bf.entry.rating = s
	bf.notify()
}

func (bf *BookForm) GetRating() string {
	return bf.entry.rating 
}

func (bf *BookForm) SetCompletedAt(t *time.Time) {
	bf.entry.completedAt = t
	bf.notify()
}
func (bf *BookForm) GetCompletedAt() *time.Time {
	return bf.entry.completedAt
}

func (bf *BookForm) SetBorrower(s string) {
	bf.entry.borrower = s
	bf.notify()
}
func (bf *BookForm) GetBorrower() string{
	return bf.entry.borrower
}

func (bf *BookForm) SetLoanedAt(t *time.Time) {
	bf.entry.loanedAt = t
	bf.notify()
}
func (bf *BookForm) GetLoanedAt() *time.Time {
	return bf.entry.loanedAt
}

func (bf *BookForm) SetLoaned(t bool) {
	bf.entry.isLoaned = t
	bf.notify()
}

func (bf *BookForm) IsLoaned() bool {
	return bf.entry.isLoaned
}

func (bf *BookForm) SetCompleted(t bool) {
	bf.entry.isCompleted = t
	bf.notify()
}

func (bf *BookForm) IsCompleted() bool {
	return bf.entry.isCompleted
}

func (bf *BookForm) GetBookEntry() (*models.BookEntry, error) {
	
	if err := bf.validate(); err != nil {
		return nil, err
	}
	
	book := models.BookEntry{}

	book.Title = bf.entry.title
	book.Author = bf.entry.author
	book.Genre = bf.entry.genre

	book.IsLoaned = bf.entry.isLoaned
	if bf.entry.isLoaned {
		date := bf.entry.loanedAt
		book.LoanedAt = *date
		book.Borrower = bf.entry.borrower
	}

	book.IsCompleted = bf.entry.isCompleted
	if bf.entry.isCompleted {
		date := bf.entry.completedAt
		book.CompletedAt = *date

		rating, err := parseRating(bf.entry.rating)
		if err != nil {
			return nil, err
		}
		book.Rating = rating
	}
	return &book, nil
}

func (bf *BookForm) validate() error {
	title := bf.entry.title
	author := bf.entry.author
	genre := bf.entry.genre

	if title == "" {
		return errors.New("missing title")
	}
	if author == "" {
		return errors.New("missing auther")
	}
	if genre == "" {
		return errors.New("missing genre")
	}

	isLoaned := bf.entry.isLoaned
	isRead := bf.entry.isCompleted

	if isLoaned {
		date := bf.entry.loanedAt
		borrower := bf.entry.borrower

		if borrower == "" {
			return errors.New("missing borrower")
		}
		if date == nil {
			return errors.New("missing borrower date")
		}
	}

	if isRead {
		completed := bf.entry.completedAt
		rating := bf.entry.rating

		if completed == nil {
			return errors.New("missing completion date")
		}

		ratings := Ratings()
		rank := slices.Index(ratings, rating)
		if rank == 0 {
			return errors.New("ratting not selected")
		}
		if rank == -1 {
			return errors.New("invalid rating")

		}
	}
	return nil
}

func (bf *BookForm) AddListener(fn func()) {
	if bf.l == nil {
		bf.l = make([]func(), 0)
	}
	bf.l = append(bf.l, fn)
}

func (bf *BookForm) notify() {
	for _, fn := range bf.l {
		fn()
	}
}

