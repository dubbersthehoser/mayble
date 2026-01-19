package viewmodel

import (
	"errors"
	"time"

	"fyne.io/fyne/v2/data/binding"
)

type BookForm struct {
	bus       *eventBus
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

	Valid   binding.Bool
	Error   binding.String
	Success binding.String
}


func NewBookForm(err, success binding.String) *BookForm {
	bf := &BookForm{
		bus: &eventBus{},

		Title: binding.NewString(),
		Author: binding.NewString(),
		Genre: binding.NewString(),

		Valid: binding.NewBool(),
		Success: success,
		Error: err,

		IsRead:    binding.NewBool(),
		Rating:    binding.NewString(),
		Completed: binding.NewString(),

		IsLoaned: binding.NewBool(),
		Borrower: binding.NewString(),
		Date: binding.NewString(),
	}
	return bf
}

func (bf *BookForm) validate() error {
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

	dateFormat := "02/01/2006"

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
		_, err := time.Parse(dateFormat, date)
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
		_, err := time.Parse(dateFormat, completed)
		if err != nil {
			return errors.New("Invalid Completion Date (DD/MM/YYYY)")
		}
		// convert rating to int
		_ = rating
	}
	return nil
}

func (bf *BookForm) Submit() {
	bf.bus.publish(event{
		name: "BookFormSubmit",
		data: bf,
	})
	err := bf.validate()
	if err != nil {
		_ = bf.Valid.Set(false)
		_ = bf.Success.Set("")
		_ = bf.Error.Set(err.Error())
		return 
	}

	_ = bf.Valid.Set(true)
	_ = bf.Success.Set("Added")

	_ = bf.Title.Set("")
	_ = bf.Author.Set("")
	_ = bf.Genre.Set("")
	_ = bf.Borrower.Set("")
	_ = bf.Date.Set("")
	_ = bf.Completed.Set("")
	_ = bf.Rating.Set("")
	_ = bf.IsLoaned.Set(false)
	_ = bf.IsRead.Set(false)
}

func (bf *BookForm) Cancel() {
	bf.bus.publish(event{
		name: "BookFormClose",
	})
}
