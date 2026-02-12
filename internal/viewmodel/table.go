package viewmodel


import (
	"time"
	"errors"
	"slices"
	"strconv"

	"fyne.io/fyne/v2/data/binding"

	repo "github.com/dubbersthehoser/mayble/internal/repository"
)


type cellIndex uint

const nilIndex cellIndex = 0

type cellKind int 
const (
	cellNone = iota
	cellTable      // root grand parent of all cells
	cellHeader 
	cellData
)

type dataCell struct {
	kind      cellKind

	label     string    // header label
	name      string    // table name
	data      any       // data

	first cellIndex
	next  cellIndex
	prev  cellIndex
}

type cellList struct {
	cells []dataCell
}
func newCellList() *cellList{
	return &cellList{
		cells: make([]dataCell, 1),
	}
}
func (cl *cellList) avaliable() (cellIndex, error) {
	for i, cell := range cl.cells[1:] {
		if cell.kind == cellNone {
			return cellIndex(i), nil
		}
	}
	return cl.newCell(cellNone)
}

func (cl *cellList) wipe(i cellIndex) {
	cl.cells[i] = dataCell{}
}

func (cl *cellList) newCell(k cellKind) (cellIndex, error) {
	cell := dataCell{
		kind: k,
	}

	cl.cells = append(cl.cells, cell)
	index := len(cl.cells) - 1
	return cellIndex(index), nil
}
func (cl *cellList) get(i cellIndex) *dataCell {
	return &cl.cells[i]
}

type table struct {
	cells *cellList
	root  cellIndex
}

func NewTable(cl *cellList, headers []string) *table {
	t := &table{
		cells: cl,
		
	}

	curr, _ := t.cells.newCell(cellTable)

	t.root = curr

	var first cellIndex

	for i, h := range headers {
		cell, _ := t.cells.newCell(cellHeader)
		if i == 0 {
			t.cells.get(t.root).first = cell
			t.cells.get(t.root).prev = cell
			t.cells.get(t.root).next = cell
			curr = cell
			first = cell
		} else {
			t.cells.get(curr).next = cell
			t.cells.get(first).prev = cell
			curr = cell
		}
		t.cells.get(cell).label = h
	}
	return t
}

func (t *table) dataClear() {
	
}


















type DataItem struct {
	field string
	id   int64
	data any
}

func newDataItem(id int64, data any, field string) *DataItem {
	return &DataItem{
		id: id,
		data: data,
		field: field,
	}
}

