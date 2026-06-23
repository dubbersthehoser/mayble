package viewmodel

import (
	"errors"
	"slices"
	"time"

	"fyne.io/fyne/v2/data/binding"
	
	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/app"
)

type BookForm struct{

	s *app.Service

	Fyne struct{
		Title binding.String
		Author binding.String
		Genre binding.String

		IsLoaned binding.Bool
		Borrower binding.String
		LoanedAt binding.String

		IsCompleted binding.Bool
		CompletedAt   binding.String
		Rating      binding.String

	}

	SubmitLabel string
	OnUpdate func()
	OnCreate func()
}

func newBookForm(onUpdate, onCreate func()) *BookForm {
	bf := &BookForm{
		OnUpdate: onUpdate,
		OnCreate: onCreate,
		Fyne: struct{
			Title binding.String
			Author binding.String
			Genre binding.String

			IsLoaned binding.Bool
			Borrower binding.String
			LoanedAt binding.String

			IsCompleted binding.Bool
			CompletedAt   binding.String
			Rating      binding.String
		}{
			Title: binding.NewString(),
			Author: binding.NewString(),
			Genre: binding.NewString(),

			IsLoaned: binding.NewBool(),
			Borrower: binding.NewString(),
			LoanedAt: binding.NewString(),

			IsCompleted: binding.NewBool(),
			Rating:      binding.NewString(),
			CompletedAt: binding.NewString(),
		},
	}
	return bf
}

func (bf *BookForm) Reset() {
	_ = bf.Fyne.Title.Set("")
	_ = bf.Fyne.Author.Set("")
	_ = bf.Fyne.Genre.Set("")
	_ = bf.Fyne.CompletedAt.Set("")
	_ = bf.Fyne.Rating.Set("")
	_ = bf.Fyne.LoanedAt.Set("")
	_ = bf.Fyne.Borrower.Set("")
	_ = bf.Fyne.IsCompleted.Set(false)
	_ = bf.Fyne.IsLoaned.Set(false)
}

func (bf *BookForm) Set(book *models.BookEntry) {
	_ = bf.Fyne.Title.Set(book.Title)
	_ = bf.Fyne.Author.Set(book.Author)
	_ = bf.Fyne.Genre.Set(book.Genre)

	_ = bf.Fyne.IsCompleted.Set(book.IsCompleted)
	if book.IsCompleted {
		_ = bf.Fyne.CompletedAt.Set(book.CompletedAt.Format(dateFormat))
		_ = bf.Fyne.Rating.Set(Ratings()[book.Rating])
	}

	bf.Fyne.IsLoaned.Set(book.IsLoaned)
	if book.IsLoaned {
		_ = bf.Fyne.LoanedAt.Set(book.LoanedAt.Format(dateFormat))
		_ = bf.Fyne.Borrower.Set(book.Borrower)
	}
}

func (bf *BookForm) GetBookEntry() (*models.BookEntry, error) {
	
	if err := bf.validate(); err != nil {
		return nil, err
	}
	
	book := models.BookEntry{}

	book.Title, _ = bf.Fyne.Title.Get()
	book.Author, _ = bf.Fyne.Author.Get()
	book.Genre, _ = bf.Fyne.Genre.Get()

	book.IsLoaned, _ = bf.Fyne.IsLoaned.Get()
	if book.IsLoaned {
		var err error
		date, _ := bf.Fyne.LoanedAt.Get()
		book.LoanedAt, err = time.Parse(dateFormat, date)
		if err != nil {
			return nil, err
		}
		book.Borrower, _ = bf.Fyne.Borrower.Get()
	}

	book.IsCompleted, _ = bf.Fyne.IsCompleted.Get()
	if book.IsCompleted {
		var err error
		date, _ := bf.Fyne.CompletedAt.Get()
		book.CompletedAt, err = time.Parse(dateFormat, date)
		if err != nil {
			return nil, err
		}
		rs, _ := bf.Fyne.Rating.Get()
		rating, err := parseRating(rs)
		if err != nil {
			return nil, err
		}
		book.Rating = rating
	}
	return &book, nil
}

func (bf *BookForm) IsLoaned() bool {
	ok, _ := bf.Fyne.IsLoaned.Get()
	return ok
}
func (bf *BookForm) IsCompleted() bool {
	ok, _ := bf.Fyne.IsCompleted.Get()
	return ok
}

func (bf *BookForm) validate() error {
	title, _ := bf.Fyne.Title.Get()
	author, _ := bf.Fyne.Author.Get()
	genre, _ := bf.Fyne.Genre.Get()

	if title == "" {
		return errors.New("missing title")
	}
	if author == "" {
		return errors.New("missing auther")
	}
	if genre == "" {
		return errors.New("missing genre")
	}

	isLoaned, _ := bf.Fyne.IsLoaned.Get()
	isRead, _ := bf.Fyne.IsCompleted.Get()

	if isLoaned {
		date, _ := bf.Fyne.LoanedAt.Get()
		borrower, _ := bf.Fyne.Borrower.Get()

		if borrower == "" {
			return errors.New("missing borrower")
		}

		if date == "" {
			return errors.New("missing borrower date")
		}

		_, err := time.Parse(dateFormat, date)
		if err != nil {
			return errors.New("invalid date for loaned")
		}
	}

	if isRead {
		date, _ := bf.Fyne.CompletedAt.Get()
		rating, _ := bf.Fyne.Rating.Get()

		if date == "" {
			return errors.New("missing completion date")
		}

		_, err := time.Parse(dateFormat, date)
		if err != nil {
			return errors.New("invalid date for completion")
		}

		ratings := Ratings()
		rank := slices.Index(ratings, rating)
		if rank == 0 {
			return errors.New("rating not selected")
		}
		if rank == -1 {
			return errors.New("invalid rating")

		}
	}
	return nil
}

