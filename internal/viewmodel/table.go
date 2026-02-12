package viewmodel


import (
	"errors"
	"slices"

	"fyne.io/fyne/v2/data/binding"

	repo "github.com/dubbersthehoser/mayble/internal/repository"
)



type ValueItemView interface {
	Header()   string
	AsString() string
	ID()       int64
}


type cellIndex uint


type cellKind int 


const (
	cellNone = iota
	cellTable      // root grand parent of all cells
	cellHeader 
	cellData
)


type dataCell struct {
	kind      cellKind

	table     string
	header    string
	data      ValueItemView // data

	parent cellIndex
	first cellIndex
	next  cellIndex
	prev  cellIndex
}


type cellList struct {
	cells    []dataCell
	freeList cellIndex  // freelist of avaliable cells
}

func newCellList() *cellList{
	cl := &cellList{
		cells: make([]dataCell, 2),
		freeList: cellIndex(1),
	}
	return cl
}

func (cl *cellList) avaliable(k cellKind) (cellIndex, error) {
	// use free list to get next wiped cell.
	first := cl.get(cl.freeList).first
	if first == 0 {
		return cl.newCell(k)
	}
	cell := first
	cl.get(cl.freeList).first = cl.get(first).next
	cl.get(cell).next = 0
	cl.get(cell).kind = k
	return cell, nil
}

func (cl *cellList) wipe(i cellIndex) {
	cl.cells[i] = dataCell{}
	appendCellToParent(cl, cl.freeList, i)
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
	name  string
	cells *cellList
	root  cellIndex
}


func appendCellToParent(cells *cellList, parent, item cellIndex) {
	first := cells.get(parent).first
	if cells.get(first).kind == cellNone {
		cells.get(parent).first = item
		cells.get(item).next = item
		cells.get(item).prev = item
	} else {
		last := cells.get(first).prev
		cells.get(first).prev = item
		cells.get(item).prev = last
		cells.get(item).next = first
		cells.get(last).next = item
	}
	cells.get(item).parent = parent
}
func cellRowLength(cells *cellList, parent cellIndex) int {
	first := cells.get(parent).first
	if cells.get(first).kind == cellNone {
		return 0
	}
	curr := first
	count := 0
	for {
		curr = cells.get(curr).next
		count += 1
		if curr == first {
			break
		}
	}
	return count
}


func newTable(cl *cellList, name string, headers []string) *table {
	t := &table{
		cells: cl,
		name: name,
	}

	root, _ := t.cells.newCell(cellTable)
	t.root = root
	for _, h := range headers {
		cell, _ := t.cells.newCell(cellHeader)
		t.cells.get(cell).header = h
		t.cells.get(cell).table = name
		t.cells.get(cell).parent = root
		appendCellToParent(t.cells, root, cell)
	}
	return t
}

func (t *table) addValue(data ValueItemView) error {

	newCell, _ := t.cells.avaliable(cellData)
	t.cells.get(newCell).data = data
	t.cells.get(newCell).header = data.Header()
	t.cells.get(newCell).table = t.name

	first := t.cells.get(t.root).first
	curr := first

	for { // do loop; God's loop

		if t.cells.get(curr).header == data.Header() {
			appendCellToParent(t.cells, curr, newCell)
			return nil
		}
		
		curr = t.cells.get(curr).next

		if curr == first {
			break
		}
	}
	return errors.New("table: header not found")
}

func (t *table) headers() []string {
	
	headers := make([]string, 0)

	first := t.cells.get(t.root).first
	curr := first

	for {
		headers = append(headers, t.cells.get(curr).header)
		curr = t.cells.get(curr).next
		if curr == first {
			break
		}
	}
	return headers
}

func (t *table) clearValues() {
	headerFirst := t.cells.get(t.root).first
	headerCurr := headerFirst
	for {
		first := t.cells.get(headerCurr).first
		curr := first
		for {
			remove := curr
			curr = t.cells.get(remove).next
			t.cells.wipe(remove)
			if curr == first {
				break
			}
		}
		if headerCurr == headerFirst {
			break
		}
	}
}

