package viewmodel


import (
	"time"
	"errors"
	"strconv"

	"fyne.io/fyne/v2/data/binding"

	repo "github.com/dubbersthehoser/mayble/internal/repository"
)




type DataItem struct {
	id int64
	nextItem *DataItem
	data any
}

func newDataItem(id int64, data any) *DataItem {
	return &DataItem{
		id: id,
		data: data,
	}
}

func (di *DataItem) GetID() (int64) {
	return di.id
}

func (di *DataItem) AsView() (string, error) {
	switch v := di.data.(type) {
	case string:
		return v, nil
	case *time.Time:
		return formatDate(v), nil
	case int:
		return strconv.Itoa(v) , nil
	default:
		return "", errors.New("as view: invalid value")
	}
}

func (di *DataItem) next() *DataItem {
	return di.nextItem
}

func (di *DataItem) setNext(item *DataItem) {
	di.nextItem = item
}

type DataTable struct {
	headers    []string
	data       [][]*DataItem
	textLength map[string]int
}

func newDataTable(headers []string) *DataTable {
	dt := &DataTable{
		headers: headers,
		data: make([][]*DataItem, 0),
	}

	return dt
}

func (dt *DataTable) add(row []*DataItem) {
	if len(row) != len(dt.headers) {
		panic("row length header length missmatch")
	}

	dt.data = append(dt.data, row)
}

func (dt *DataTable) Size() (row, col int) {
	row = len(dt.data)
	if row == 0 {
		return 0, 0
	}
	col = len(dt.data[0])
	return row, col

}

// Headers list current header labels in table.
func (dt *DataTable) Headers() []string {
	return dt.headers
}

func (dt *DataTable) GetMaxTextLength(header string) int {
	s, ok := dt.textLength[header]
	if !ok {
		panic("header not found in text size")
	}
	return s
}

func (dt *DataTable) GetItem(row, col int) (*DataItem, error) {
	if row < 0 || col < 0 {
		return nil, errors.New("get item: out of range")
	}

	rowSize, colSize := dt.Size()
	
	if row >= rowSize || col >= colSize {
		return nil, errors.New("get item: out of range")
	}
	return dt.data[col][row], nil
}

func (dt *DataTable) GetString(row, col int) (string, error) {
	item, err := dt.GetItem(row, col)
	if err != nil {
		return "", err
	}
	switch v := item.data.(type) {
	case string:
		return v, nil
	case time.Time:
		return formatDate(&v), nil
	default:
		return "", errors.New("unknown data in table")
	}
}



type TableVM struct {
	
	table *DataTable

	query repo.BookSearcher

	SearchText binding.String
	SearchFrom binding.String
	OrderField binding.String
	OrderASC   binding.Bool

	join repo.BookJoin

	listeners []binding.DataListener

	avaliableTables []string
}

func NewTable(bs repo.BookSearcher) *TableVM {
	t := &TableVM{
		SearchText: binding.NewString(),
		OrderField: binding.NewString(),
		OrderASC: binding.NewBool(),

		join: repo.Main,

		listeners: make([]binding.DataListener, 0),

		avaliableTables: []string{
			"Loaned",
			"Read",
		},
	}
	return t
}

func (t *TableVM) search() {
	search, _ := t.SearchText.Get()
	from, _ := t.SearchFrom.Get()
	sortBy, _ := t.OrderField.Get()
	sortASC, _ := t.OrderASC.Get()

	direction := repo.DESC
	if sortASC {
		direction = repo.ASC
	}

	p := repo.BookSearchParams{
		Join: t.join,
		SortOrder: direction,
		SortBy: sortBy,
		Query: search,
		QueryBy: from,
	}

	rs, err := t.query.BookSearch(p)
	if err != nil {
		panic(err)
	}


	table := newDataTable(rs.Fields)

	for i := range rs.Items {
		switch v := rs.Items[i].(type) {
		case *repo.Book:
			row := []*DataItem{
				newDataItem(v.ID(), v.Title),
				newDataItem(v.ID(), v.Author),
				newDataItem(v.ID(), v.Genre),
			}
			table.add(row)
		case *repo.BookLoan:
			row := []*DataItem{
				newDataItem(v.ID(), v.Title),
				newDataItem(v.ID(), v.Author),
				newDataItem(v.ID(), v.Genre),

				newDataItem(v.ID(), v.Borrower),
				newDataItem(v.ID(), v.Loaned),
			}
			table.add(row)
		case *repo.BookRead:
			row := []*DataItem{
				newDataItem(v.ID(), v.Title),
				newDataItem(v.ID(), v.Author),
				newDataItem(v.ID(), v.Genre),

				newDataItem(v.ID(), v.Rating),
				newDataItem(v.ID(), v.Completed),
			}
			table.add(row)
		}
	}

	t.table = table
}

func (t *TableVM) Table() *DataTable {
	tbl := *t.table
	return &tbl
}

func (t *TableVM) OnJoin(table string) {
	// quiry for data with joined table
	switch table {
	case "":
		t.join = repo.Main
	case "Read":
		t.join = repo.Read
	case "Loaned":
		t.join = repo.Loaned
	}
	t.search()
	t.notify()
}

func (t *TableVM) TableJoins() []string {
	return t.avaliableTables
}

func (t *TableVM) OnDropColumns(headers []string) {
	// quiry for data without column and create a new table
	t.notify()
}

func (t *TableVM) notify() {
	for _, listener := range t.listeners {
		listener.DataChanged()
	}
}

func (t *TableVM) AddListener(l binding.DataListener) {
	t.listeners = append(t.listeners, l)
}

func (t *TableVM) RemoveListener(l binding.DataListener) {
	index := -1
	for i, listener := range t.listeners {
		if listener == l {
			index = i
		}
	}
	if index == -1 {
		return
	}
	t.listeners = append(t.listeners[:index], t.listeners[index-1:]...)
}
