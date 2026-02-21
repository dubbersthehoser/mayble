package table

import (
	"testing"
	"slices"
	"fmt"
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
		"Read",
		"Rating",
		"Borrower",
		"Loaned",
	}

	const RowCount = 12

	table := NewTable(
		name,
		header,
	)

	//
	// Check headers()
	//
	if c := slices.Compare(header, table.BaseHeaders()); c != 0 {
		t.Fatalf("expect\n\t%#v\ngot\n\t%#v\n", header, table.BaseHeaders())
	}

	
	//	
	// Check append rows
	//
	type Test  struct{
		header []string
		values []string
	}
	tests := []Test{}
	valuesTempl := []string{
		"example title",
		"example author",
		"example genre",
		"example read",
		"example rating",
		"example borrower",
		"example loaned",
	}
	for range RowCount {
		entry := Test{
			header: header,
			values: make([]string, len(header)),
		}
		for i, templ := range valuesTempl {
			entry.values[i] = fmt.Sprintf("%s %d", templ, i)
		}
		tests = append(tests, entry)
	}

	for i, c := range tests {
		err := table.AppendRow(int64(i), c.values)
		unexpectedError(t, err)
	}
	//
	// Check size
	//
	erow, ecol := RowCount, len(header)
	row, col := table.Size()
	if row != erow {
		t.Fatalf("expect %d, got %d", erow, row)
	}
	if col != ecol {
		t.Fatalf("expect %d, got %d", ecol, col)
	}

	//
	// Check Value
	//
	cell := table.GetCell(2, 1)
	v := cell.Value()
	ev := tests[2].values[1]
	if v != ev {
		t.Fatalf("expect %s, got %s", ev, v)
	}

	//	
	// Check ID
	//
	id := cell.ID()
	eid := int64(2)
	if id != eid {
		t.Fatalf("expect %d, got %d", eid, id)
	}

	//
	// Check hide header
	//

	hideTests := []struct{
		input  []string
		expect []string
	}{
		{
			input: []string{
			},
			expect: []string{
				"Title",
				"Author",
				"Genre",
				"Read",
				"Rating",
				"Borrower",
				"Loaned",
			},
		},
		{
			input: []string{
				"Author",
			},
			expect: []string{
				"Title",
				"Genre",
				"Read",
				"Rating",
				"Borrower",
				"Loaned",
				"Author",
			},
		},
		{
			input: []string{
				"Title",
				"Author",
				"Read",
			},
			expect: []string{
				"Genre",
				"Rating",
				"Borrower",
				"Loaned",
				"Title",
				"Read",
				"Author",
			},
		},
	}
	for i, c := range hideTests {
		table.SetHidden(c.input)
		if r := slices.Compare(c.expect, table.Headers()); r != 0 {
			t.Fatalf("[%d] expect\n\t%#v\ngot\n\t%#v", i, c.expect, table.Headers())
		}
	}

	//
	// Check Clear
	//
	table.ClearValues()
	erow, ecol = 0, len(header)
	row, col = table.Size()
	if row != erow {
		t.Fatalf("expect %d, got %d", erow, row)
	}
	if col != ecol {
		t.Fatalf("expect %d, got %d", ecol, col)
	}
	
}