func (t *table) getValue(row int, header string) (ValueItemView, error) {
	rootHeader := cellIndex(0)
	firstHeader := t.cells.get(t.root).first
	currHeader := firstHeader
	for {
		if t.cells.get(currHeader).header == header {
			rootHeader = currHeader
		}
		currHeader = t.cells.get(currHeader).next
		if currHeader == firstHeader {
			break
		}
	}

	first := t.cells.get(rootHeader).first
	curr := first
	count := 0
	for {
		if count == row {
			return t.cells.get(curr).data, nil
		}
		count += 1
		curr = t.cells.get(curr).next
		if curr == first {
			break
		}
	}
	return nil, errors.New("table: value not found")
}

func (t *table) size() (row int, col int) {
	col = len(t.headers())
	first := t.cells.get(t.root).first
	row = cellRowLength(t.cells, first)
	return row, col
}


func searchTable(t *table, s string) (float64, int, int) {
	firstHeader := t.cells.get(t.root).first
	currHeader := firstHeader

	row := 0
	col := 0

	for {
		currHeader = t.cells.get(currHeader).next
		
		first := t.cells.get(currHeader).first
		curr := first
		for {
			data := t.cells.get(curr).data
			if data.AsString() == s {
				return 1.0, row, col
			}
			row += 1
			curr = t.cells.get(curr).next
			if curr == first {
				break
			}
		}
		col += 1
		
		if currHeader == firstHeader {
			break
		}
	}
	return 0.0, row, col
} 



func getResultFields(r repo.Resultable) []string {
	switch r.(type) {
	case *repo.Book:
		return []string{
			"Title",
			"Author",
			"Genre",
		}
	case *repo.BookLoan:
		return []string{
			"Title",
			"Author",
			"Genre",
			"Borrower",
			"Loaned",
		}
	case *repo.BookRead:
		return []string{
			"Title",
			"Author",
			"Genre",
			"Rating",
			"Completed",
		}
	default:
		return []string{}
	}
}


func getTableHeaders(name repo.BookJoin) []string {
	switch name {
	case repo.Main:
		return getResultFields(&repo.Book{})
	case repo.Loaned:
		return getResultFields(&repo.BookLoan{})
	case repo.Read:
		return getResultFields(&repo.BookRead{})
	default:
		return []string{}
	}
}



type loadedItem struct {
	field  string
	column int
	table *TableVM
}

func newLoadedItem(t *TableVM, column int, field string) *loadedItem {
	return &loadedItem{
		field: field,
		table: t,
		column: column,
	}
} 

func (bi *loadedItem) ID() int64 {
	data := bi.table.loaded[bi.column]
	return data.ID()
}

func (bi *loadedItem) AsString() string {
	r := bi.table.loaded[bi.column]
	switch data := r.(type) {
	case *repo.Book:
		switch bi.field {
		case "Title":
			return data.Title
		case "Author":
			return data.Author
		case "Genre":
			return data.Genre
		default:
			panic("bookItem: field not found")
		}
	case *repo.BookLoan:
		switch bi.field {
		case "Title":
			return data.Title
		case "Author":
			return data.Author
		case "Genre":
			return data.Genre
		case "Borrower":
			return data.Borrower
		case "Loaned":
			return formatDate(data.Loaned)
		default:
			panic("bookItem: field not found")
		}
	case *repo.BookRead:
		switch bi.field {
		case "Title":
			return data.Title
		case "Author":
			return data.Author
		case "Genre":
			return data.Genre
		case "Rating":
			return formatRating(data.Rating)
		case "Loaned":
			return formatDate(data.Completed)
		default:
			panic("bookItem: field not found")
		}
	default:
		return "PlaceHolder"
	}
}



// itemSelected manage selected item form table.
type ItemSelected struct {
	listeners []binding.DataListener
	item ValueItemView
}

