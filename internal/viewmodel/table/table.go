package table

import (
	"cmp"
	"errors"
	"log"
	"slices"
	"sync"
	"fmt"

	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/search"
	"github.com/dubbersthehoser/mayble/internal/viewmodel/display"
)

type BookData struct {
	data    []models.BookEntry
	rowToID map[int]int64
	idToRow map[int64]int
	mu      sync.RWMutex
}

func NewBookData(books []models.BookEntry) *BookData {
	bd := &BookData{
		rowToID: make(map[int]int64),
		idToRow: make(map[int64]int),
		data: books,
	}
	return bd
}

func (bd *BookData) Get(p Point) (string, error) {
	bd.mu.RLock()
	defer bd.mu.RUnlock()
	if p.Row >= len(bd.data) || p.Col < 0 {
		return "", errors.New("point.row out of range")
	}
	fields := display.EntryValues(&bd.data[p.Row])
	if p.Col >= len(fields) || p.Col < 0 {
		return "", errors.New("point.col out of range")
	}
	return fields[p.Col], nil
}

func (bd *BookData) RowToID(row int) (int64, error) {
	bd.mu.RLock()
	defer bd.mu.RUnlock()
	id, ok := bd.rowToID[row]
	if !ok {
		return 0, fmt.Errorf("row_to_id %d: row not found")
	}
	return id, nil
}

func (bd*BookData) IDToRow(id int64) (int, error) {
	bd.mu.RLock()
	defer bd.mu.RUnlock()
	row, ok := bd.idToRow[id]
	if !ok {
		return 0, fmt.Errorf("id_to_row %d: id not found", id)
	}
	return row, nil
}

func (bd*BookData) Size() (rows, cols int) {
	bd.mu.RLock()
	defer bd.mu.RUnlock()
	rows = len(bd.data)
	if rows == 0 {
		return 0, 0
	}
	cols = len(models.BookEntryFields())
	return
}

type Point struct {
	Col int
	Row int
}

func (p Point) IsSelectable() bool {
	return p.Col >= 0 && p.Row >= 0
}

type Table struct {
	sheet     *Sheet
	settings  *Settings
	sorting   *Sorting
	searching *Searching
	selected  *Selected

	OnNewSheet func()
	OnNewHeader func()

	eb *eventBus
}

func NewTable(cfg *config.Config, books []models.BookEntry) *Table {
	eb := newEventBus()
	t := &Table{
		sheet:       NewSheet(getHeaderSet(cfg), books),
		sorting:     NewSorting(eb, cfg.UI.TableSortBy, cfg.UI.TableAscending),
		settings:    NewSettings(cfg),
		selected:    newSelected(),
		OnNewSheet:  func() {},
		OnNewHeader: func() {},
	}
	return t
}

func (t *Table) NewSheet(books []models.BookEntry) {
	t.sheet = NewSheet(getHeaderSet(t.settings.cfg), books)
	t.OnNewSheet()
}

func (t *Table) ChangeHeader()


//
// Sheet
//

var sheetVersion int64

type Sheet struct {
	data         *BookData
	headerSet    []int
	sortedLookup map[int]int64
	
	version int64
}

func NewSheet(headerSet []int, books []models.BookEntry) *Sheet {
	sheetVersion += 1
	s := &Sheet{
		data:    NewBookData(books),
		sortedLookup: make(map[int]int64),
		headerSet: headerSet,
		version: sheetVersion,
	}

	clear(s.sortedLookup)
	for sRow := range books {
		id := books[sRow].ID
		s.sortedLookup[sRow] = id
	}

	return s
}

func (s *Sheet) RowToID(row int) (int64, bool) {
	id, ok := s.sortedLookup[row]
	return id, ok
}

func (s *Sheet) Get(p Point) (string, error) {
	p.Col = s.headerSet[p.Col]
	id := s.sortedLookup[p.Row]
	var err error
	p.Row, err = s.data.IDToRow(id)
	if err != nil {
		return "", err
	}
	return s.data.Get(p)
}

func (s *Sheet) Size() (rows, cols int) {
	rows, _  = s.data.Size()
	cols = len(s.headerSet)
	return
}

func (s *Sheet) Header() []string {
	header := make([]string, 0)
	for idx, label := range models.BookEntryFields() {
		if slices.Contains(s.headerSet, idx) {
			header = append(header, label)
		}
	}
	return header
}

//
// Sorting
//

type Sorting struct {
	eb *eventBus
	column string
	asc    bool
}

func NewSorting(eb *eventBus, column string, asc bool) *Sorting {
	s := &Sorting{
		eb: eb,
		column: column,
		asc: asc,
	}
	return s
}

