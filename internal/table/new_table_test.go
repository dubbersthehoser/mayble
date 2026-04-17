package table

import (
	"testing"
	"slices"
)

func Test_nTable(t *testing.T) {

	headers := []string{
		"Title",
		"Author",
		"Genre",
		"Rating",
		"Read",
	}
	
	table := NewnTable("TEST", headers)

	{
		expect := headers
		actual := table.Headers()

		if !slices.Equal(expect, actual) {
			t.Fatalf("expected\n%#v\n  got\n%#v", expect, actual)
		}
	}

	rowValues := [][]string{
		{
			"title_0",
			"author_0",
			"genre_0",
			"rating_0",
			"read_0",
		},
		{
			"title_1",
			"author_1",
			"genre_1",
			"rating_1",
			"read_1",
		},
		{
			"title_2",
			"author_2",
			"genre_2",
			"rating_2",
			"read_2",
		},
	}

	for id, values := range rowValues {
		table.AppendRow(int64(id), values)
	}

	{
		expect := "title_2"
		actual := table.GetCell(2, 0).value
		if expect != actual {
			t.Fatalf("expect '%s', got '%s'", expect, actual)
		}

		expect = "STUB"
		actual = table.GetCell(100, 100).value
		if expect != actual {
			t.Fatalf("expect '%s', got '%s'", expect, actual)
		}
	}
}


