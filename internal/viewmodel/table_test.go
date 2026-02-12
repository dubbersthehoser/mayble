package viewmodel


import (
	"time"
	"slices"
	"testing"
)

func unexpectedError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
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
	
	table := newDataTable([]string{
		"Title",
		"Author",
		"Genre",
	})

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
		items := []*DataItem{
			newDataItem(int64(i), row[0]),
			newDataItem(int64(i), row[1]),
			newDataItem(int64(i), row[2]),
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
