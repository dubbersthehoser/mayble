package controller

import (
	"time"
	"errors"
)


// TODO move validation to a core, or app package

func ValidateTitle(title string) error {
	if title == "" {
		return errors.New("must have an title")
	}
	return nil
}
func ValidateAuthor(author string) error {
	if author == "" {
		return errors.New("must have an author")
	}
	return nil
}
func ValidateGenre(genre string) error {
	if genre == "" {
		return errors.New("must have an genre")
	}
	return nil
}
func ValidateLoanName(name string) error {
	if name == "" {
		return errors.New("must have loan name")
	}
	return nil
}
func ValidateLoanDate(date *time.Time) error {
	if date == nil || date.IsZero() {
		return errors.New("must have loan date")
	}
	return nil
}