func (s *Sorting) SetOrderBy(l string) {
	s.column = l
}

func (s *Sorting) SetAscending(t bool) {
	s.asc = t
}

func (s *Sorting) GetOrderBy() string {
	return s.column
}

func (s *Sorting) GetAscending() bool {
	return s.asc
}

func (s *Sorting) Sort() {
	s.eb.Notify(eventSort{asc: s.asc, column: s.column})
}

//
// Settings
//

type Settings struct {
	eb       *eventBus
	cfg      *config.Config
	onHidden []func()
}

func NewSettings(cfg *config.Config) *Settings {
	cs := &Settings{
		cfg: cfg,
	}
	return cs
}

func (ts *Settings) HeaderMinWidth() float32 {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.cfg.UI.TableMinWidth
}

func (ts *Settings) HeaderHeight() float32 {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.cfg.UI.TableHeaderHeight
}

func (ts *Settings) IsLoanHidden() bool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return isLoanHidden(ts.cfg)
}

func (ts *Settings) IsReadHidden() bool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return isReadHidden(ts.cfg)
}

func (ts *Settings) IsIDHidden() bool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return isIDHidden(ts.cfg)
}

func (ts *Settings) SetIDHidden(t bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	header := ts.cfg.UI.Headers[models.IdxID]
	header.IsHidden = t
	ts.cfg.UI.Headers[models.IdxID] = header
}

func (ts *Settings) ToggleHiddenID() bool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	header := ts.cfg.UI.Headers[models.IdxID]
	t := !header.IsHidden
	ts.SetIDHidden(t)
	return t
}

func (ts *Settings) SetLoanHidden(t bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	loaned := ts.cfg.UI.Headers[models.IdxLoanedAt]
	borrower := ts.cfg.UI.Headers[models.IdxBorrower]

	loaned.IsHidden = t
	borrower.IsHidden = t

	ts.cfg.UI.Headers[models.IdxLoanedAt] = loaned
	ts.cfg.UI.Headers[models.IdxBorrower] = borrower
}

func (ts *Settings) ToggleHiddenLoan() bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	loaned := ts.cfg.UI.Headers[models.IdxLoanedAt]
	borrower := ts.cfg.UI.Headers[models.IdxBorrower]
	t := !(loaned.IsHidden || borrower.IsHidden)
	ts.SetLoanHidden(t)
	return t
}

func (ts *Settings) SetReadHidden(t bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	rating := ts.cfg.UI.Headers[models.IdxRating]
	completed := ts.cfg.UI.Headers[models.IdxCompletedAt]

	rating.IsHidden = t
	completed.IsHidden = t

	ts.cfg.UI.Headers[models.IdxRating] = rating
	ts.cfg.UI.Headers[models.IdxCompletedAt] = completed
}

func (ts *Settings) ToggleHiddenRead() bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	rating := ts.cfg.UI.Headers[models.IdxRating]
	completed := ts.cfg.UI.Headers[models.IdxCompletedAt]
	t := !(completed.IsHidden || rating.IsHidden)
	ts.SetReadHidden(t)
	return t
}

func (ts *Settings) SetWidth(label string, width float32) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	idx := slices.Index(models.BookEntryFields(), label)
	if idx == -1 {
		log.Printf("Error: invalid header label '%s'", label)
		return
	}
	h, ok := ts.cfg.UI.Headers[idx]
	if !ok {
		log.Printf("Warning: column not found '%s'", label)
		return
	}
	h.Width = width
	ts.cfg.UI.Headers[idx] = h
}

func (ts *Settings) GetWidth(label string) float32 {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	idx := slices.Index(models.BookEntryFields(), label)
	if idx == -1 {
		log.Printf("Error: invalid header label '%s'", label)
		return ts.cfg.UI.TableMinWidth
	}
	h := ts.cfg.UI.Headers[idx]
	width := h.Width
	return width
}

//
// Searching
//

const columnAll = "All"

type Searchable struct {
	table *Table
}

func NewSearchable(t *Table) *Searchable {
	s := &Searchable{
		table: t,
	}

	return s
}

func (s *Searchable) Options() []string {
	o := []string{
		"All",
	}
	o = append(o, s.table.sheet.Header()...)
	return o
}

func (s *Searchable) SetBy(c string) {
	var col int
	if c == "All" {
		col = -1
	} else {
		col = slices.Index(s.table.sheet.Header(), c)
	}

	s.table.eb.Notify(eventSearchBy{
		column: col,
	})
}

type Searching struct {

	sheet *Sheet

	eb *eventBus

	byColumn string

	// The row that is selected from the scored matches.
	row    int
	scored []Point

	l []func()
}

