package table

import (
	"errors"
	"fmt"
	"cmp"
	"slices"
)

type Table struct {
	headerOrder map[string]int // Keep original header locations.
	rowCount    int            // Keep track of rows in table.
	name        string
	first       *Header        // First header in table list.
}

func NewTable(name string, headers []string) *Table {
	t := &Table{
		name: name,
		headerOrder: make(map[string]int),
	}

	for i, h := range headers {
		t.headerOrder[h] = i
	}

	for _, header := range headers {
		h := newHeader(t, header)
		if t.first != nil {
			t.first.appendHeader(h)
		} else {
			t.first = h
		}
	}
	return t
}

type Header struct {
	name   string
	table  *Table
	hidden bool
	next   *Header
	prev   *Header

	value *Cell
}
func newHeader(t *Table, name string) *Header {
	h := &Header{
		name: name,
		table: t,
	}
	h.next = h
	h.prev = h
	return h
}
func stubHeader(t *Table) *Header {
	h := newHeader(t, "STUB")
	h.appendValue(-1, "STUB")
	return h
}

func (h *Header) IsHidden() bool {
	return h.hidden
}

func (h *Header) Name() string {
	return h.name
}

func (h *Header) appendHeader(n *Header) {
	prev := h.prev
	next := h

	n.next = next
	n.prev = prev

	next.prev = n
	prev.next = n
}

func (h *Header) appendValue(id int64, v string) {
	h.appendCell(newCell(h, id, v))
}

func (h *Header) appendCell(c *Cell) {
	if h.value == nil {
		h.value = c
		return
	}
	h.value.append(c)
}

func (h *Header) getCell(idx int) *Cell {
	curr := h.value
	i := 0
	for {
		if i == idx {
			return curr
		}
		i += 1
		curr = curr.next
		if curr == h.value {
			return stubCell(h)
		}
	}
}

type Cell struct {
	header *Header
	id     int64
	value  string
	next   *Cell
	prev   *Cell
}
func newCell(h *Header, id int64, v string) *Cell {
	c := &Cell{
		header: h,
		id: id,
		value: v,
	}
	c.next = c
	c.prev = c
	return c
}
func stubCell(h *Header) *Cell {
	return newCell(h, -1, "STUB") 
}

func (c *Cell) ID() int64 {
	return c.id
}

func (c *Cell) Value() string {
	return c.value
}

func (c *Cell) Header() string {
	return c.header.name
}

// append add cell as if c is head of list, and v is added to the end of the list.
func (c *Cell) append(v *Cell) {
	next := c
	prev := c.prev
	v.next = next
	v.prev = prev
	next.prev = v
	prev.next = v
}

// IsHidden check wheather cell is hidden.
func (c *Cell) IsHidden() bool {
	return c.header.hidden
}


// Headers returns the current header order.
func (t *Table) Headers() []string {
	headers := make([]string, 0)
	if t.first == nil {
		return headers
	}
	curr := t.first
	for {
		headers = append(headers, curr.name)
		curr = curr.next
		if t.first == curr {
			fmt.Printf("%#v\n", headers) //!
			return headers
		}
	}
}

// AppendRow add a row to table with entry id and its values.
// returns errors when the number of values don't match number of table headers.
func (t *Table) AppendRow(id int64, values []string) error {
	if len(values) != len(t.headerOrder) {
		return errors.New("table invalid value count")
	}
	curr := t.first
	i := 0
	for {
		curr.appendValue(id, values[i])
		i+=1
		curr = curr.next
		if t.first == curr {
			t.rowCount += 1
			return nil
		}
	}
}

// BaseHeaders lists original header order.
func (t *Table) BaseHeaders() []string {
	l := make([]string, len(t.headerOrder))
	for h, v := range t.headerOrder {
		l[v] = h
	}
	return l
}

// HiddenHeaders list hidden headers.
func (t *Table) HiddenHeaders() []string {
	headers := make([]string, 0)
	curr := t.first
	for {
		if curr.IsHidden() {
			headers = append(headers, curr.name)
		}
		curr = curr.next
		if curr == t.first {
			return headers
		}
	}
}

// VisableHeaders list shown headers.
func (t *Table) VisableHeaders() []string {
	headers := make([]string, 0)
	curr := t.first
	for {
		if !curr.IsHidden() {
			headers = append(headers, curr.name)
		}
		curr = curr.next
		if curr == t.first {
			return headers
		}
	}
}

// ClearValues remove values from table.
func (t *Table) ClearValues() error {
	curr := t.first
	for {
		curr.value = nil
		curr = curr.next
		if curr == t.first {
			t.rowCount = 0
			return nil
		}
	}
}

// GetCell return cell from table at row and col.
// returns stubbed values if not found
func (t *Table) GetCell(row, col int) *Cell {
	header := t.GetHeader(col)
	return header.getCell(row)
}

func (t *Table) GetHeader(col int) *Header {
	curr := t.first
	i :=  0
	for {
		if i == col {
			return curr
		}
		i += 1
		curr = curr.next
		if curr == t.first {
			return stubHeader(t)
		}
	}
}

func (t *Table) Name() string {
	return t.name
}

func (t *Table) SetHidden(headers []string) {

	hidden := make([]*Header, 0)
	shown := make([]*Header, 0)

	curr := t.first
	for {
		curr.hidden = slices.Contains(headers, curr.name)
		if curr.hidden {
			hidden = append(hidden, curr)
		} else {
			shown = append(shown, curr)
		}
		curr = curr.next
		if curr == t.first {
			break
		}
	}

	// mantain the column ording for shown values.
	slices.SortFunc(shown, func(a, b *Header) int {
		APlace := t.headerOrder[a.name]
		BPlace := t.headerOrder[b.name]
		return cmp.Compare(APlace, BPlace)
	})

	final := slices.Concat(shown, hidden)
	t.first = final[0]
	t.first.next = t.first
	t.first.prev = t.first
	println(t.first.name)
	for _, h := range final[1:] {
		println(h.name)
		t.first.appendHeader(h)
	}
}

func (t *Table) Size() (row int, col int) {
	row, col = t.rowCount, len(t.headerOrder)
	return
}

// WalkVisableValues 
func WalkVisableValues(t *Table, fn func(row, col int, c *Cell) bool) {
	
	currHeader := t.first
	var col, row int

	for {
		if currHeader.IsHidden() {
			return
		}

		currCell := currHeader.value
		for {
			if fn(row, col, currCell) {
				return
			}
			currCell = currCell.next
			row += 1
			if currCell == currHeader.value {
				break
			}
		}

		currHeader = currHeader.next
		row = 0
		col += 1

		if currHeader == t.first {
			return
		}
	}
}
