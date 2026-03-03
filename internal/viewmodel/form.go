package viewmodel

import (
	"errors"
	"slices"
	"fmt"
	"log"

	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/bus"
	repo "github.com/dubbersthehoser/mayble/internal/repository"

)

type BookForm struct {
	id        int64
	Title     binding.String
	Author    binding.String
	Genre     binding.String

	IsLoaned binding.Bool
	Borrower binding.String
	Date     binding.String

	IsRead    binding.Bool
	Rating    binding.String
	Completed binding.String
}

func NewBookForm() *BookForm {
	return &BookForm{
		id: -1,
		Title: binding.NewString(),
		Author: binding.NewString(),
		Genre: binding.NewString(),

		IsRead:    binding.NewBool(),
		Rating:    binding.NewString(),
		Completed: binding.NewString(),

		IsLoaned: binding.NewBool(),
		Borrower: binding.NewString(),
		Date: binding.NewString(),
	}
}

func (bf *BookForm) Set(book *repo.BookEntry) {
	bf.id = book.ID
	_ = bf.Title.Set(book.Title)
	_ = bf.Author.Set(book.Author)
	_ = bf.Genre.Set(book.Genre)

	if repo.Loaned & book.Variant != 0 {
		_ = bf.IsLoaned.Set(true)
		_ = bf.Date.Set(formatDate(&book.Loaned))
		_ = bf.Borrower.Set(book.Borrower)
	}

	if repo.Read & book.Variant != 0 {
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

func (s *SubmissionList) popSubmission() (*repo.BookEntry, error) {
	if len(s.submissions) == 0 {
		return nil, errors.New("empty list")
	}
	top := s.submissions[len(s.submissions)-1]
	s.submissions = s.submissions[:len(s.submissions)-1]
	s.l.notify()
	return &top, nil
}

func (s *SubmissionList) appendSubmission(book *repo.BookEntry) {
	s.submissions = append(s.submissions, *book)
	s.l.notify()
}

func (s *SubmissionList) addSubmission(bf *BookForm) error {
	//book := repo.BookEntry{}

	//isOnLoan, _ := bf.IsLoaned.Get()
	//isRead, _ := bf.IsRead.Get()

	//if isOnLoan {
	//	book.Variant |= repo.Loaned
	//	date, _ := bf.Date.Get()
	//	timeDate, _ := parseDate(date)
	//	book.Loaned = *timeDate
	//	book.Borrower, _ = bf.Borrower.Get()
	//}

	//if isRead {
	//	book.Variant |= repo.Read
	//	sdate, _ := bf.Completed.Get()
	//	srating, _ := bf.Rating.Get()

	//	date, _ := parseDate(sdate)
	//	rating := slices.Index(Ratings(), srating)

	//	book.Read = *date
	//	book.Rating = rating
	//}

	//book.Variant |= repo.Book

	//book.Title, _ = bf.Title.Get()
	//book.Author, _ = bf.Author.Get()
	//book.Genre, _ = bf.Genre.Get()

	book := bf.ToBookEntry()

	s.submissions = append(s.submissions, *book)
	s.l.notify()
	s.bus.Notify(bus.Event{
		Name: msgUserInfo,
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
	s.form.Set(&book)
	s.Remove(idx)
}





func (s *SubmissionList) AddListener(l binding.DataListener) {
	s.l.AddListener(l)
}

type CreateBookForm struct {
	bus       *bus.Bus
	sl        *SubmissionList
	repo      repo.BookCreator
	Genres    *UniqueGenres
	BookForm
}

func NewCreateBookForm(vms *vmService) *CreateBookForm {
	bf := &CreateBookForm{
		sl: NewSubmissionList(vms.bus),
		bus: vms.bus,
		Genres: vms.genres,
		repo: vms.app.bookCreator,
		BookForm: *NewBookForm(),
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

		println(completed)

		if completed == "" {
			return errors.New("Missing Completion Date")
		}
		_, err := parseDate(completed)
		if err != nil {
			return errors.New("Invalid Completion Date (DD/MM/YYYY)")
		}
		ratings := Ratings()
		rank := slices.Index(ratings[1:], rating)
		if rank == -1 {
			return errors.New("Invalid Rating")

		}
	}
	return nil
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
		Name: msgUserInfo,
		Data: "Added submission",
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
			Data: "No submissions to submit",
		})
		return
	}

	if bf.repo == nil {
		bf.bus.Notify(bus.Event{
			Name: msgUserInfo,
			Data: "Not implemented.",
		})
		return
	}

	failed := make([]repo.BookEntry, 0)
	for {
		book, err := bf.sl.popSubmission()
		if err != nil {
			break
		}
		err = bf.repo.CreateBook(book)
		if err != nil {
			failed = append(failed, *book)
			log.Println(fmt.Errorf("form.Submit: %w", err))
		}
	}
	bf.sl.Clear()
	for _, f := range failed {
		bf.sl.appendSubmission(&f)
	}
	if len(failed) > 0 {
		bf.bus.Notify(bus.Event{
			Name: msgUserError,
			Data: "Submission failed be added.",
		})
		return
	}

	bf.bus.Notify(bus.Event{
		Name: msgUserSuccess,
		Data: "Books Added!",
	})

	bf.bus.Notify(bus.Event{
		Name: msgDataChanged,
	})
}


type EditBookVM struct {
	BookForm
	Genres  *UniqueGenres
	bus     *bus.Bus
	IsOpen  binding.Bool
	vms     *vmService
}

func NewEditBookVM(vms *vmService, isOpen binding.Bool) *EditBookVM {
	ed := &EditBookVM{
		bus:      vms.bus,
		BookForm: *NewBookForm(),
		Genres:   vms.genres,
		IsOpen:   isOpen,
		vms:      vms,
	}
	return ed
}

func (ed *EditBookVM) Submit() {
	err := validate(&ed.BookForm)
	if err != nil {
		ed.bus.Notify(bus.Event{
			Name: msgUserError,
			Data: err.Error(),
		})
		return 
	}

	book := ed.BookForm.ToBookEntry()

	ed.vms.app.bookUpdator.UpdateBook(book)
	
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