func NewSearching(s *Sheet, eb *eventBus) *Searching {
	sr := &Searching{
		sheet:      s,
		eb: eb,
	}
	return sr
}

func (s *Searching) has() bool {
	return len(s.scored) != 0
}

func (s *Searching) ByColumn(label string) {
	s.byColumn = label
}

func (s *Searching) Prev() {
	if !s.has() {
		return
	}
	s.row -= 1
	if s.row < 0 {
		s.row = len(s.scored) - 1
	}
	s.eb.Notify(eventSelected{
		has: true,
		point: s.scored[s.row],
	})
}

func (s *Searching) Next() {
	s.row += 1
	if !s.has() {
		return
	}
	if s.row == len(s.scored) {
		s.row = 0
	}
	s.eb.Notify(eventSelected{
		has: true,
		point: s.scored[s.row],
	})
}

func (s *Searching) search(sh *Sheet, by string, pattern string) {
	var trv search.Traverser
	if s.byColumn == columnAll {
		trv = newTableTraverse(sh)
	} else {
		idx := slices.Index(sh.Header(), by)
		trv = newColumnTraverse(sh, idx)
	}
	srch := (&search.Searcher{}).Set(trv, pattern)
	s.searchToScored(srch)

	if !s.has() {
		s.eb.Notify(eventSelected{
			has: false,
		})
	}

	s.row = -1
	s.Next()
}

func (s *Searching) searchToScored(srch *search.Searcher) {
	type result struct {
		row, col, score int
	}
	results := make([]result, 0)
	for srch.Next() {
		point := srch.Point()
		score := srch.Score()
		if score == -1 {
			continue
		}
		r := result{
			row:   point.Row,
			col:   point.Col,
			score: srch.Score(),
		}
		results = append(results, r)
	}
	if len(results) == 0 {
		s.row = 0
		s.scored = s.scored[:0]
	}

	slices.SortFunc(results, func(a, b result) int {
		r := cmp.Compare(a.score, b.score)
		if r == 0 {
			return cmp.Compare(a.row, b.row)
		}
		return r * -1
	})
	s.row = 0
	s.scored = s.scored[:0]
	for _, r := range results {
		p := Point{
			Row: r.row,
			Col: r.col,
		}
		s.scored = append(s.scored, p)
	}
}

//
// Selected
//

type Selected struct {
	eb    *eventBus
}

func newSelected() *Selected {
	es := &Selected{
		eb: newEventBus(),
	}
	return es
}

func (es *Selected) Select(p Point) {
	es.eb.Notify(eventSelected{point: p, has: true})
}
func (es *Selected) Unselect() {
	es.eb.Notify(eventSelected{has: false})
}

func (es *Selected) AddListener(fn func(has bool, p Point)) {
	es.eb.Subscribe("eventSelected", func(v any) {
		e := v.(eventSelected)
		fn(e.has, e.point)
	})
}

func isLoanHidden(cfg *config.Config) bool {
	loaned := cfg.UI.Headers[models.IdxLoanedAt]
	borrower := cfg.UI.Headers[models.IdxBorrower]
	return loaned.IsHidden && borrower.IsHidden

}

func isIDHidden(cfg *config.Config) bool {
	header := cfg.UI.Headers[models.IdxID]
	return header.IsHidden
}

func isReadHidden(cfg *config.Config) bool {
	rating := cfg.UI.Headers[models.IdxRating]
	completed := cfg.UI.Headers[models.IdxCompletedAt]
	return rating.IsHidden && completed.IsHidden
}

func removeHiddenColumns(cfg *config.Config) []int {
	indexs := make([]int, 0)
	if isLoanHidden(cfg) {
		indexs = append(indexs, models.IdxBorrower, models.IdxLoanedAt)
	}
	if isReadHidden(cfg) {
		indexs = append(indexs, models.IdxRating, models.IdxCompletedAt)
	}
	if isIDHidden(cfg) {
		indexs = append(indexs, models.IdxID)
	}
	return indexs
}

func includedColumns(cfg *config.Config) []int {
	indexes := make([]int, 0 )
	for idx, h := range cfg.UI.Headers {
		if !h.IsHidden{
			indexes = append(indexes, idx)
		}
	}
	return indexes
}

func getHeaderSet(cfg *config.Config) []int {
	set := make([]int, 0)
	for idx, h := range cfg.UI.Headers {
		if !h.IsHidden {
			set = append(set, idx)
		}
	}
	return set
}

func columnLookaside(cfg *config.Config, col int) (int, error) {
	count := -1
	for idx, h := range cfg.UI.Headers {
		if !h.IsHidden {
			count += 1
		}

		if col == count {
			return idx, nil
		}
	}
	return 0, fmt.Errorf("lookaside %d: invalid column", col)
}

