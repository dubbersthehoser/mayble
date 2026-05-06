package viewmodel

import (
	"errors"
	"fmt"
	"log"
	"slices"

	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/viewmodel/submissions"
	"github.com/dubbersthehoser/mayble/internal/models"
	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

// A BookForm the view model for book editing and creation.
type BookForm struct {
	id     int64
	Title  binding.String
	Author binding.String
	Genre  binding.String

	IsLoaned binding.Bool
	Borrower binding.String
	Date     binding.String

	IsRead    binding.Bool
	Rating    binding.String
	Completed binding.String
}

func NewBookForm() *BookForm {
	return &BookForm{
		id:     -1,
		Title:  binding.NewString(),
		Author: binding.NewString(),
		Genre:  binding.NewString(),

		IsRead:    binding.NewBool(),
		Rating:    binding.NewString(),
		Completed: binding.NewString(),

		IsLoaned: binding.NewBool(),
		Borrower: binding.NewString(),
		Date:     binding.NewString(),
	}
}

func (bf *BookForm) reset() {
	bf.id = -1
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

func (bf *BookForm) validate() error {
	title, _ := bf.Title.Get()
	author, _ := bf.Author.Get()
	genre, _ := bf.Genre.Get()

	if title == "" {
		return errors.New("missing title")
	}
	if author == "" {
		return errors.New("missing auther")
	}
	if genre == "" {
		return errors.New("missing genre")
	}

	isLoaned, _ := bf.IsLoaned.Get()
	isRead, _ := bf.IsRead.Get()

	if isLoaned {
		date, _ := bf.Date.Get()
		borrower, _ := bf.Borrower.Get()

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
		completed, _ := bf.Completed.Get()
		rating, _ := bf.Rating.Get()

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

func (bf *BookForm) Set(book *models.BookEntry) {
	bf.reset()
	bf.id = book.ID
	_ = bf.Title.Set(book.Title)
	_ = bf.Author.Set(book.Author)
	_ = bf.Genre.Set(book.Genre)

	if book.IsLoaned {
		_ = bf.IsLoaned.Set(true)
		_ = bf.Date.Set(formatDate(&book.LoanedAt))
		_ = bf.Borrower.Set(book.Borrower)
	}

	if book.IsCompleted {
		_ = bf.IsRead.Set(true)
		_ = bf.Rating.Set(formatRating(book.Rating))
		_ = bf.Completed.Set(formatDate(&book.CompletedAt))
	}
}

func (bf *BookForm) ToBookEntry() *models.BookEntry {
	book := &models.BookEntry{ID: bf.id}

	isOnLoan, _ := bf.IsLoaned.Get()
	isRead, _ := bf.IsRead.Get()

	if isOnLoan {
		book.IsLoaned = true
		date, _ := bf.Date.Get()
		timeDate, _ := parseDate(date)
		book.LoanedAt = *timeDate
		book.Borrower, _ = bf.Borrower.Get()
	}

	if isRead {
		book.IsCompleted = true
		sdate, _ := bf.Completed.Get()
		srating, _ := bf.Rating.Get()

		date, _ := parseDate(sdate)
		rating := slices.Index(Ratings(), srating)

		book.CompletedAt = *date
		book.Rating = rating
	}

	book.Title, _ = bf.Title.Get()
	book.Author, _ = bf.Author.Get()
	book.Genre, _ = bf.Genre.Get()
	return book
}

// A SubmissionList view model of list of form submissions.
type SubmissionList struct {
	form        *BookForm
	bus         *bus.Bus
	l           *listener
	sub         *submissions.List
	Limit       binding.String
}

func fmtFormLimit(count, max int) string {
	return fmt.Sprintf("Limit %d/%d", count, max)
}

func NewSubmissionList(bus *bus.Bus, form *BookForm) *SubmissionList {

	FormLimit := 16

	sl := &SubmissionList{
		l:           &listener{},
		bus:         bus,
		sub:         submissions.NewList(FormLimit),
		form:        form,
		Limit:       binding.NewString(),
	}
	sl.updateLimit()
	return sl
}

func (s *SubmissionList) updateLimit() {
	_ = s.Limit.Set(fmtFormLimit(s.sub.Length(), s.sub.Cap()))
}

func (s *SubmissionList) Clear() {
	s.sub.Clear()
	s.updateLimit()
}

func (s *SubmissionList) pop() (*models.BookEntry, error) {
	top := s.sub.Pop()
	if top == nil {
		return top, errors.New("empty submission list")
	}
	s.updateLimit()
	s.l.notify()
	return top, nil
}

func (s *SubmissionList) add(bf *BookForm) error {
	book := bf.ToBookEntry()
	err := s.sub.Append(*book)
	if err != nil {
		return err
	}
	s.updateLimit()
	s.l.notify()
	return nil
}

func (s *SubmissionList) Length() int {
	return s.sub.Length()
}

func (s *SubmissionList) GetView(idx int) string {
	book, err := s.sub.Get(idx)
	if err != nil {
		return "STUB: out of range"
	}
	prefix := fmt.Sprintf("%s, %s, %s", book.Title, book.Author, book.Genre)
	switch {
	case book.IsLoaned && book.IsCompleted:
		return fmt.Sprintf("%s (loaned) (completed)", prefix)
	case book.IsLoaned:
		return fmt.Sprintf("%s (loaned)", prefix)
	case book.IsCompleted:
		return fmt.Sprintf("%s (completed)", prefix)
	default:
		return prefix
	}
}

func (s *SubmissionList) Remove(idx int) {
	s.bus.Notify(bus.Event{
		Name: msgUserInfo,
		Data: "Form removed",
	})
	err := s.sub.Remove(idx)
	if err == nil {
		s.l.notify()
	}
}

func (s *SubmissionList) Edit(idx int) {
	if s.form == nil {
		log.Println("submission_list.edit: nil form")
		return
	}
	book, err := s.sub.Get(idx)
	if err != nil {
		log.Println("submission_list.edit: index out of range")
		return
	}
	s.form.Set(book)
	s.Remove(idx)
}

func (s *SubmissionList) AddListener(l binding.DataListener) {
	s.l.AddListener(l)
}


// A BookSubmissionForm is the view model of submission page.
type BookSubmissionForm struct {
	bus    *bus.Bus
	sl     *SubmissionList
	repo   repo.BookCreator
	Genres *UniqueGenres
	BookForm
}

func NewBookSubmissionForm(b *bus.Bus, c repo.BookCreator, g *UniqueGenres) *BookSubmissionForm {
	bf := &BookSubmissionForm{
		bus:      b,
		Genres:   g,
		repo:     c,
		BookForm: *NewBookForm(),
	}
	bf.sl = NewSubmissionList(b, &bf.BookForm)
	return bf
}

func (bf *BookSubmissionForm) AddSubmission() {
	err := bf.BookForm.validate()
	if err != nil {
		bf.bus.Notify(bus.Event{
			Name: msgUserError,
			Data: err.Error(),
		})
		return
	}

	err = bf.sl.add(&bf.BookForm)
	if err != nil {
		bf.bus.Notify(bus.Event{
			Name: msgUserError,
			Data: "At form limit. Please submit",
		})
		return 
	}
	

	bf.bus.Notify(bus.Event{
		Name: msgUserInfo,
		Data: "Added Form",
	})
	bf.reset()
}

func (bf *BookSubmissionForm) SubmissionList() *SubmissionList {
	return bf.sl
}

func (bf *BookSubmissionForm) Submit() {

	if bf.repo == nil {
		bf.bus.Notify(bus.Event{
			Name: msgUserInfo,
			Data: "Not implemented.",
		})
		return
	}

	if bf.sl.Length() == 0 {
		bf.bus.Notify(bus.Event{
			Name: msgUserInfo,
			Data: "No submissions to submit",
		})
		return
	}

	for {
		book, err := bf.sl.pop()
		if err != nil {
			break
		}
		_, err = bf.repo.CreateBook(book)
		if err != nil {
			log.Println(fmt.Errorf("form.Submit: %w", err))
			bf.bus.Notify(bus.Event{
				Name: msgUserError,
				Data: "Submission failed",
			})
			return
		}
	}
	bf.sl.Clear()
	bf.bus.Notify(bus.Event{
		Name: msgUserSuccess,
		Data: "Books Added!",
	})

	bf.bus.Notify(bus.Event{
		Name: msgDataChanged,
	})
}

// A EditBookVM is a entry edit form view model.
type EditBookVM struct {
	BookForm
	bus     *bus.Bus
	IsOpen  binding.Bool
	updator repo.BookUpdator
	Genres  *UniqueGenres
}

func NewEditBookVM(b *bus.Bus, u repo.BookUpdator, g *UniqueGenres,  isOpen binding.Bool) *EditBookVM {
	ed := &EditBookVM{
		bus:      b,
		BookForm: *NewBookForm(),
		updator:  u,
		IsOpen:   isOpen,
		Genres:   g,
	}
	return ed
}

func (ed *EditBookVM) Submit() {
	err := ed.BookForm.validate()
	if err != nil {
		ed.bus.Notify(bus.Event{
			Name: msgUserError,
			Data: err.Error(),
		})
		return
	}

	book := ed.BookForm.ToBookEntry()

	ed.updator.UpdateBook(book)

	ed.bus.Notify(bus.Event{
		Name: msgDataChanged,
	})
	ed.Close()
}

func (ed *EditBookVM) Set(b *models.BookEntry) {
	ed.BookForm.Set(b)
}

func (ed *EditBookVM) Close() {
	ed.BookForm.reset()
	ed.IsOpen.Set(false)
}
