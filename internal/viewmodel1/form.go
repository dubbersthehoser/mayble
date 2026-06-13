package viewmodel

import (
	"errors"
	"slices"
	
	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/app"
)

type BookForm struct{

	s *app.Service

	entry struct{
		title, author, genre string
		isCompleted, isLoaned bool
		rating, completed string
		borrower, loaned  string
	}

	SubmitLabel string
	OnUpdate func()
	OnCreate func()
}

func (bf *BookForm) Reset() {
	bf.entry.title = ""
	bf.entry.author = ""
	bf.entry.genre = ""
	bf.entry.rating = ""
	bf.entry.completed = ""
	bf.entry.borrower = ""
	bf.entry.isCompleted = false
	bf.entry.isLoaned = false
}

func (bf *BookForm) SetTitle(s string) {
	bf.entry.title = s
}

func (bf *BookForm) SetAuthor(s string) {
	bf.entry.author = s
}

func (bf *BookForm) SetGenre(s string) {
	bf.entry.genre = s
}

func (bf *BookForm) SetRating(s string) {
	bf.entry.rating = s
}

func (bf *BookForm) SetCompletedAt(s string) {
	bf.entry.completed = s
}

func (bf *BookForm) SetBorrower(s string) {
	bf.entry.borrower = s
}

func (bf *BookForm) SetLoanedAt(s string) {
	bf.entry.loaned = s
}

func (bf *BookForm) SetLoaned(t bool) {
	bf.entry.isLoaned = t
}

func (bf *BookForm) SetCompleted(t bool) {
	bf.entry.isCompleted = t
}

func (bf *BookForm)GetBookEntry() (*BookForm, error) {
	
	if err := bf.validate(); != nil {
		return nil, err
	}
	
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
		date := bf.entry.loaned
		borrower := bf.entry.borrower

		if borrower == "" {
			return errors.New("missing borrower")
		}
		if date == "" {
			return errors.New("missing borrower date")
		}
		_, err := parseDate(date)
		if err != nil {
			return errors.New("invalid borrower date")
		}
	}

	if isRead {
		completed := bf.entry.completed
		rating := bf.entry.rating

		if completed == "" {
			return errors.New("missing completion date")
		}
		_, err := parseDate(completed)
		if err != nil {
			return errors.New("invalid completion date")
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

