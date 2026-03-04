package database

import (
	"testing"
	"time"

	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

func Test_bookBuilder(t *testing.T) {
	
	t.Run("invalid builds", func(t *testing.T) {
		builder := newBookBuilder()
		builder.SetID(432)
		builder.SetTitle("A Title")
		builder.SetAuthor("A Author")
		builder.SetGenre("A Genre")
		_, err := builder.Build()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		builder.SetID(-127)
		_, err = builder.Build()
		if err == nil {
			t.Fatalf("expected error")
		}

		builder.SetID(324)
		builder.SetTitle("")
		_, err = builder.Build()
		if err == nil {
			t.Fatalf("expected error")
		}
		builder.SetTitle("A Title")
		builder.SetAuthor("")
		_, err = builder.Build()
		if err == nil {
			t.Fatalf("expected error")
		}
		builder.SetAuthor("A Author")
		builder.SetGenre("")
		_, err = builder.Build()
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("build", func(t *testing.T) {
		testBuilderBuild(t)
	})
}

func testBuilderBuild(t *testing.T) {
	type input struct {
			id       int64
			title,
			author,
			genre,
			loanDate, 
			borrower, 
			readDate  string
			rating    int64
	}
	tests := []struct{
		name   string
		input  input
		expect repo.BookEntry
	}{
		{
			name: "just a book",
			input: input{
				id: 132,
				title: "A Title",
				author: "A Author",
				genre: "A Genre",
			},
			expect: repo.BookEntry{
				ID: 132,
				Title: "A Title",
				Author: "A Author",
				Genre: "A Genre",
			},
		},
		{
			name: "just loaned book",
			input: input{
				id: 132,
				title: "A Title",
				author: "A Author",
				genre: "A Genre",
				loanDate: "2020-02-03",
				borrower: "Lane",
			},
			expect: repo.BookEntry{
				Variant: repo.Loaned,
				ID: 132,
				Title: "A Title",
				Author: "A Author",
				Genre: "A Genre",
				Loaned: time.Date(2020, time.February, 3, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
			},
		},
		{
			name: "just read book",
			input: input{
				id: 132,
				title: "A Title",
				author: "A Author",
				genre: "A Genre",
				readDate: "2020-02-03",
				rating: 3,
			},
			expect: repo.BookEntry{
				Variant: repo.Read,
				ID: 132,
				
				Title: "A Title",
				Author: "A Author",
				Genre: "A Genre",
				Read: time.Date(2020, time.February, 3, 0, 0, 0, 0, time.UTC),
				Rating: 3,
			},
		},
		{
			name: "just book and read loaned",
			input: input{
				id: 132,
				title: "A Title",
				author: "A Author",
				genre: "A Genre",
				readDate: "2020-02-03",
				rating: 3,
				loanDate: "2020-02-03",
				borrower: "Lane",
			},
			expect: repo.BookEntry{
				Variant: repo.Read | repo.Loaned,
				ID: 132,
				Title: "A Title",
				Author: "A Author",
				Genre: "A Genre",
				Read: time.Date(2020, time.February, 3, 0, 0, 0, 0, time.UTC),
				Rating: 3,
				Loaned: time.Date(2020, time.February, 3, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
			},
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			builder := newBookBuilder()
			actual, err := builder.SetID(c.input.id).
				SetTitle(c.input.title).
				SetAuthor(c.input.author).
				SetGenre(c.input.genre).
				SetReadDate(c.input.readDate).
				SetRating(c.input.rating).
				SetLoanedDate(c.input.loanDate).
				SetBorrower(c.input.borrower).
				Build()

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if actual == nil {
				t.Fatalf("unexpected nil")
			}

			if *actual != c.expect {
				t.Fatalf("expect\n%v\n  got\n%v", c.expect, *actual)
			}
		})
	}
}
