package viewmodel

import (
	"errors"
	"slices"
	"fmt"

	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/bus"
	repo "github.com/dubbersthehoser/mayble/internal/repository"

)

type BookForm struct {
	Title     binding.String
	Author    binding.String
	Genre     binding.String

	UniqGenre binding.StringList

	IsLoaned binding.Bool
	Borrower binding.String
	Date     binding.String

	IsRead    binding.Bool
	Rating    binding.String
	Completed binding.String
}


type SubmissionList struct {
	form       *CreateBookForm
	bus        *bus.Bus
	l          *listener
	submissions []repo.BookEntry
}
func NewSubmissionList(bus *bus.Bus) *SubmissionList{
	return &SubmissionList{
		l: &listener{},
		bus: bus,
		submissions: make([]repo.BookEntry, 0),
	}
}

func (s *SubmissionList) Clear() {
	s.submissions = s.submissions[:0]
}

func (s *SubmissionList) addSubmission(bf *BookForm) error {
	book := repo.BookEntry{}

	isOnLoan, _ := bf.IsLoaned.Get()
	isRead, _ := bf.IsRead.Get()

	if isOnLoan {
		book.Variant |= repo.Loaned
		date, _ := bf.Date.Get()
		timeDate, _ := parseDate(date)
		book.Loaned = *timeDate
		book.Borrower, _ = bf.Borrower.Get()
	}

	if isRead {
		book.Variant |= repo.Read
		sdate, _ := bf.Completed.Get()
		srating, _ := bf.Rating.Get()

		date, _ := parseDate(sdate)
		rating := slices.Index(Ratings(), srating)

		book.Read = *date
		book.Rating = rating
	}

	book.Variant |= repo.Book

	book.Title, _ = bf.Title.Get()
	book.Author, _ = bf.Author.Get()
	book.Genre, _ = bf.Genre.Get()

	s.submissions = append(s.submissions, book)
	s.l.notify()
	s.bus.Notify(bus.Event{
		Name: msgUserSuccess,
		Data: "Added to submission",
	})
	return nil
}

func (s *SubmissionList) Length() int {
	return len(s.submissions)
}

func (s *SubmissionList) Get(idx int) string {
	if idx >= s.Length() || idx < 0 {
		return "STUB: out of range"
	}
	book := s.submissions[idx]
	switch book.Variant {
	case repo.BookReadAndLoaned:
		return fmt.Sprintf("%s %s %s (loaned) (read)", book.Title, book.Author, book.Genre)
	case repo.BookLoaned:
		return fmt.Sprintf("%s %s %s (loaned)", book.Title, book.Author, book.Genre)
	case repo.BookRead:
		return fmt.Sprintf("%s %s %s (read)", book.Title, book.Author, book.Genre)
	case repo.Book:
		return fmt.Sprintf("%s %s %s", book.Title, book.Author, book.Genre)
	default:
		return "STUB: variant not found"
	}
}

func (s *SubmissionList) Remove(idx int) {
	if idx >= s.Length() || idx < 0 {
		return
	}

	s.submissions = append(s.submissions[:idx], s.submissions[idx+1:]...)

	s.bus.Notify(bus.Event{
		Name: msgUserInfo,
		Data: "Submission removed",
	})

	s.l.notify()
}

func (s *SubmissionList) Edit(idx int) {
	if s.form == nil {
		return
	}

	if idx >= s.Length() || idx < 0 {
		return
	}

	book := s.submissions[idx]

	_ = s.form.Title.Set(book.Title)
	_ = s.form.Author.Set(book.Author)
	_ = s.form.Genre.Set(book.Genre)

	if repo.Loaned & book.Variant != 0 {
		_ = s.form.IsLoaned.Set(true)
		_ = s.form.Date.Set(formatDate(&book.Loaned))
		_ = s.form.Borrower.Set(book.Borrower)
	}

	if repo.Read & book.Variant != 0 {
		_ = s.form.IsRead.Set(true)
		_ = s.form.Rating.Set(formatRating(book.Rating))
		_ = s.form.Completed.Set(formatDate(&book.Read))
	}
	s.Remove(idx)

}

func (s *SubmissionList) AddListener(l binding.DataListener) {
	s.l.AddListener(l)
}

type CreateBookForm struct {
	bus       *bus.Bus
	sl        *SubmissionList
	BookForm
}
func NewCreateBookForm(bus *bus.Bus) *CreateBookForm {
	bf := &CreateBookForm{
		sl: NewSubmissionList(bus),
		bus: bus,
		BookForm: BookForm{
			Title: binding.NewString(),
			Author: binding.NewString(),
			Genre: binding.NewString(),

			IsRead:    binding.NewBool(),
			Rating:    binding.NewString(),
			Completed: binding.NewString(),

			IsLoaned: binding.NewBool(),
			Borrower: binding.NewString(),
			Date: binding.NewString(),
		},
	}
	bf.sl.form = bf
	return bf
}

func validate(bf *BookForm) error {
	title, _ := bf.Title.Get()
	author, _ := bf.Author.Get()
	genre, _ := bf.Genre.Get()

	if title == "" {
		return errors.New("Missing Title")
	}
	if author == "" {	
		return errors.New("Missing Auther")
	} 
	if genre == "" {
		return errors.New("Missing Genre")
	}

	isLoaned, _ := bf.IsLoaned.Get()
	isRead, _ := bf.IsRead.Get()

	if isLoaned {
		date, _ := bf.Date.Get()
		borrower, _ := bf.Borrower.Get()

		if borrower == "" {
			return errors.New("Missing Borrower")
		}
		if date == "" {
			return errors.New("Missing Borrower Date")
		}
		_, err := parseDate(date)
		if err != nil {
			return errors.New("Invalid Borrower Date (DD/MM/YYYY)")
		}
	}

	if isRead {
		completed, _ := bf.Completed.Get()
		rating, _ := bf.Rating.Get()

		if completed == "" {
			return errors.New("Missing Completion Date")
		}
		_, err := parseDate(completed)
		if err != nil {
			return errors.New("Invalid Completion Date (DD/MM/YYYY)")
		}
		// convert rating to int
		ratings := Ratings()
		rank := slices.Index(ratings[1:], rating)
		if rank == -1 {
			return errors.New("Invalid Rating")

		}
	}
	return nil
}

func (bf *BookForm) reset() {
	_ = bf.Title.Set("")
	_ = bf.Author.Set("")
	_ = bf.Genre.Set("")
	_ = bf.Borrower.Set("")
	_ = bf.Date.Set("")
	_ = bf.Completed.Set("")
	_ = bf.Rating.Set(Ratings()[0])
	_ = bf.IsLoaned.Set(false)
	_ = bf.IsRead.Set(false)
}

func (bf *CreateBookForm) AddSubmission() {
	err := validate(&bf.BookForm)
	if err != nil {
		bf.bus.Notify(bus.Event{
			Name: msgUserError,
			Data: err.Error(),
		})
		return 
	}

	bf.bus.Notify(bus.Event{
		Name: msgUserSuccess,
		Data: "Added",
	})

	bf.sl.addSubmission(&bf.BookForm)
	bf.reset()
}

func (bf *CreateBookForm) SubmissionList() *SubmissionList {
	return bf.sl
}

func (bf *CreateBookForm) Submit() {
	if bf.sl.Length() == 0  {
		bf.bus.Notify(bus.Event{
			Name: msgUserInfo,
			Data: "No Submissions",
		})
	}
}

