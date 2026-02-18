package viewmodel


import (
	"time"
	"slices"
	"testing"

	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

func unexpectedError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}


func TextCellPool(t *testing.T) {
	cells := newCellPool()

	idx := cells.create(cellTable)

	if cells.get(idx).kind != cellTable {
		t.Fatalf("expect kind %d, got %d", cellTable, cells.get(idx).kind)
	}

	cells.get(idx).table = "TableA"

	if cells.get(idx).table != "TableA" {
		t.Fatalf("expect %s, got %s", "TableA", cells.get(idx).table)
	}
}


func TestTable(t *testing.T) {
	
	name := "TableA"
	header := []string{
		"Title",
		"Author",
		"Genre",
	}

	table := newTable(
		name,
		header,
	)

	//
	// Check headers()
	//
	if c := slices.Compare(header, table.headers()); c != 0 {
		t.Fatalf("expect\n\t%#v\ngot\n\t%#v\n", header, table.headers())
	}

	
	//	
	// Check append rows
	//
	tests := []struct{
		 	header []string
			values []string
		}{
			{
				header: header,
				values: []string{
					"example title 1",
					"example author 1",
					"example genre 1",
				},
			},
			{
				header: header,
				values: []string{
					"example title 2",
					"example author 2",
					"example genre 2",
				},
			},
			{
				header: header,
				values: []string{
					"example title 3",
					"example author 3",
					"example genre 3",
				},
			},
			{
				header: header,
				values: []string{
					"example title 4",
					"example author 4",
					"example genre 4",
				},
			},
	}

	for i, c := range tests {
		err := table.appendRow(int64(i), c.values)
		unexpectedError(t, err)
	}
	//
	// Check size
	//
	erow, ecol := 4, 3
	row, col := table.size()
	if row != erow {
		t.Fatalf("expect %d, got %d", erow, row)
	}
	if col != ecol {
		t.Fatalf("expect %d, got %d", ecol, col)
	}

	//
	// Check Value
	//
	cell := table.getCell(2, 1)
	v := table.getValue(cell)
	ev := tests[2].values[1]
	if v != ev {
		t.Fatalf("expect %s, got %s", ev, v)
	}

	//	
	// Check ID
	//
	id, err := table.getID(cell)
	unexpectedError(t, err)
	eid := int64(2)
	if id != eid {
		t.Fatalf("expect %d, got %d", eid, id)
	}

	//
	// Check hide header
	//
	hide := []string{"Author"}
	table.setHidden(hide)
	eh := []string{"Title", "Genre", "Author"}
	if r := slices.Compare(eh, table.headers()); r != 0 {
		t.Fatalf("expect\n\t%#v\ngot\n\t%#v", eh, table.headers())
	}
	ok := table.isHidden(cell)
	unexpectedError(t, err)
	if !ok {
		t.Fatalf("expect %t, got %t", true, ok)
	}
	hide = []string{"Author", "Title"}
	table.setHidden(hide)
	eh = []string{"Genre", "Title", "Author", }
	if r := slices.Compare(eh, table.headers()); r != 0 {
		t.Fatalf("expect\n\t%#v\ngot\n\t%#v", eh, table.headers())
	}
	hide = []string{"Author"}
	table.setHidden(hide)
	eh = []string{"Title", "Genre", "Author"}
	if r := slices.Compare(eh, table.headers()); r != 0 {
		t.Fatalf("expect\n\t%#v\ngot\n\t%#v", eh, table.headers())
	}
	hide = []string{}
	table.setHidden(hide)
	eh = []string{"Title", "Author", "Genre"}
	if r := slices.Compare(eh, table.headers()); r != 0 {
		t.Fatalf("expect\n\t%#v\ngot\n\t%#v", eh, table.headers())
	}


	//
	// Check Clear
	//
	table.clearValues()
	erow, ecol = 0, 3
	row, col = table.size()
	if row != erow {
		t.Fatalf("expect %d, got %d", erow, row)
	}
	if col != ecol {
		t.Fatalf("expect %d, got %d", ecol, col)
	}
	
}



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
