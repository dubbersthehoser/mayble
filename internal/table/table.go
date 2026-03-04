package table

import (
	"errors"
	"slices"
	"cmp"
)

// cellIndex handle for a cell.
type CellIndex uint

// noneIndex a handle to the None cell.
const NoneIndex CellIndex = 0

type CellKind int 

const (
	CellNone CellKind = iota // A free, avaliable, or stub cell.
	cellFree                 // Set as a free list cell. (this is to prevent a bug)
	cellTable                // Root grand parent of all cells
	CellHeader               // Cell representing a table's header
	CellView                 // Cell representing a table's view data
)

// DataCell is a genric table cell foreach kind of cell that can be linked
// to other cells intrusively, via cellPool.
type DataCell struct {
	kind CellKind

	hidden bool
	table  string

	header string

	value string
	id   int64

	parent CellIndex
	first  CellIndex
	next   CellIndex
	prev   CellIndex
}

func (dc *DataCell) Header() string {
	return dc.header
}

func (dc *DataCell) ID() int64 {
	return dc.id
}

func (dc *DataCell) Value() string {
	return dc.value
} 

func (dc *DataCell) Kind() CellKind {
	return dc.kind
}


// cellPool manages the collection of cells, their handels and creation.
type cellPool struct {
	cells    []DataCell
	nextFree CellIndex // free list of cells
}

// newCellPool create a new cell pool with stub cell.
func newCellPool() *cellPool{
	cl := &cellPool{
		cells: make([]DataCell, 1),
	}
	cl.cells[0].value = "STUB"
	cl.cells[0].header = "STUB"
	return cl
}

// create return an index of a new cell with k kind.
func (cl *cellPool) create(k CellKind) CellIndex {
	if cl.nextFree == NoneIndex {
		cell := DataCell{
			kind: k,
		}
		cl.cells = append(cl.cells, cell)
		index := len(cl.cells) - 1
		return CellIndex(index)
	}
	first := cl.nextFree
	cl.nextFree = cl.cells[first].next
	cl.cells[first].next = NoneIndex
	cl.cells[first].prev = NoneIndex
	cl.cells[first].kind = k
	return first
}

// destroy wipe cell data and add it to the free list.
func (cl *cellPool) destroy(idx CellIndex) {
	cl.cells[idx] = DataCell{
		kind: cellFree,
	}
	next := cl.nextFree
	cl.cells[idx].next = next 
	cl.nextFree = idx
}

// get returns data cell at cell index.
func (cl *cellPool) get(i CellIndex) *DataCell {
	if i >= CellIndex(len(cl.cells)) {
		i = NoneIndex
	}
	return &cl.cells[i]
}

