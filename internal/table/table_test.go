package table

import (
	"fmt"
	"slices"
	"testing"
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

const RowCount = 3

func headers() []string {
	return []string{
		"Title",
		"Author",
		"Genre",
		"Read",
		"Rating",
		"Loaned",
		"Borrower",
	}
}

type entryGen struct {
	count    int
	template []string
}

func newEntryGen() *entryGen {
	return &entryGen{
		template: []string{
			"example_title",
			"example_author",
			"example_genre",
			"example_read",
			"example_rating",
			"example_loaned",
			"example_borrower",
		},
	}
}

func (e *entryGen) Gen() []string {
	c := slices.Clone(e.template)
	for i := range c {
		c[i] = fmt.Sprintf("%s_%d", c[i], e.count)
	}
	e.count += 1
	return c
}

func (e *entryGen) Range(n int) [][]string {
	g := make([][]string, 0)
	for range n {
		g = append(g, e.Gen())
	}
	return g
}

func TestTable(t *testing.T) {

	var data [][]string
	gen := newEntryGen()

	for range RowCount {
		data = append(data, gen.Gen())
	}

	table := NewTable(
		"Test Table",
		headers(),
	)

	t.Run(
		"AppendRow",
		func(t *testing.T) {
			test_TableAppendRow(table, data, t)
			// Check Value
			cell := table.GetCell(2, 1)
			v := cell.Value()
			ev := data[2][1]
			if v != ev {
				t.Fatalf("expect %s, got %s", ev, v)
			}

			// Check ID
			id := cell.ID()
			eid := int64(2)
			if id != eid {
				t.Fatalf("expect %d, got %d", eid, id)
			}
		},
	)

	t.Run(
		"SetHidden",
		func(t *testing.T) {
			testSetHidden(table, t)
		},
	)

	t.Run(
		"WalkVisableValues",
		func(t *testing.T) {
			testWalkVisableValues(table, data, t)
		},
	)

	// Check Clear
	table.ClearValues()
	erow, ecol := 0, len(headers())
	row, col := table.Size()
	if row != erow {
		t.Fatalf("expect %d, got %d", erow, row)
	}
	if col != ecol {
		t.Fatalf("expect %d, got %d", ecol, col)
	}

	t.Run(
		"AppendRowAfterClear",
		func(t *testing.T) {
			test_TableAppendRow(table, data, t)
			// Check Value
			cell := table.GetCell(2, 1)
			v := cell.Value()
			ev := data[2][1]
			if v != ev {
				t.Fatalf("expect %s, got %s", ev, v)
			}

			// Check ID
			id := cell.ID()
			eid := int64(2)
			if id != eid {
				t.Fatalf("expect %d, got %d", eid, id)
			}
		},
	)

}

func test_TableAppendRow(table *Table, data [][]string, t *testing.T) {

	// Check headers
	if c := slices.Compare(headers(), table.BaseHeaders()); c != 0 {
		t.Fatalf("expect\n\t%#v\ngot\n\t%#v\n", headers(), table.BaseHeaders())
	}

	// Check append rows
	type Test struct {
		header []string
		values []string
	}

	tests := []Test{}
	for i := range data {
		test := Test{
			header: headers(),
			values: data[i],
		}
		tests = append(tests, test)
	}

	for i, c := range tests {
		err := table.AppendRow(int64(i), c.values)
		unexpectedError(t, err)
	}

	// Check size
	erow, ecol := RowCount, len(headers())
	row, col := table.Size()
	if row != erow {
		t.Fatalf("expect %d, got %d", erow, row)
	}
	if col != ecol {
		t.Fatalf("expect %d, got %d", ecol, col)
	}
}

func testSetHidden(table *Table, t *testing.T) {
	// Check hide header

	hideTests := []struct {
		input  []string
		expect []string
	}{
		{
			input: []string{},
			expect: []string{
				"Title",
				"Author",
				"Genre",
				"Read",
				"Rating",
				"Loaned",
				"Borrower",
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
				"Loaned",
				"Borrower",
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
				"Loaned",
				"Borrower",
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
	table.SetHidden([]string{})
}
func testWalkVisableValues(table *Table, data [][]string, t *testing.T) {

	var expect [][]string
	for _, s := range data {
		item := slices.Clone(s)
		expect = append(expect, item)
	}

	actual := make([][]string, 0)

	WalkVisableValues(table, func(row, col int, cell *DataCell) {
		if len(actual) == row {
			actual = append(actual, make([]string, 0))
		}
		actual[row] = append(actual[row], cell.Value())
	})

	if len(actual) != len(expect) {
		t.Fatalf("exepct length %d, got %d", len(expect), len(actual))
	}

	for row := range actual {
		if len(actual[row]) != len(expect[row]) {
			t.Fatalf("row=%d: exepct length %d, got %d", row, len(expect[row]), len(actual[row]))
		}
		for col := range actual {
			if actual[row][col] != expect[row][col] {
				t.Fatalf("col=%d, row=%d: expect %s, got %s", col, row, actual[row][col], expect[row][col])
			}
		}
	}

	table.SetHidden([]string{"Title"})
	defer table.SetHidden([]string{})

	for row := range expect {
		expect[row] = expect[row][1:]
	}

	actual = make([][]string, 0)

	WalkVisableValues(table, func(row, col int, cell *DataCell) {
		if len(actual) == row {
			actual = append(actual, make([]string, 0))
		}
		actual[row] = append(actual[row], cell.Value())
	})

	if len(actual) != len(expect) {
		t.Fatalf("exepct length %d, got %d", len(expect), len(actual))
	}

	for row := range actual {
		if len(actual[row]) != len(expect[row]) {
			t.Fatalf("row=%d: exepct length %d, got %d", row, len(expect[row]), len(actual[row]))
		}
		for col := range actual {
			if actual[row][col] != expect[row][col] {
				t.Fatalf("col=%d, row=%d: expect %s, got %s", col, row, actual[row][col], expect[row][col])
			}
		}
	}
}
