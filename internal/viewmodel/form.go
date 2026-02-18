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

	UniqueGenres binding.StringList

	IsLoaned binding.Bool
	Borrower binding.String
	Date     binding.String

	IsRead    binding.Bool
	Rating    binding.String
	Completed binding.String
}
func NewBookForm() *BookForm {
	return &BookForm{
		Title: binding.NewString(),
		Author: binding.NewString(),
		Genre: binding.NewString(),
		UniqueGenres: binding.NewStringList(),

		IsRead:    binding.NewBool(),
		Rating:    binding.NewString(),
		Completed: binding.NewString(),

		IsLoaned: binding.NewBool(),
		Borrower: binding.NewString(),
		Date: binding.NewString(),
	}
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
	case repo.BookReadAndLoaned:
		return fmt.Sprintf("%s (loaned) (read)", prefix)
	case repo.BookLoaned:
		return fmt.Sprintf("%s (loaned)", prefix)
	case repo.BookRead:
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


func (bf *BookForm) Set(book *repo.BookEntry) {
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



func (s *SubmissionList) AddListener(l binding.DataListener) {
	s.l.AddListener(l)
}

type CreateBookForm struct {
	bus       *bus.Bus
	sl        *SubmissionList
	repo      repo.BookCreator
	genres     repo.GenreRetriever
	BookForm
}

func NewCreateBookForm(b *bus.Bus, repo repo.BookCreator, genres repo.GenreRetriever) *CreateBookForm {
	bf := &CreateBookForm{
		sl: NewSubmissionList(b),
		bus: b,
		repo: repo,
		genres: genres, 
		BookForm: *NewBookForm(),
	}

	updateGenres := func() {
		genres, err := genres.GetUniqueGenres()
		if err != nil {
			return 
		}
		for i := range bf.BookForm.UniqueGenres.Length() {
			v, _ := bf.BookForm.UniqueGenres.GetValue(i)
			_ = bf.BookForm.UniqueGenres.Remove(v)
		}
		for i := range genres {
			_ = bf.BookForm.UniqueGenres.Append(genres[i])
		}

	}

	b.Subscribe(bus.Handler{
		Name: msgDataChanged,
		Handler: func(e *bus.Event) {
			updateGenres()
		},
	})

	updateGenres()

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
		Name: msgUserInfo,
		Data: "Added submission",
	})

	bf.sl.addSubmission(&bf.BookForm)
	bf.reset()
}

func (bf *CreateBookForm) SubmissionList() *SubmissionList {
	return bf.sl
}

func (bf *CreateBookForm) GetUniqueGenres() []string {
	genres, err := bf.genres.GetUniqueGenres()
	if err != nil {
		return []string{
			"_STUB_",
		}
	}
	return genres
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
		}
	}
	bf.sl.Clear()
	bf.sl.submissions = failed
	for _, f := range failed {
		bf.sl.appendSubmission(&f)
	}
	if len(failed) > 0 {
		bf.bus.Notify(bus.Event{
			Name: msgUserError,
			Data: "Submission failed be added.",
		})
	} else {
		bf.bus.Notify(bus.Event{
			Name: msgUserSuccess,
			Data: "Books Added!",
		})
	}
}


type UniqueGenres struct {
	list binding.StringList
	genres repo.GenreRetriever
}

func NewUniqueGenres(b *bus.Bus, l binding.StringList, g repo.GenreRetriever) *UniqueGenres {
	ug := &UniqueGenres{
		list: l,
		genres: g,
	}
	b.Subscribe(bus.Handler{
		Name: msgDataChanged,
		Handler: func(e *bus.Event) {
			ug.Update()
		},
	})
	return ug
}

func (u *UniqueGenres) Update() {
	genres, err := u.genres.GetUniqueGenres()
	if err != nil {
		return 
	}
	for i := range u.list.Length() {
		v, _ := u.list.GetValue(i)
		_ = u.list.Remove(v)
	}
	for i := range genres {
		_ = u.list.Append(genres[i])
	}

}


type EditBookVM struct {
	BookForm
	genres    *UniqueGenres
	bus        *bus.Bus
	IsOpen     binding.Bool

}

func NewEditBookVM(b *bus.Bus, isOpen binding.Bool, repo repo.BookUpdator, g repo.GenreRetriever) *EditBookVM {
	ed := &EditBookVM{
		bus: b,
		BookForm: *NewBookForm(),
		IsOpen: isOpen,
	}
	ed.genres = NewUniqueGenres(b, ed.UniqueGenres, g)
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
	ed.bus.Notify(bus.Event{
		Name: msgUserInfo,
		Data: "not implemented",
	})
}

func (ed *EditBookVM) Close() {
	ed.BookForm.reset()
	ed.IsOpen.Set(false)
}
