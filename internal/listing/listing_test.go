package listing

import (
	"testing"
	"time"

	"github.com/dubbersthehoser/mayble/internal/app"
)

func TestDateToString(t *testing.T) {
	tests := []struct{
		input time.Time
		expect string
	}{
		{
			input: time.Date(2000, time.Month(3), 5, 0, 0, 0, 0, time.UTC),
			expect: "5/3/2000",
		},
		{
			input: time.Date(1999, time.Month(12), 13, 0, 0, 0, 0, time.UTC),
			expect: "13/12/1999",
		},
	}

	for i, test := range tests {
		expect := test.expect
		actual := DateToString(&test.input)
		if expect != actual {
			t.Fatalf("case %d, expect %s, got %s", i, expect, actual)
		}
	}
}

func TestStringToOrderBy(t *testing.T) {
	tests := []struct{
		input string
		expect OrderBy
	}{
		{
			input: "Title",
			expect: ByTitle,
		},
		{
			input: "AuthoR",
			expect: ByAuthor,
		},
		{
			input: "GeNre",
			expect: ByGenre,
		},
		{
			input: "RaTTING",
			expect: ByRatting,
		},
		{
			input: "borrower",
			expect: ByBorrower,
		},
		{
			input: "DATE",
			expect: ByDate,
		},
		{
			input: "",
			expect: ByNothing,
		},
	}
	for i, test := range tests {
		expect := test.expect
		actual, err := StringToOrderBy(test.input)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		if expect != actual {
			t.Fatalf("case %d, expect %#v, got %#v", i, expect, actual)
		}
	}

	errTests := []string{
		"nothing",
		"borower",
		"tield",
		"data",
		"gener",
		"titel",
	}
	for i, test := range errTests {
		_, err := StringToOrderBy(test)
		if err == nil {
			t.Fatalf("case %d, expected error with %s", i, test)
		}
	}
}

func TestRattingToString(t *testing.T) {
	rattingStrings := GetRattingStrings()
	tests := []struct{
		input  int
		expect string
	}{
		{
			input: 0,
			expect: rattingStrings[0],
		},
		{
			input: 1,
			expect: rattingStrings[1],
		},
		{
			input: 2,
			expect: rattingStrings[2],
		},
		{
			input: 3,
			expect: rattingStrings[3],
		},
		{
			input: 4,
			expect: rattingStrings[4],
		},
		{
			input: 5,
			expect: rattingStrings[5],
		},
	}

	for i, test := range tests {
		expect  := test.expect
		actual, err := RattingToString(test.input)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		if expect != actual {
			t.Fatalf("case %d, expect %s, got %s", i, expect, actual)
		}
	}

	errTests := []int{
		-1,
		-16,
		6,
		10,
		8,
		1234,
		-2134,
	}

	for i, test := range errTests {
		_, err := RattingToString(test)
		if err == nil {
			t.Fatalf("case %d, expected error with %d", i, test)
		}
	}
}

func TestRattingToInt(t *testing.T) {
	rattingStrings := GetRattingStrings()

	tests := []struct{
		input  string
		expect int
	}{
		{
			expect: 0,
			input: rattingStrings[0],
		},
		{
			expect: 1,
			input: rattingStrings[1],
		},
		{
			expect: 2,
			input: rattingStrings[2],
		},
		{
			expect: 3,
			input: rattingStrings[3],
		},
		{
			expect: 4,
			input: rattingStrings[4],
		},
		{
			expect: 5,
			input: rattingStrings[5],
		},
	}

	for i, test := range tests {
		expect  := test.expect
		actual, err := RattingToInt(test.input)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		if expect != actual {
			t.Fatalf("case %d, expect %d, got %d", i, expect, actual)
		}
	}
	errTests := []string{
		"ratting",
		"five",
		"one",
		"star",
		"bird",
	}

	for i, test := range errTests {
		_, err := RattingToInt(test)
		if err == nil {
			t.Fatalf("case %d, expected error with %s", i, test)
		}
	}
}