func (di *DataItem) Field() string {
	return di.field
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


type DataTable struct {
	name       string
	headers    []string
	exclude    []string
	data       [][]DataItem
	textLength map[string]int
}

func newDataTable(name string, headers []string) *DataTable {
	dt := &DataTable{
		name: name,
		headers: headers,
		data: make([][]DataItem, 0),
		textLength: make(map[string]int),
	}
	for _, h := range headers {
		dt.textLength[h] = len(h)
	}
	return dt
}

func (dt *DataTable) deleteData() {
	dt.data = dt.data[:0]
}

func (dt *DataTable) add(row []DataItem) {
	if len(row) != len(dt.Headers()) {
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
	if len(dt.exclude) == 0 {
		return dt.headers
	}
	set := []string{}
	for _, h := range dt.headers {
		if slices.Contains(dt.exclude, h) {
			continue
		}
		set = append(set, h)
	}
	return set
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

	if len(dt.exclude) == 0 {
		return &dt.data[row][col], nil
	}

	headers := dt.Headers()
	col = slices.Index(dt.headers, headers[col])
	return &dt.data[row][col], nil
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


// itemSelected manage selected item form table.
type ItemSelected struct {
	listeners []binding.DataListener
	item *DataItem
}

// Get return item when item isn't nil otherwise returns an error.
func (is *ItemSelected) Get() (*DataItem, error) {
	if is.item == nil {
		return nil, errors.New("item not selected")
	}
	return is.item, nil
}

func (is *ItemSelected) Set(di *DataItem) (error) {
	is.item = di
	is.notify()
	return nil
}

func (is *ItemSelected) notify() {
	for _, l := range is.listeners {
		l.DataChanged()
	}
}

func (is *ItemSelected) AddListener(l binding.DataListener) {
	is.listeners = append(is.listeners, l)
	
}

func (is *ItemSelected) RemoveListener(l binding.DataListener) {
	index := slices.Index(is.listeners, l)
	if index == -1 {
		return
	}
	is.listeners = append(is.listeners[:index], is.listeners[index-1:]...)
}

func NewItemSelected() *ItemSelected {
	return &ItemSelected{
		listeners: make([]binding.DataListener, 0),
	}
}

type QueryVM struct {
	
	table      repo.BookJoin

	SearchText binding.String
	SearchFrom binding.String

	OrderField binding.String
	OrderASC   binding.Bool
}

type TablesVM struct {
	
	table    repo.BookJoin
	tables   map[repo.BookJoin]DataTable

	selected *ItemSelected

	repo  repo.BookSearcher
	Query *QueryVM

	listeners []binding.DataListener

}

func NewTablesVM(bs repo.BookSearcher) *TablesVM {
	t := &TablesVM{
		table:     repo.Main,
		tables:    make(map[repo.BookJoin]DataTable),
		selected: &ItemSelected{},
		repo: bs,
	}
	return t
}

func (t *TablesVM) LoadTables() error {
	for _, name := range t.TableNames() {
		t.tables[repo.BookJoin(name)] = DataTable{name: string(name)}
	}

	return nil
}

func (t *TablesVM) TableName() string {
	return string(t.table)
}

func (t *TablesVM) Table() *DataTable {
	table := t.tables[t.table]
	return &table
}

func (t *TablesVM) TableNames() []string {
	return []string{
		string(repo.Main),
		string(repo.Loaned),
		string(repo.Read),
	}
}

func (t *TablesVM) search() {
	search, _ := t.Query.SearchText.Get()
	from, _ := t.Query.SearchFrom.Get()
	sortBy, _ := t.Query.OrderField.Get()
	sortASC, _ := t.Query.OrderASC.Get()

	direction := repo.DESC
	if sortASC {
		direction = repo.ASC
	}

	p := repo.BookSearchParams{
		Join: t.table,
		SortOrder: direction,
		SortBy: sortBy,
		Query: search,
		QueryBy: from,
	}

	rs, err := t.repo.BookSearch(p)
	if err != nil {
		panic(err)
	}

	table := t.tables[t.table]
	table.deleteData()
	for i := range rs.Items {
		switch v := rs.Items[i].(type) {
		case *repo.Book:
			row := []DataItem{
				*newDataItem(v.ID(), v.Title, rs.Fields[0]),
				*newDataItem(v.ID(), v.Author, rs.Fields[1]),
				*newDataItem(v.ID(), v.Genre, rs.Fields[2]),
			}
			table.add(row)
		case *repo.BookLoan:
			row := []DataItem{
				*newDataItem(v.ID(), v.Title, rs.Fields[0]),
				*newDataItem(v.ID(), v.Author, rs.Fields[1]),
				*newDataItem(v.ID(), v.Genre, rs.Fields[2]),

				*newDataItem(v.ID(), v.Borrower, rs.Fields[3]),
				*newDataItem(v.ID(), v.Loaned, rs.Fields[4]),
			}
			table.add(row)
		case *repo.BookRead:
			row := []DataItem{
				*newDataItem(v.ID(), v.Title, rs.Fields[0]),
				*newDataItem(v.ID(), v.Author, rs.Fields[1]),
				*newDataItem(v.ID(), v.Genre, rs.Fields[2]),

				*newDataItem(v.ID(), v.Rating, rs.Fields[3]),
				*newDataItem(v.ID(), v.Completed, rs.Fields[4]),
			}
			table.add(row)
		}
	}
	t.tables[t.table] = table
}

func (t *TablesVM) notify() {
	for _, listener := range t.listeners {
		listener.DataChanged()
	}
}

func (t *TablesVM) AddListener(l binding.DataListener) {
	if t.listeners == nil {
		t.listeners = make([]binding.DataListener, 0)
	}
	t.listeners = append(t.listeners, l)
}

func (t *TablesVM) RemoveListener(l binding.DataListener) {
	if t.listeners == nil {
		return
	}
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
