package viewmodel


import (
	"time"
	"slices"
	"testing"

	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

// TODO fix and re-write tests

func TestVariantToTableName(t *testing.T) {
	tests := []struct{
		v repo.Variant
		expect string
	}{
		{v: repo.Book, expect: "All Books"},
		{v: repo.Book | repo.Loaned, expect: "On Loaned"},
		{v: repo.Book | repo.Read, expect: "Read"},
	}

	for i, c := range tests {
		actual := VariantToTableName(c.v)
		if c.expect != actual {
			t.Fatalf("case[%d] expect %s, got %s", i, c.expect, actual)
		}
	}
}


func TestEntryValues(t *testing.T) {
	entry := repo.BookEntry{
		Title: "Expect Title",
		Author: "Expect Author",
		Genre: "Expect Genre",

		Borrower: "Expect Borrower",
		Loaned: time.Date(2020, 11, 01, 0, 0, 0, 0, time.UTC),

		Rating: 3,
		Read: time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC),
	}
	
	tests := []struct{
		v      repo.Variant
		expect []string
	}{
		{
			v: repo.Book,
			expect: []string{
				entry.Title,
				entry.Author,
				entry.Genre,
			},
		},
		{
			v: repo.Book | repo.Loaned,
			expect: []string{
				entry.Title,
				entry.Author,
				entry.Genre,
				entry.Borrower,
				formatDate(&entry.Loaned),
			},
		},
		{
			v: repo.Book | repo.Read,
			expect: []string{
				entry.Title,
				entry.Author,
				entry.Genre,
				formatRating(entry.Rating),
				formatDate(&entry.Read),
			},
		},
		{
			v: repo.Book | repo.Read | repo.Loaned,
			expect: []string{},
		},
	}

	for i, c := range tests {
		entry.Variant = c.v
		actual := EntryValues(&entry)
		if r := slices.Compare(c.expect, actual); r != 0 {
			t.Fatalf("[%d] expect\n\t%#v\ngot\n\t%#v", i, c.expect, actual)
		}
	}
}


func TestVariantFields(t *testing.T) {
	tests := []struct{
		v      repo.Variant
		expect []string
	}{
		{
			v: repo.Book,
			expect: []string{
				"Title",
				"Author",
				"Genre",
			},
		},
	}

	for i, c := range tests {
		actual := VariantFields(c.v)
		if r := slices.Compare(c.expect, actual); r != 0 {
			t.Fatalf("[%d] expect\n\t%#vgot\n\t%#v", i, c.expect, actual)
		}
	}
}
