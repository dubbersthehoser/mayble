package controller

import (
	"fmt"
	"time"
	"errors"
)

func GetRattingStrings() []string {
	return []string{"TBR", "⭐", "⭐⭐", "⭐⭐⭐", "⭐⭐⭐⭐", "⭐⭐⭐⭐⭐"}
}

func RattingToInt(ratting string) int {
	for i, str := range GetRattingStrings() {
		if str == ratting {
			return i
		}
	}
	panic("invalid ratting string was passed")
}

func rattingToString(i int) string {
	str := GetRattingStrings()[i]
	return str
}

func dateToString(date *time.Time) string {
	return fmt.Sprintf("%d/%d/%d", date.Day(), date.Month(), date.Year())
}

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
	if date == nil {
		return errors.New("must have loan date")
	}
	return nil
}

