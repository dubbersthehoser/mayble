package viewmodel


import (
	"errors"
	"slices"
	"cmp"
	"fmt"

	"fyne.io/fyne/v2/data/binding"

	repo "github.com/dubbersthehoser/mayble/internal/repository"
	"github.com/dubbersthehoser/mayble/internal/config"
)


// cellIndex handle for a cell.
type cellIndex uint

// noneIndex a handle to the None cell.
const noneIndex cellIndex = 0

type cellKind int 

const (
	cellNone = iota // A free, avaliable, or stub cell.
	cellTable       // Root grand parent of all cells
	cellHeader      // Cell representing a table's header
	cellView        // Cell representing a table's view data
)

// dataCell is a genric table cell type foreach kind of cell that can be linked
// to other cells intrusively, via cellPool.
type dataCell struct {
	kind cellKind

	hidden bool
	table  string
	header string

	view string
	id   int64
	v    repo.Variant

	parent cellIndex
	first  cellIndex
	next   cellIndex
	prev   cellIndex
}


// cellPool manages the collection of cells, their handels and creation.
type cellPool struct {
	cells    []dataCell
	nextFree cellIndex // single linked free list of cells
}

func newCellPool() *cellPool{
	cl := &cellPool{
		cells: make([]dataCell, 1),
	}
	cl.cells[0].view = "STUB"
	cl.cells[0].header = "STUB"
	return cl
}

func (cl *cellPool) create(k cellKind) cellIndex {
	first := cl.nextFree
	if first == 0 {
		cell := dataCell{
			kind: k,
		}
		cl.cells = append(cl.cells, cell)
		index := len(cl.cells) - 1
		return cellIndex(index)
	}
	cl.nextFree = cl.cells[first].next
	cl.cells[first].next = noneIndex
	return first
}

func (cl *cellPool) destroy(idx cellIndex) {
	cl.cells[idx] = dataCell{}
	next := cl.nextFree
	cl.cells[idx].next = next 
	cl.cells[next].prev = idx
	cl.nextFree = idx
}

func (cl *cellPool) get(i cellIndex) *dataCell {
	if i >= cellIndex(len(cl.cells)) {
		i = noneIndex
	}
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
	name          string
	headerOrder   map[string]int // maintain original header locations.
	cells         *cellPool
	root          cellIndex      // first left most header in table.
	rowCount      int            // Keep track of rows
}

func newTable(name string, headers []string) *table {
	t := &table{
		cells: newCellPool(),
		name: name,
	}

	t.cells.get(noneIndex).view = "STUB"
	t.cells.get(noneIndex).header = "STUB"
	t.cells.get(noneIndex).id = 0

	t.headerOrder = make(map[string]int)

	for i, h := range headers {
		t.headerOrder[h] = i
	}


	root  := t.cells.create(cellTable)
	t.root = root
	for _, h := range headers {
		cell := t.cells.create(cellHeader)
		t.cells.get(cell).header = h
		t.cells.get(cell).table = name
		t.cells.get(cell).parent = root
		appendCellToParent(t.cells, root, cell)
	}
	return t
}