func TestOrderBookLoans(t *testing.T) {
	type Input struct{
		order Ordering
		by    OrderBy
		items []app.BookLoan
	}
	tests := []struct{
		input  Input
		expect []app.BookLoan
	}{
		{ // by title
			input: Input{
				order: ASC,
				by: ByTitle,
				items: []app.BookLoan{
					app.BookLoan{
						ID: 0,
						Title: "zack",
					},
					app.BookLoan{
						ID: 1,
						Title: "anderson",
					},
					app.BookLoan{
						ID: 2,
						Title: "rake",
					},
				},
			},
			expect: []app.BookLoan{
				app.BookLoan{
					ID: 1,
					Title: "anderson",
				},
				app.BookLoan{
					ID: 2,
					Title: "rake",
				},
				app.BookLoan{
					ID: 0,
					Title: "zack",
				},
			},
		},
		{ // by author
			input: Input{
				order: ASC,
				by: ByAuthor,
				items: []app.BookLoan{
					app.BookLoan{
						ID: 0,
						Author: "zack",
					},
					app.BookLoan{
						ID: 1,
						Author: "anderson",
					},
					app.BookLoan{
						ID: 2,
						Author: "rake",
					},
				},
			},
			expect: []app.BookLoan{
				app.BookLoan{
					ID: 1,
					Author: "anderson",
				},
				app.BookLoan{
					ID: 2,
					Author: "rake",
				},
				app.BookLoan{
					ID: 0,
					Author: "zack",
				},
			},
		},
		{ // by genre
			input: Input{
				order: ASC,
				by: ByGenre,
				items: []app.BookLoan{
					app.BookLoan{
						ID: 0,
						Genre: "zack",
					},
					app.BookLoan{
						ID: 1,
						Genre: "anderson",
					},
					app.BookLoan{
						ID: 2,
						Genre: "rake",
					},
				},
			},
			expect: []app.BookLoan{
				app.BookLoan{
					ID: 1,
					Genre: "anderson",
				},
				app.BookLoan{
					ID: 2,
					Genre: "rake",
				},
				app.BookLoan{
					ID: 0,
					Genre: "zack",
				},
			},
		},
		{ // by borrower
			input: Input{
				order: ASC,
				by: ByBorrower,
				items: []app.BookLoan{
					app.BookLoan{
						ID: 0,
						IsOnLoan: true,
						Borrower: "zack",
					},
					app.BookLoan{
						ID: 1,
						IsOnLoan: true,
						Borrower: "anderson",
					},
					app.BookLoan{
						ID: 2,
						IsOnLoan: true,
						Borrower: "rake",
					},
				},
			},
			expect: []app.BookLoan{
				app.BookLoan{
					ID: 1,
					IsOnLoan: true,
					Borrower: "anderson",
				},
				app.BookLoan{
					ID: 2,
					IsOnLoan: true,
					Borrower: "rake",
				},
				app.BookLoan{
					ID: 0,
					IsOnLoan: true,
					Borrower: "zack",
				},
			},
		},
		{ // by ratting
			input: Input{
				order: ASC,
				by: ByRatting,
				items: []app.BookLoan{
					app.BookLoan{
						ID: 0,
						Ratting: 5,
						Borrower: "zack",
					},
					app.BookLoan{
						ID: 1,
						Ratting: 0,
						Borrower: "anderson",
					},
					app.BookLoan{
						ID: 2,
						Ratting: 3,
						Borrower: "rake",
					},
				},
			},
			expect: []app.BookLoan{
				app.BookLoan{
					ID: 1,
					Borrower: "anderson",
				},
				app.BookLoan{
					ID: 2,
					Borrower: "rake",
				},
				app.BookLoan{
					ID: 0,
					Borrower: "zack",
				},
			},
		},
		{ // by date
			input: Input{
				order: ASC,
				by: ByDate,
				items: []app.BookLoan{
					app.BookLoan{
						ID: 0,
						IsOnLoan: true,
						Date: time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
					},
					app.BookLoan{
						ID: 1,
						IsOnLoan: true,
						Date: time.Date(2010, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
					},
					app.BookLoan{
						ID: 2,
						IsOnLoan: true,
						Date: time.Date(2015, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			expect: []app.BookLoan{
				app.BookLoan{
					ID: 1,
				},
				app.BookLoan{
					ID: 2,
				},
				app.BookLoan{
					ID: 0,
				},
			},
		},
		{ // by date DEC
			input: Input{
				order: DEC,
				by: ByDate,
				items: []app.BookLoan{
					app.BookLoan{
						ID: 0,
						IsOnLoan: true,
						Date: time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
					},
					app.BookLoan{
						ID: 1,
						IsOnLoan: true,
						Date: time.Date(2010, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
					},
					app.BookLoan{
						ID: 2,
						IsOnLoan: true,
						Date: time.Date(2015, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			expect: []app.BookLoan{
				app.BookLoan{
					ID: 0,
				},
				app.BookLoan{
					ID: 2,
				},
				app.BookLoan{
					ID: 1,
				},
			},
		},
		{ // by Nothing
			input: Input{
				order: ASC,
				by: ByNothing,
				items: []app.BookLoan{
					app.BookLoan{
						ID: 0,
						Title: "z",
						Author: "z",
						Genre: "z",
						Ratting: 1,
						IsOnLoan: true,
						Borrower: "z", 
						Date: time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
					},
					app.BookLoan{
						ID: 1,
						Title: "b",
						Author: "b",
						Genre: "b",
						Ratting: 1,
						IsOnLoan: true,
						Borrower: "b", 
						Date: time.Date(2010, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
					},
					app.BookLoan{
						ID: 2,
						Title: "a",
						Author: "a",
						Genre: "a",
						Ratting: 0,
						IsOnLoan: true,
						Borrower: "a", 
						Date: time.Date(2015, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			expect: []app.BookLoan{
				app.BookLoan{
					ID: 0,
				},
				app.BookLoan{
					ID: 1,
				},
				app.BookLoan{
					ID: 2,
				},
			},
		},
		{ // by id
			input: Input{
				order: ASC,
				by: ByID,
				items: []app.BookLoan{
					app.BookLoan{
						ID: 1,
						Title: "z",
						Author: "z",
						Genre: "z",
						Ratting: 1,
						IsOnLoan: true,
						Borrower: "z", 
						Date: time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
					},
					app.BookLoan{
						ID: 2,
						Title: "b",
						Author: "b",
						Genre: "b",
						Ratting: 1,
						IsOnLoan: true,
						Borrower: "b", 
						Date: time.Date(2010, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
					},
					app.BookLoan{
						ID: 0,
						Title: "a",
						Author: "a",
						Genre: "a",
						Ratting: 0,
						IsOnLoan: true,
						Borrower: "a", 
						Date: time.Date(2015, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			expect: []app.BookLoan{
				app.BookLoan{
					ID: 0,
				},
				app.BookLoan{
					ID: 1,
				},
				app.BookLoan{
					ID: 2,
				},
			},
		},
	}

	for i, test := range tests {
		input := test.input
		actual := OrderBookLoans(input.items, input.by, input.order)
		expect := test.expect

		if len(expect) != len(actual) {
			t.Fatalf("case %d, expect length %d, got %d", i, len(expect), len(actual))
		}

		for index := range expect {
			if actual[index].ID != expect[index].ID {
				t.Fatalf("case %d, index %d, expect %d, got %d", i, index, expect[index].ID, actual[index].ID)
			}
		}

	}
}

func TestBookLoanToListed(t *testing.T) {

	date := time.Now()
	
	tests := []struct{
		input app.BookLoan
		expect BookLoan
	}{
		{
			input: app.BookLoan{
				Title: "title_1",
				Author: "author_1",
				Genre: "genre_1",
				Ratting: 1,
				IsOnLoan: false,
				Borrower: "borrower_1",
				Date: date,
			},
			expect: BookLoan{
				Title: "title_1",
				Author: "author_1",
				Genre: "genre_1",
				Ratting: "⭐",
				Borrower: "n/a",
				Date: "n/a",
			},
		},
		{
			input: app.BookLoan{
				Title: "title_2",
				Author: "author_2",
				Genre: "genre_2",
				Ratting: 2,
				IsOnLoan: true,
				Borrower: "borrower_2",
				Date: time.Date(2000, time.Month(11), 10, 0, 0, 0, 0, time.UTC),
			},
			expect: BookLoan{
				Title: "title_2",
				Author: "author_2",
				Genre: "genre_2",
				Ratting: "⭐⭐",
				Borrower: "borrower_2",
				Date: "10/11/2000",
			},
		},
	}

	for i, test := range tests {
		actual := BookLoanToListed(&test.input)
		expect := test.expect
		if actual.Title != expect.Title {
			t.Fatalf("case %d, expect %s, got %s", i, expect.Title, actual.Title)
		}
		if actual.Author != expect.Author {
			t.Fatalf("case %d, expect %s, got %s", i, expect.Author, actual.Author)
		}
		if actual.Genre != expect.Genre {
			t.Fatalf("case %d, expect %s, got %s", i, expect.Genre, actual.Genre)
		}
		if actual.Ratting != expect.Ratting {
			t.Fatalf("case %d, expect %s, got %s", i, expect.Ratting, actual.Ratting)
		}
		if actual.Borrower != expect.Borrower {
			t.Fatalf("case %d, expect %s, got %s", i, expect.Borrower, actual.Borrower)
		}
		if actual.Date != expect.Date {
			t.Fatalf("case %d, expect %s, got %s", i, expect.Date, actual.Date)
		}
	}
}
