package viewmodel

import (
	"errors"
	"fmt"
	"log"
	"slices"

	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/viewmodel/submissions"
	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

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
		rank := slices.Index(ratings[1:], rating)
		if rank == -1 {
			return errors.New("invalid rating")

		}
	}
	return nil
}

func (bf *BookForm) Set(book *repo.BookEntry) {
	bf.reset()
	bf.id = book.ID
	_ = bf.Title.Set(book.Title)
	_ = bf.Author.Set(book.Author)
	_ = bf.Genre.Set(book.Genre)

	if repo.Loaned&book.Variant != 0 {
		_ = bf.IsLoaned.Set(true)
		_ = bf.Date.Set(formatDate(&book.Loaned))
		_ = bf.Borrower.Set(book.Borrower)
	}

	if repo.Read&book.Variant != 0 {
		_ = bf.IsRead.Set(true)
		_ = bf.Rating.Set(formatRating(book.Rating))
		_ = bf.Completed.Set(formatDate(&book.Read))
	}
}

func (bf *BookForm) ToBookEntry() *repo.BookEntry {
	book := &repo.BookEntry{ID: bf.id}

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

	book.Title, _ = bf.Title.Get()
	book.Author, _ = bf.Author.Get()
	book.Genre, _ = bf.Genre.Get()
	return book
}

type SubmissionList struct {
	form        *BookForm
	bus         *bus.Bus
	l           *listener
	sub         *submissions.List
	Limit       binding.String
}

func fmtFormLimit(count, max int) string {
	return fmt.Sprintf("Form Limit %d/%d", count, max)
}

func NewSubmissionList(bus *bus.Bus, form *BookForm) *SubmissionList {
	sl := &SubmissionList{
		l:           &listener{},
		bus:         bus,
		sub:         submissions.NewList(25),
		form:        form,
		Limit:       binding.NewString(),
	}
	return sl
}

func (s *SubmissionList) Clear() {
	s.sub.Clear()
}

func (s *SubmissionList) pop() (*repo.BookEntry, error) {
	top := s.sub.Pop()
	if top == nil {
		return top, errors.New("empty submission list")
	}
	s.l.notify()
	return top, nil
}

func (s *SubmissionList) append(book *repo.BookEntry) {
	s.sub.Append(*book)
	s.l.notify()
}

func (s *SubmissionList) add(bf *BookForm) error {
	book := bf.ToBookEntry()
	s.append(book)
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
	prefix := fmt.Sprintf("\"%s\" \"%s\" \"%s\"", book.Title, book.Author, book.Genre)
	switch book.Variant {
	case repo.Read | repo.Loaned:
		return fmt.Sprintf("%s (loaned) (read)", prefix)
	case repo.Loaned:
		return fmt.Sprintf("%s (loaned)", prefix)
	case repo.Read:
		return fmt.Sprintf("%s (read)", prefix)
	case repo.Book:
		return prefix
	default:
		return "STUB: variant not found"
	}
}

func (s *SubmissionList) Remove(idx int) {
	s.bus.Notify(bus.Event{
		Name: msgUserInfo,
		Data: "Submission removed",
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

	bf.bus.Notify(bus.Event{
		Name: msgUserInfo,
		Data: "Added submission",
	})
	bf.sl.add(&bf.BookForm)
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

func (ed *EditBookVM) Set(b *repo.BookEntry) {
	ed.BookForm.Set(b)
}

func (ed *EditBookVM) Close() {
	ed.BookForm.reset()
	ed.IsOpen.Set(false)
}
