package validation

import (
	"github.com/dubbersthehoser/mayble/internal/app"
)

type ValidationErr struct {
	field       string
	description string
}
func (v *ValidationErr) Description() string {
	return v.description
}
func (v *ValidationErr) Field() string {
	return v.field
}
func (v *ValidationErr) Error() string {
	return fmt.Sprintf("data: %s, %s", v.field, v.description)
}

func ValidateBookLoan(b *app.BookLoan) error {
	var (
		InvalidTitle        error = ValidTitle(b.Title)
		InvalidAuthor       error = ValidAuthor(b.Author)
		InvalidGenre        error = ValidGenre(b.Genre)
		InvalidRatting      error = ValidRatting(b.Ratting)
		InvalidBorrowerName error
		InvalidBorrowDate   error
	)
	if b.IsOnLoan() {
		InvalidBorrowerName = ValidBorrowerName(b.Loan.Borrower)
		InvalidBorrowDate = ValidBorrowDate(&b.Loan.Date)
	}
	switch {
	case InvalidTitle != nil:
		return InvalidTitle

	case InvalidAuthor != nil:
		return InvalidAuthor

	case InvalidGenre != nil:
		return InvalidGenre

	case InvalidRatting != nil:
		return InvalidRatting

	case b.IsOnLoan() && InvalidBorrowerName != nil:
		return InvalidBorrowerName

	case b.IsOnLoan() && InvalidBorrowDate != nil:
		return InvalidBorrowDate
	}
	return nil
}

func ValidID(id int64) error {
	err := &ValidationErr{}
	if id < 0 {
		err.field= "ID"
		err.description = "can't be a negative value"
		return err
	}
	return nil
}
func ValidTitle(s string) error {
	err := &ValidationErr{}
	if s == ""  {
		err.field= "Title"
		err.description = "can't be empty string"
		return err
	}
	return nil
}
func ValidAuthor(s string) error {
	err := &ValidationErr{}
	if s == ""  {
		err.field= "Author"
		err.description = "can't be empty string"
		return err
	}
	return nil
}
func ValidGenre(s string) error {
	err := &ValidationErr{}
	if s == ""  {
		err.field = "Genre"
		err.description = "can't be empty string"
		return err
	}
	return nil
}
func ValidRatting(ratting int) error {
	err := &ValidationErr{}
	if ratting < 0 || ratting >= 6 {
		err.field = "Ratting"
		err.description = "out of range of 0-5"
		return err
	}
	return nil
}
func ValidBorrowerName(s string) error {
	err := &ValidationErr{}
	if s == ""  {
		err.field = "Borrower"
		err.description = "can't be empty string"
		return err
	}
	return nil
}
func ValidBorrowDate(date *time.Time) error {
	err := &ValidationErr{}
	if date == nil {
		err.field = "Date"
		err.description = "can't be nil value"
		return err
	}
	if date.IsZero() {
		err.field = "Date"
		err.description = "can't be zero value"
	}
	return nil
}