// Get return item when item isn't nil otherwise returns an error.
func (is *ItemSelected) Get() (ValueItemView, error) {
	if is.item == nil {
		return nil, errors.New("item not selected")
	}
	return is.item, nil
}

func (is *ItemSelected) Set(vi ValueItemView) (error) {
	is.item = vi
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
	table      *repo.BookJoin
	SearchText binding.String
	SearchFrom binding.String
	OrderField binding.String
	OrderASC   binding.Bool
}

func newQueryVM(table *repo.BookJoin) *QueryVM {
	return &QueryVM{
		table: table,
		OrderField: binding.NewString(),
		OrderASC:   binding.NewBool(),
		SearchText: binding.NewString(),
	}
} 


type TableVM struct {
	selected *ItemSelected
	repo     repo.BookSearcher
	Query    *QueryVM
	loaded   []repo.Resultable
	view     *table
	
	l        *listener
}

func NewTableVM(table repo.BookJoin, headers []string, query *QueryVM) *TableVM {
	return &TableVM{
		view: newTable(newCellList(), string(table), headers),
		loaded: make([]repo.Resultable, 0),
		selected: &ItemSelected{},
		Query: query,
		l: &listener{},
	}
}

func (t *TableVM) AddListener(l binding.DataListener) {
	t.l.AddListener(l)
}

func (t *TableVM) Headers() []string {
	return t.view.headers()
}

func (t *TableVM) load() error {
	
	param := repo.BookSearchParams{
		Join: *t.Query.table,
	}

	rSet, err := t.repo.BookSearch(param)
	if err != nil {
		return err
	}

	t.loaded = rSet.Items
	return nil
}

func (t *TableVM) Size() (int, int) {
	return t.view.size()
}

func (t *TableVM) Get(row, col int) (ValueItemView, error) {
	name := t.view.headers()[col]
	v, err := t.view.getValue(row, name)
	if err != nil {
		return nil, err
	}
	return v, nil
}



type TablesVM struct {
	
	table    repo.BookJoin
	tables   map[repo.BookJoin]TableVM

	repo  repo.BookSearcher
	Query *QueryVM

	l *listener
}

func NewTablesVM(bs repo.BookSearcher) *TablesVM {
	t := &TablesVM{
		table:     repo.Main,
		tables:    make(map[repo.BookJoin]TableVM),
		repo: bs,
		l: &listener{},
	}
	t.Query = newQueryVM(&t.table)
	t.LoadTables()
	return t
}

func (t *TablesVM) LoadTables() error {
	for _, name := range t.TableNames() {
		table := repo.BookJoin(name)
		t.tables[table] = *NewTableVM(table, getTableHeaders(table), t.Query)
	}
	return nil
}

func (t *TablesVM) TableName() string {
	return string(t.table)
}

func (t *TablesVM) GetTable(s string) *TableVM {
	table := t.tables[repo.BookJoin(s)]
	return &table
}

func (t *TablesVM) Table() *TableVM {
	table := t.tables[t.table]
	return &table
}

func (t *TablesVM) SetTable(s string) {
	t.table = repo.BookJoin(s)
	t.notify()
}

func (t *TablesVM) TableNames() []string {
	return []string{
		string(repo.Main),
		string(repo.Loaned),
		string(repo.Read),
	}
}

func (t *TablesVM) notify() {
	t.l.notify()
}

func (t *TablesVM) AddListener(l binding.DataListener) {
	t.l.AddListener(l)
}



type listener struct {
	listeners []binding.DataListener
}

func (t *listener) notify() {
	for _, listener := range t.listeners {
		listener.DataChanged()
	}
}

func (t *listener) AddListener(l binding.DataListener) {
	if t.listeners == nil {
		t.listeners = make([]binding.DataListener, 0)
	}
	t.listeners = append(t.listeners, l)
}

func (t *listener) RemoveListener(l binding.DataListener) {
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
