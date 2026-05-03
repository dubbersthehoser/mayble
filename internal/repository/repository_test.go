package repository

import (
	"testing"
)

func TestHeaderIndexing(t *testing.T) {
	
	headers := BookEntryFields()

	if len(headers) != 7 {
		t.Fatalf("expect %d, got %d", 7, len(headers))
	}

	tests := []struct{
		index int
		expect string
	}{
		{
			index: 0,
			expect: "Title",
		},
		{
			index: 1,
			expect: "Author",
		},
		{
			index: 2,
			expect: "Genre",
		},
		{
			index: 3,
			expect: "Rating",
		},
		{
			index: 4,
			expect: "Read",
		},
		{
			index: 5,
			expect: "Borrower",
		},
		{
			index: 6,
			expect: "Loaned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expect, func(t *testing.T) {
			if headers[tt.index] != tt.expect {
				t.Fatalf("expect %s, got %s", tt.expect, headers[tt.index])
			}
		})
	}

}