// appendCellToParent add item cell to parent's children.
func appendCellToParent(cells *cellPool, parent, item CellIndex) {
	first := cells.get(parent).first
	if cells.get(first).kind == CellNone {
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


// A Table stores cell values of a table in a tree like structure.
//
// This structure helps keep track of hidden columns in the table by having the
// ablity to rearrange columns.
// The structure's root node is a table node its has one child node being the
// top-left header. The header links to it's the next header and its next's
// next and so forth until linking back around to the first header. Then each
// header has children of their values for that column with the same ring 
// linking as the headers.
type Table struct {
	name          string
	headerOrder   map[string]int        // Keep original header locations.
	cells         *cellPool             // Pool storing all the cells.
	root          CellIndex             // The root table table cell.
	rowCount      int                   // Keep track of rows in table.
}

func NewTable(name string, headers []string) *Table {
	t := &Table{
		cells: newCellPool(),
		name: name,
	}

	t.cells.get(NoneIndex).value = "STUB"
	t.cells.get(NoneIndex).header = "STUB"
	t.cells.get(NoneIndex).id = -127

	t.headerOrder = make(map[string]int)

	for i, h := range headers {
		t.headerOrder[h] = i
	}

	root  := t.cells.create(cellTable)
	t.root = root
	for _, h := range headers {
		cell := t.cells.create(CellHeader)
		t.cells.get(cell).header = h
		t.cells.get(cell).table = name
		t.cells.get(cell).parent = root
		appendCellToParent(t.cells, root, cell)
	}
	return t
}

// IsHidden returns whether cell is hidden.
func (t *Table) IsHidden(cell *DataCell) bool {
	switch cell.kind {
	case CellHeader:
		return cell.hidden
	case CellView:
		return t.cells.get(cell.parent).hidden
	}
	return false
}

// Name get the name of table.
func (t *Table) Name() string {
	return t.name
}


// Headers list current header order in their cell based order.
func (t *Table) Headers() []string {
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

// BaseHeader lists headers in their orignal order.
func (t *Table) BaseHeaders() []string {
	l := make([]string, len(t.headerOrder))
	for h, v  := range t.headerOrder {
		l[v] = h
	}
	return l
}


// VisableHeader list of visable headers.
func (t *Table) VisableHeaders() []string {
	headers := make([]string, 0)
	first := t.cells.get(t.root).first
	curr := first
	for {
		header := t.cells.get(curr)
		if !header.hidden {
			headers = append(headers, header.header)
		}
		curr = header.next
		if curr == first {
			break
		}
	}
	return headers
}

// HiddenHeaders list hidden headers.
func (t *Table) HiddenHeaders() []string {
	headers := make([]string, 0)
	first := t.cells.get(t.root).first
	curr := first
	for {
		header := t.cells.get(curr)
		if header.hidden {
			headers = append(headers, header.header)
		}
		curr = header.next
		if curr == first {
			break
		}
	}
	return headers
}

// addValue to column with header seting its value to s and its id.
func (t *Table) addValue(id int64, header, s string) error {

	newCell := t.cells.create(CellView)

	t.cells.get(newCell).value = s
	t.cells.get(newCell).header = header
	t.cells.get(newCell).table = t.name
	t.cells.get(newCell).id = id

	first := t.cells.get(t.root).first
	curr := first
	for { 
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


// AppendRow add a row of values to table, marked by its id.
func (t *Table) AppendRow(id int64, values []string) error {
	if len(t.Headers()) != len(values) {
		return errors.New("missmatch headers to values")
	}

	// use the original header order to add value.
	for i, header := range t.Headers() {
		value := values[i]
		err := t.addValue(id, header, value)
		if err != nil {
			return err
		}
	}
	t.rowCount += 1
	return nil
}


// clearColumnValues clear data cells form header.
func (t *Table) clearColumnValues(header CellIndex) {
	first := t.cells.get(header).first
	curr := first
	for {
		if t.cells.get(curr).kind != CellView {
			panic("invalid cell kind")
		}
		remove := curr
		curr = t.cells.get(curr).next
		t.cells.destroy(remove)
		if curr == first {
			break
		}
	}
	t.cells.get(header).first = NoneIndex
}

var clearCount int = 0

// clearValue clear all values form table, while retaining headers.
func (t *Table) ClearValues() error {
	clearCount += 1
	first := t.cells.get(t.root).first
	curr := first
	if first == NoneIndex { // nothing to remove
		return nil
	}

	for {
		if t.cells.get(curr).kind != CellHeader {
			panic("invalid cell kind")
		}
		t.clearColumnValues(curr)
		curr = t.cells.get(curr).next
		if curr == first {
			break
		}
	}
	t.rowCount = 0
	return nil
}

// GetCell with row to column number.
func (t *Table) GetCell(row, col int) *DataCell {
	hCell := t.GetHeaderCell(col)
	first := hCell.first
	curr := first
	count := 0
	for {
		if count == row {
			return t.cells.get(curr)
		}
		count += 1
		curr = t.cells.get(curr).next
		if curr == first {
			break
		}
	}
	return t.cells.get(NoneIndex)
}

// GetHeaderCell find the cell for column col.
func (t *Table) GetHeaderCell(col int) *DataCell {
	first := t.cells.get(t.root).first
	curr := first
	count := 0
	for {
		if count == col {
			return t.cells.get(curr)
		}
		count+=1
		curr = t.cells.get(curr).next
		if first == curr {
			break
		}
	}
	return t.cells.get(NoneIndex)
}

// Size get the size of the table including the header.
func (t *Table) Size() (row int, col int) {
	col = len(t.Headers())
	return t.rowCount, col
}

// SetHidden given headers to hidden and retaining, or restoring hidden headers.
func (t *Table) SetHidden(headers []string) {
	
	hiddenCells := make([]CellIndex, 0)

	showCells := make([]CellIndex, 0)

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

	// mantain the column ording.
	slices.SortFunc(showCells, func(a, b CellIndex) int {
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

// WalkVisableValues fallow all shown values from ui.
func WalkVisableValues(t *Table, fn func(vRow int, vCol int, cell *DataCell)) {
	var vRow, vCol int
	firstHeader := t.cells.get(t.root).first
	currHeader := firstHeader
	for {
		first := t.cells.get(currHeader).first
		if t.IsHidden(t.cells.get(first)) {
			return
		}
		curr := first
		for {
			fn(vRow, vCol, t.cells.get(curr))
			curr = t.cells.get(curr).next
			vRow += 1
			if curr == first {
				break
			}
		}
		vRow = 0
		vCol += 1
		currHeader = t.cells.get(currHeader).next
		if currHeader == firstHeader {
			break
		}
	}
}
