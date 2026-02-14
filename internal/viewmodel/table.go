package viewmodel


import (
	"errors"
	"slices"
	"cmp"

	"fyne.io/fyne/v2/data/binding"

	repo "github.com/dubbersthehoser/mayble/internal/repository"
)



type cellIndex uint


type cellKind int 


const (
	cellNone = iota
	cellTable      // root grand parent of all cells
	cellHeader 
	cellView
	cellData
)


type dataCell struct {
	kind      cellKind

	hidden    bool
	table     string
	header    string

	view      string
	id        int64

	parent cellIndex
	first  cellIndex
	next   cellIndex
	prev   cellIndex
}


type cellPool struct {
	cells    []dataCell
	freeList cellIndex  // freelist of avaliable cells
}

func newCellList() *cellPool{
	cl := &cellPool{
		cells: make([]dataCell, 2),
		freeList: cellIndex(1),
	}
	return cl
}

func (cl *cellPool) avaliable(k cellKind) (cellIndex, error) {
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

func (cl *cellPool) wipe(i cellIndex) {
	cl.cells[i] = dataCell{}
	appendCellToParent(cl, cl.freeList, i)
}

func (cl *cellPool) newCell(k cellKind) (cellIndex, error) {
	cell := dataCell{
		kind: k,
	}

	cl.cells = append(cl.cells, cell)
	index := len(cl.cells) - 1
	return cellIndex(index), nil
}

func (cl *cellPool) get(i cellIndex) *dataCell {
	return &cl.cells[i]
}

func appendCellToParent(cells *cellPool, parent, item cellIndex) {
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

func cellRowLength(cells *cellPool, parent cellIndex) int {
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



type table struct {
	name     string
	headerOrder map[string]int
	cells    *cellPool
	root     cellIndex
	rowCount int
}

func newTable(cl *cellPool, name string, headers []string) *table {
	t := &table{
		cells: cl,
		name: name,
	}

	t.headerOrder = make(map[string]int)

	for i, h := range headers {
		t.headerOrder[h] = i
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

func (t *table) addValue(header, s string) error {

	newCell, _ := t.cells.avaliable(cellView)
	t.cells.get(newCell).view = s
	t.cells.get(newCell).header = header
	t.cells.get(newCell).table = t.name

	first := t.cells.get(t.root).first
	curr := first

	for { // do loop; God's loop
		if t.cells.get(curr).header == header {
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
func (t *table) appendRow(headers, values []string) error {
	if len(headers) != len(values) {
		return errors.New("missmatch headers to values")
	}
	for i := range headers {
		header := headers[i]
		value := values[i]
		t.addValue(header, value)
	}
	t.rowCount += 1
	return nil
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

func (t *table) getCell(row, col int) (cellIndex, error) {
	headerIdx, err := t.getHeaderCell(col)
	if err != nil {
		return 0, err
	}
	first := t.cells.get(headerIdx).first
	curr := first
	count := 0
	for {
		if count == row {
			return curr, nil
		}
		count += 1
		curr = t.cells.get(curr).next
		if curr == first {
			break
		}
	}
	return 0, errors.New("table: row not found")
}

func (t *table) getID(idx cellIndex) (int64, error) {
	cell := t.cells.get(idx)
	if cell.kind != cellView {
		return -1, errors.New("table: invalid cell kind")
	}
	return cell.id, nil
}

func (t *table) getValue(idx cellIndex) (string, error) {
	cell := t.cells.get(idx)
	if cell.kind != cellView {
		return "", errors.New("table: invalid cell kind")
	}
	return cell.view, nil
}

func (t *table) getHeaderCell(col int) (cellIndex, error) {
	first := t.cells.get(t.root).first
	curr := first
	count := 0
	for {
		if count == col {
			return curr, nil
		}
		count+=1
		curr = t.cells.get(curr).next
		if first == curr {
			break
		}
	}
	return 0, errors.New("table: header not found")
}


func (t *table) isHidden(idx cellIndex) (bool, error) {
	cell := t.cells.get(idx)
	switch cell.kind {
	case cellView:
		return t.cells.get(cell.parent).hidden, nil
	case cellHeader:
		return cell.hidden, nil
	default:
		return false, errors.New("invalid cell kind")
	}
}

func (t *table) size() (row int, col int) {
	col = len(t.headers())
	first := t.cells.get(t.root).first
	row = cellRowLength(t.cells, first)
	return t.rowCount, col
}

func (t *table) setHidden(headers []string) {

	hiddenCells := make([]cellIndex, 0)

	showCells := make([]cellIndex, 0)

	firstHeader := t.cells.get(t.root).first
	currHeader := firstHeader
	for {
		cell := t.cells.get(currHeader)
		idx := slices.Index(headers, cell.header)
		cell.hidden = idx != -1 // true when in hidden headers
		if cell.hidden {
			hiddenCells = append(hiddenCells, currHeader)
		} else {
			showCells = append(showCells, currHeader)
		}
		currHeader = cell.next
		if currHeader == firstHeader {
			break
		}
	}

	slices.SortFunc(showCells, func(a, b cellIndex) int {
		APlace := t.headerOrder[t.cells.get(a).header]
		BPlace := t.headerOrder[t.cells.get(b).header]
		return cmp.Compare(APlace, BPlace)
	})


	// rearrange columns to keep shown columns to the left and hidden to the right.
	parent := t.root
	t.cells.get(parent).first = 0
	for _, idx := range(showCells) {
		appendCellToParent(t.cells, parent, idx)
	}
	for _, idx := range(hiddenCells) {
		appendCellToParent(t.cells, parent, idx)
	}
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
		return []string{"N/A"}
	}
}

func TableLabel(name string) string {
	r := repo.BookJoin(name)
	switch r {
	case repo.Main:
		return "All Books"
	case repo.Loaned:
		return "On Loan"
	case repo.Read:
		return "Read"
	default:
		return "n/a"
	}
}
func TableName(label string) string {
	switch label {
	case "Books":
		return string(repo.Main)
	case "On Loan":
		return string(repo.Loaned)
	case "Read":
		return string(repo.Read)
	default:
		return "n/a"
	}
	
}


func resultAsString(r repo.Resultable, field string) string {
	switch data := r.(type) {
	case *repo.Book:
		switch field {
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
		switch field {
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
		switch field {
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
//type ItemSelected struct {
//	listeners []binding.DataListener
//	id        int64 
//}
//
//// Get return item when item isn't nil otherwise returns an error.
//func (is *ItemSelected) Get() (int64, error) {
//	if is.id == -1 {
//		return is.id, errors.New("item not selected")
//	}
//	return is.id, nil
//}
//
//func (is *ItemSelected) Set(id int64) (error) {
//	is.id = id
//	is.notify()
//	return nil
//}
//
//func (is *ItemSelected) notify() {
//	for _, l := range is.listeners {
//		l.DataChanged()
//	}
//}
//
//func (is *ItemSelected) AddListener(l binding.DataListener) {
//	is.listeners = append(is.listeners, l)
//	
//}
//
//func (is *ItemSelected) RemoveListener(l binding.DataListener) {
//	index := slices.Index(is.listeners, l)
//	if index == -1 {
//		return
//	}
//	is.listeners = append(is.listeners[:index], is.listeners[index-1:]...)
//}
//
//func NewItemSelected() *ItemSelected {
//	return &ItemSelected{
//		listeners: make([]binding.DataListener, 0),
//	}
//}



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
	repo     repo.BookSearcher
	Query    *QueryVM
	table    *table
	
	l        *listener
}

func NewTableVM(table repo.BookJoin, headers []string, query *QueryVM, store repo.BookSearcher) *TableVM {
	t := &TableVM{
		table: newTable(newCellList(), string(table), headers),
		Query: query,
		repo: store,
		l: &listener{},
	}

	if table == repo.Main {
		t.load()
	}

	return t
}

func (t *TableVM) SetHidden(hide []string) {
	t.table.setHidden(hide)
	t.l.notify()
}


func (t *TableVM) Headers() []string {
	return t.table.headers()
}

func (t *TableVM) load() error {
	
	param := repo.BookSearchParams{
		Join: *t.Query.table,
	}

	rSet, err := t.repo.BookSearch(param)
	if err != nil {
		return err
	}

	if len(rSet.Items) == 0 {
		return nil
	}

	for _, result := range rSet.Items {
		fields := getResultFields(result)
		values := make([]string, len(fields))
		for i, field := range fields {
			value := resultAsString(result, field)
			values[i] = value
		}
		err := t.table.appendRow(fields, values)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TableVM) Size() (int, int) {
	return t.table.size()
}

func (t *TableVM) Get(row, col int) (string, error) {
	idx, err := t.table.getCell(row, col)
	if err != nil {
		return "N/A", err
	}
	return t.table.getValue(idx)
}

func (t *TableVM) GetID(row, col int) (int64, error) {
	idx, err := t.table.getCell(row, col)
	if err != nil {
		return -1, err
	}
	return t.table.getID(idx)
}


func (t *TableVM) IsItemHidden(row, col int) (bool, error) {
	idx, err := t.table.getCell(row, col)
	if err != nil {
		return false, err
	}
	return t.table.isHidden(idx)
}
func (t *TableVM) IsHeaderHidden(col int) (bool, error) {
	idx, err := t.table.getHeaderCell(col)
	if err != nil {
		return false, err
	}
	return t.table.isHidden(idx)

}

func (t *TableVM) AddListener(l binding.DataListener) {
	t.l.AddListener(l)
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
	_ = t.LoadTables()
	return t
}

func (t *TablesVM) LoadTables() error {
	for _, name := range t.TableNames() {
		table := repo.BookJoin(name)
		t.tables[table] = *NewTableVM(
			table,
			getTableHeaders(table),
			t.Query,
			t.repo,
		)
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