func (t *table) addValue(id int64, header, s string) error {

	newCell := t.cells.create(cellView)
	t.cells.get(newCell).view = s
	t.cells.get(newCell).header = header
	t.cells.get(newCell).table = t.name
	t.cells.get(newCell).id = id

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

func (t *table) appendRow(id int64, values []string) error {
	headers := t.headers()
	if len(headers) != len(values) {
		return errors.New("missmatch headers to values")
	}
	for i := range headers {
		header := headers[i]
		value := values[i]
		t.addValue(id, header, value)
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
	if headerCurr == noneIndex { // nothing to remove
		return
	}
	t.rowCount -= 1
	for {
		first := t.cells.get(headerCurr).first
		curr := first
		for {
			remove := curr
			curr = t.cells.get(curr).next
			t.cells.destroy(remove)
			if curr == first {
				break
			}
		}
		t.rowCount -= 1
		headerCurr = t.cells.get(headerCurr).next
		if headerCurr == headerFirst {
			break
		}
	}
}

func (t *table) getCell(row, col int) cellIndex {
	headerIdx := t.getHeaderCell(col)
	first := t.cells.get(headerIdx).first
	curr := first
	count := 0
	for {
		if count == row {
			return curr
		}
		count += 1
		curr = t.cells.get(curr).next
		if curr == first {
			break
		}
	}
	return noneIndex
}

func (t *table) getID(idx cellIndex) (int64, error) {
	cell := t.cells.get(idx)
	if cell.kind != cellView {
		return -1, errors.New("table: invalid cell kind")
	}
	return cell.id, nil
}

func (t *table) getValue(idx cellIndex) string {
	cell := t.cells.get(idx)
	return cell.view
}

func (t *table) getHeaderCell(col int) cellIndex {
	first := t.cells.get(t.root).first
	curr := first
	count := 0
	for {
		if count == col {
			return curr
		}
		count+=1
		curr = t.cells.get(curr).next
		if first == curr {
			break
		}
	}
	return noneIndex
}

func (t *table) isHidden(idx cellIndex) bool {
	cell := t.cells.get(idx)
	switch cell.kind {
	case cellView:
		return t.cells.get(cell.parent).hidden
	case cellHeader:
		return cell.hidden
	default:
		return t.cells.get(noneIndex).hidden
	}
}

func (t *table) size() (row int, col int) {
	col = len(t.headers())
	first := t.cells.get(t.root).first
	row = cellRowLength(t.cells, first)
	return t.rowCount, col
}

// setHidden set give headers to hidden and retaining, or restoring excluded headers.
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

func (t *table) WalkAllValues(fn func(*dataCell)) {
	firstHeader := t.cells.get(t.root).first
	currHeader := firstHeader
	for {
		first := t.cells.get(currHeader).first
		curr := first
		for {
			fn(t.cells.get(curr))
			curr = t.cells.get(curr).next
			if curr == first {
				break
			}
		}
		currHeader = t.cells.get(currHeader).next
		if currHeader == firstHeader {
			break
		}
	}
}



func VariantToTableName(v repo.Variant) string {
	switch v {
	case repo.Book:
		return "All Books"
	case repo.BookLoaned:
		return "On Loan"
	case repo.BookRead:
		return "Read"
	default:
		return "N/A"
	}
}

func VariantFields(v repo.Variant) []string {
	switch v {
	case (repo.Book):
		return []string{
			"Title",
			"Author",
			"Genre",
		}
	case (repo.BookLoaned):
		return []string{
			"Title",
			"Author",
			"Genre",
			"Borrower",
			"Loaned",
		}
	case (repo.BookRead):
		return []string{
			"Title",
			"Author",
			"Genre",
			"Rating",
			"Read",
		}
	default:
		return []string{}
	}
}

func EntryValues(e *repo.BookEntry) []string {
	switch e.Variant {
	case (repo.Book):
		return []string{
			e.Title,
			e.Author,
			e.Genre,
		}
	case (repo.Book|repo.Loaned):
		return []string{
			e.Title,
			e.Author,
			e.Genre,
			e.Borrower,
			formatDate(&e.Loaned),
		}
	case (repo.Book|repo.Read):
		return []string{
			e.Title,
			e.Author,
			e.Genre,
			formatRating(e.Rating),
			formatDate(&e.Read),
		}
	default:
		return []string{}
	}
}

type TableVM struct {
	repo     repo.BookRetriever
	config   *config.Config
	v        repo.Variant

	SortBy     binding.String
	SortOrder  binding.String
	SearchText binding.String
	
	table    *table
	
	l        *listener
}


func NewTableVM(table repo.Variant, config *config.Config, store repo.BookRetriever) *TableVM {
	t := &TableVM{

		table: newTable(VariantToTableName(table), VariantFields(table)),
		repo: store,
		config: config,
		v: table,

		SortBy: binding.NewString(),
		SortOrder: binding.NewString(),
		SearchText: binding.NewString(),

		l: &listener{},
	}

	if table == repo.Book {
		t.load()
	}

	return t
}

const MinColWidth float32 = 100.0

// cleanColumnIndex transform a column index given by the ui and change it to a
// stable index of that column.
func cleanColumnIndex(t *TableVM, col int) int {
	idx := t.table.getHeaderCell(col)
	label := t.table.cells.get(idx).header
	i := slices.Index(VariantFields(t.v), label)
	return i
}

func (t *TableVM) StoreColumnWidth(col int, width float32) {
	if t.config == nil {
		return
	}
	if width < MinColWidth {
		width = MinColWidth
	}
	table := t.config.GetUITable(fmt.Sprint(t.v))
	i := cleanColumnIndex(t, col)
	_ = table.SetColWidth(i, width)
}


func (t *TableVM) GetColumnWidth(col int) float32 {
	if t.config == nil {
		return MinColWidth
	}
	table := t.config.GetUITable(fmt.Sprint(t.v))
	i := cleanColumnIndex(t, col)
	width := table.GetColWidth(i)
	if width < MinColWidth {
		width = MinColWidth
	}
	return width
}

func (t *TableVM) SetHidden(hide []string) {
	t.table.setHidden(hide)
	t.l.notify()
}


func (t *TableVM) Headers() []string {
	return t.table.headers()
}


func (t *TableVM) load() error {

	items, err := t.repo.GetAllBooks(t.v)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return nil
	}

	for _, item := range items {
		err := t.table.appendRow(
			item.ID,
			EntryValues(&item),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TableVM) Size() (int, int) {
	return t.table.size()
}

func (t *TableVM) Get(row, col int) string {
	idx := t.table.getCell(row, col)
	return t.table.getValue(idx)
}

func (t *TableVM) GetID(row, col int) (int64, error) {
	idx := t.table.getCell(row, col)
	return t.table.getID(idx)
}

func (t *TableVM) GetLabel(row, col int) (string, error) {
	idx := t.table.getCell(row, col)
	v := t.table.cells.get(idx).v
	return VariantToTableName(v), nil
}

// IsItemHidden check whether cell item is hidden.
func (t *TableVM) IsItemHidden(row, col int) bool {
	idx := t.table.getCell(row, col)
	return t.table.isHidden(idx)
}

// IsHeaderHidden check whether Header is hidden.
func (t *TableVM) IsHeaderHidden(col int) bool {
	idx := t.table.getHeaderCell(col)
	return t.table.isHidden(idx)

}

func (t *TableVM) AddListener(l binding.DataListener) {
	t.l.AddListener(l)
}

type TablesVM struct {
	config *config.Config
	table  string
	tables map[string]TableVM

	repo  repo.BookRetriever

	l *listener
}

func NewTablesVM(cfg *config.Config, q repo.BookRetriever) *TablesVM {
	t := &TablesVM{
		table:     VariantToTableName(repo.Book),
		tables:    make(map[string]TableVM),
		repo: q,
		l: &listener{},
	}
	t.loadTables()
	return t
}

func (t *TablesVM) loadTables() {
	names := []struct{
		name string
		v    repo.Variant
	}{
		{
			name: VariantToTableName(repo.Book),
			v: repo.Book,
		},
		{
			name: VariantToTableName(repo.Book|repo.Loaned),
			v: repo.BookLoaned,
		},
		{
			name: VariantToTableName(repo.Book|repo.Read),
			v: repo.BookRead,
		},
	}
	for _, name := range names {
		t.tables[name.name] = *NewTableVM(
			name.v,
			t.config,
			t.repo,
		)
	}
}

func (t *TablesVM) TableName() string {
	return t.table
}

func (t *TablesVM) GetTable(s string) *TableVM {
	table := t.tables[s]
	return &table
}

func (t *TablesVM) SetTable(s string) {
	t.table = s
	t.notify()
}

func (t *TablesVM) Variants() []repo.Variant {
	return []repo.Variant{
		repo.Book,
		repo.BookLoaned,
		repo.BookRead,
	}
}

func (t *TablesVM) TableNames() []string {
	vs := t.Variants()
	names := make([]string, len(vs))
	for i, v := range vs {
		names[i] = VariantToTableName(v)
	}
	return names
}

func (t *TablesVM) notify() {
	t.l.notify()
}

func (t *TablesVM) AddListener(l binding.DataListener) {
	t.l.AddListener(l)
}


