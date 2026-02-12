package viewmodel


import (
	"time"
	"slices"
	"testing"
	"strconv"
)

func unexpectedError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TextCellList(t *testing.T) {
	cells := newCellList()

	idx, err := cells.newCell(cellTable)
	unexpectedError(t, err)

	if cells.get(idx).kind != cellTable {
		t.Fatalf("expect kind %d, got %d", cellTable, cells.get(idx).kind)
	}

	cells.get(idx).table = "TableA"

	if cells.get(idx).table != "TableA" {
		t.Fatalf("expect %s, got %s", "TableA", cells.get(idx).table)
	}
}


type testDataItem struct {
	id int64
	header string
}
func (td *testDataItem) AsString() string {
	return strconv.Itoa(int(td.id))
}
func (td *testDataItem) Header() string {
	return td.header
}
func (td *testDataItem) ID() int64 {
	return td.id
}


func TestTable(t *testing.T) {
	
	name := "TableA"
	header := []string{
		"Title",
		"Author",
		"Genre",
		"Borrower",
		"Rating",
	}
	table := newTable(
		newCellList(),
		name,
		header,
	)

	if c := slices.Compare(header, table.headers()); c != 0 {
		t.Fatalf("expect\n\t%#v\ngot\n\t%#v\n", header, table.headers())
	}

	tests := []struct{
		id     int64
		header string
		row    int
	}{
		{ id: 1242, header: "Title", row: 0},
		{ id: 4323, header: "Title", row: 1},
		{ id: 431, header: "Title", row: 2},
		{ id: 432, header: "Title", row: 3},

		{ id: 33, header: "Author", row: 0},
		{ id: 8423, header: "Author", row: 1},
		{ id: 840, header: "Author", row: 2},
		{ id: 84342, header: "Author", row: 3},

		{ id: 1021, header: "Genre", row: 0},
		{ id: 1324, header: "Genre", row: 1},
		{ id: 1324, header: "Genre", row: 2},
		{ id: 1324, header: "Genre", row: 3},

		{ id: 1, header: "Rating", row: 0},
		{ id: 13, header: "Rating", row: 1},
		{ id: 24, header: "Rating", row: 2},
		{ id: 32, header: "Rating", row: 3},
	}

	for i, c := range tests {
		value := &testDataItem{
			id: c.id,
			header: c.header,
		}
		table.addValue(value)
		v, err := table.getValue(c.row, c.header)
		unexpectedError(t, err)
		if v.ID() != value.ID() {
			t.Fatalf("%d: expect %d, got %d", i, value.ID(), v.ID())
		}
	}
	row, col := table.size()
	titleLen := 4
	if row != titleLen {
		t.Fatalf("expect %d, got %d", titleLen, row)
	}
	if col != len(header) {
		t.Fatalf("expect %d, got %d", titleLen, col)
	}
	table.clearValues()
	row, col = table.size()
	if row != 0 {
		t.Fatalf("expect %d, got %d", 0, row)
	}
	if col != len(header) {
		t.Fatalf("expect %d, got %d", 0, col)
	}
}


func TestDataItem(t *testing.T) {
	item := newDataItem(123, "hello", "")

	v, err := item.AsView()
	unexpectedError(t, err)

	if v != "hello" {
		t.Fatalf("expect '%s', got '%s'", "hello", v)
	}

	if item.GetID() != 123 {
		t.Fatalf("expect %d, got %d", 123, item.GetID())
	}

	item2 := newDataItem(32, 21, "")
	v, err = item2.AsView()
	unexpectedError(t, err)
	if v != "21" { // Note: This will be changed to use stars.
		t.Fatalf("expect \"21\", got \"%s\"", v)
	}

	date := time.Date(2000, 02, 01, 0,0,0,0, time.UTC)
	item3 := newDataItem(43, &date, "")
	v, err = item3.AsView()
	unexpectedError(t, err)
	if v != "01/02/2000" {
		t.Fatalf("expect '01/02/2000', got %s", v)
	}

}

func TestDataTable(t *testing.T) {
	
	table := newDataTable(
		"Main",
		[]string{
			"Title",
			"Author",
			"Genre",
		},
	)

	rows := [][]string{
		{
		"book 0",
		"author 0",
		"genre 0",
		},
		{
		"book 1",
		"author 1",
		"genre 1",
		},
		{
		"book 2",
		"author 2",
		"genre 2",
		},
	}

	for i, row := range rows {
		items := []DataItem{
			*newDataItem(int64(i), row[0], "Title"),
			*newDataItem(int64(i), row[1], "Author"),
			*newDataItem(int64(i), row[2], "Genre"),
		}
		table.add(items)
	}

	table.exclude = []string{}

	v, err := table.GetString(2, 0)
	unexpectedError(t, err)

	expect := rows[2][0]
	if v != expect {
		t.Fatalf("expect '%s', got '%s'", expect, v)
	}

	table.exclude = []string{
		"Author",
	}

	ha := table.Headers()
	he := []string{
		"Title",
		"Genre",
	}

	if r := slices.Compare(ha, he); r != 0 {
		t.Fatalf("expect headers\n\t%#v\ngot\n\t%#v", he, ha)
	}

	v, err = table.GetString(1, 1)
	unexpectedError(t, err)

	expect = rows[1][2]
	if v != expect {
		t.Fatalf("expect '%s', got '%s'", expect, v)
	}
}
