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

type GetAllBooks func() ([]models.BookEntry, error)

type BookData struct {
	data    []models.BookEntry
	rowToID map[int]int64
	idToRow map[int64]int
	mu      sync.RWMutex
}

func NewBookData() *BookData {
	bd := &BookData{
		rowToID: make(map[int]int64),
		idToRow: make(map[int64]int),
		data: make([]models.BookEntry, 0),

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

func (bd *BookData) Set(data []models.BookEntry)  {
	bd.mu.Lock()
	defer bd.mu.Unlock()
	bd.data = data
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

	// muSheet mutex for when Sheet is being touched by mutiple threads.
	// Inpaticular the main UI thread, search, and recreating Sheet threads.
	// NOTE:(dth): Best for [view.Table] be the only one to call [Sheet] methods.
	muSheet sync.RWMutex
}

func NewTable(cfg *config.Config, fetch GetAllBooks) *Table {
	t := &Table{
		sheet:    NewSheet(cfg, fetch),
		sorting:  NewSorting(cfg),
		settings: NewSettings(cfg),
		selected: newSelected(),
	}

	t.searching = NewSearching(t)

	t.searching.AddListener(func() {
		if t.searching.Has() {
			t.selected.Select(t.searching.Selected())
		} else {
			t.selected.Unselect()
		}
	})

	t.sorting.AddListener(func() {
		t.sheet.Load()
		t.selected.Unselect()
	})

	t.settings.AddOnHidden(func() {
		t.selected.Unselect()
	})

	return t
}

func (t *Table) NewSheet() {
	t.muSheet.Lock()
	defer t.muSheet.Unlock()
	old := t.sheet
	t.sheet = NewSheet(old.cfg, old.fetch)
	t.sheet.Load()
}

func (t *Table) GetDataPoint(p Point) (string, error) {
	t.muSheet.Lock()
	defer t.muSheet.Unlock()
	return t.sheet.Get(p)
}

func (t *Table) Search(s string) {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	if s == "" {
		t.selected.Unselect()
		return
	}
	t.searching.search(s)
}

func (t *Table) NextSearchedItem() {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	t.searching.Next()
}

func (t *Table) PrevSearchedItem() {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	t.searching.Prev()
}

func (t *Table) SearchableOptions() []string {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	return t.searching.Searchable.GetOptions()
}

func (t *Table) SearchableSetBy(s string) {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	t.searching.Searchable.SetBy(s)
}

func (t *Table) AddListenerOnHidden(h func()) {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	t.settings.AddOnHidden(h)
}

func (t *Table) Size() (int, int) {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	return t.sheet.Size()
}

func (t *Table) Header() []string {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	return t.sheet.Header()
}

func (t *Table) GetColumnWidth(label string) float32 {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	return t.settings.GetWidth(label)
}

func (t *Table) SetColumnWidth(label string, width float32) {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	t.settings.SetWidth(label, width)
}

func (t *Table) GetSortingOrderBy() string {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	return t.sorting.GetOrderBy()
}
func (t *Table) SetSortingOrderBy(by string) {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	t.sorting.SetOrderBy(by)
}

func (t *Table) GetSortingAscending() bool {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	return t.sorting.GetAscending()
}

func (t *Table) SetSortingAscending(b bool) {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	t.sorting.SetAscending(b)
}

func (t *Table) Sort() {
	t.muSheet.Lock()
	defer t.muSheet.Unlock()
	t.sorting.Sort()
}

func (t *Table) SelectPoint(p Point) {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	t.selected.Select(p)
}

func (t *Table) UnselectPoint() {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	t.selected.Unselect()
}

func (t *Table) AddSelectedListener(h func(has bool, p Point)) {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	t.selected.AddListener(h)
}

func (t *Table) AddListenerOnDataChanged(h func()) {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	t.sheet.AddListener(h)
}

func (t *Table) HeaderHeight() float32 {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	return t.settings.HeaderHeight()
}

func (t *Table) HeaderMinWidth() float32 {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	return t.settings.HeaderMinWidth()
}

func (t *Table) NotifyOnColumnHidden() {
	t.muSheet.RLock()
	defer t.muSheet.RUnlock()
	t.settings.NotifyOnHidden()
}

//
// Sheet
//

type Sheet struct {
	fetch GetAllBooks
	cfg   *config.Config
	data  *BookData

	sortedLookup map[int]int64

	l []func()
}

func NewSheet(cfg *config.Config, fetch GetAllBooks) *Sheet {
	s := &Sheet{
		fetch:   fetch,
		cfg:     cfg,
		data:    NewBookData(),
		sortedLookup: make(map[int]int64),
	}
	return s
}

func (s *Sheet) RowToID(row int) (int64, bool) {
	id, ok := s.sortedLookup[row]
	return id, ok
}

func (s *Sheet) Get(p Point) (string, error) {
	col, err := columnLookaside(s.cfg, p.Col)
	if err != nil {
		return "", err
	}
	p.Col = col
	id := s.sortedLookup[p.Row]
	p.Row, err = s.data.IDToRow(id)
	if err != nil {
		return "", err
	}
	return s.data.Get(p)
}

func (s *Sheet) Size() (rows, cols int) {
	rows, _  = s.data.Size()
	cols = len(includedColumns(s.cfg))
	return
}

func (s *Sheet) Header() []string {
	header := models.BookEntryFields()
	removeIdx := removeHiddenColumns(s.cfg)
	for _, idx := range removeIdx {
		header = slices.Delete(header, idx, idx+1)
	}
	return header
}

func (s *Sheet) Load() error {

	by := s.cfg.UI.TableSortBy
	asc := s.cfg.UI.TableAscending

	books, err := s.fetch()
	if err != nil {
		return err
	}

	// load the unsorted data, sort, then add to sort lookup.
	_ = s.data.Load(books)
	if err := app.SortBooks(books, by, asc); err != nil {
		return err
	}
	clear(s.sortedLookup)
	for sRow := range books {
		id := books[sRow].ID
		s.sortedLookup[sRow] = id
	}

	s.notify()
	return nil
}

func (s *Sheet) AddListener(fn func()) {
	if s.l == nil {
		s.l = make([]func(), 0)
	}
	s.l = append(s.l, fn)
}

func (s *Sheet) notify() {
	for _, fn := range s.l {
		fn()
	}
}

//
// Sorting
//

type Sorting struct {
	cfg *config.Config

	l []func()
}

func NewSorting(cfg *config.Config) *Sorting {
	s := &Sorting{
		cfg: cfg,
	}
	return s
}

func (s *Sorting) SetOrderBy(l string) {
	idx := slices.Index(models.BookEntryFields(), l)
	if idx == -1 {
		return
	}
	s.cfg.UI.TableSortBy = idx
}

func (s *Sorting) SetAscending(t bool) {
	s.cfg.UI.TableAscending = t
}

func (s *Sorting) GetOrderBy() string {
	return models.BookEntryFields()[s.cfg.UI.TableSortBy]
}

func (s *Sorting) GetAscending() bool {
	return s.cfg.UI.TableAscending
}

func (s *Sorting) Sort() {
	s.notify()
}

// AddListener for Sort call.
func (s *Sorting) AddListener(fn func()) {
	if s.l == nil {
		s.l = make([]func(), 0)
	}
	s.l = append(s.l, fn)
}

func (s *Sorting) notify() {
	for _, fn := range s.l {
		fn()
	}
}

//
// Settings
//

type Settings struct {
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
	return ts.cfg.UI.TableMinWidth
}

func (ts *Settings) HeaderHeight() float32 {
	return ts.cfg.UI.TableHeaderHeight
}

func (ts *Settings) IsLoanHidden() bool {
	return isLoanHidden(ts.cfg)
}

func (ts *Settings) IsReadHidden() bool {
	return isReadHidden(ts.cfg)
}

func (ts *Settings) IsIDHidden() bool {
	return isIDHidden(ts.cfg)
}

func (ts *Settings) SetIDHidden(t bool) {
	header := ts.cfg.UI.Headers[models.IdxID]
	header.IsHidden = t
	ts.cfg.UI.Headers[models.IdxID] = header
	ts.NotifyOnHidden()
}

func (ts *Settings) ToggleHiddenID() bool {
	header := ts.cfg.UI.Headers[models.IdxID]
	t := !header.IsHidden
	ts.SetIDHidden(t)
	return t
}

func (ts *Settings) SetLoanHidden(t bool) {

	loaned := ts.cfg.UI.Headers[models.IdxLoanedAt]
	borrower := ts.cfg.UI.Headers[models.IdxBorrower]

	loaned.IsHidden = t
	borrower.IsHidden = t

	ts.cfg.UI.Headers[models.IdxLoanedAt] = loaned
	ts.cfg.UI.Headers[models.IdxBorrower] = borrower

	ts.NotifyOnHidden()
}

func (ts *Settings) ToggleHiddenLoan() bool {
	loaned := ts.cfg.UI.Headers[models.IdxLoanedAt]
	borrower := ts.cfg.UI.Headers[models.IdxBorrower]
	t := !(loaned.IsHidden || borrower.IsHidden)
	ts.SetLoanHidden(t)
	return t
}

func (ts *Settings) SetReadHidden(t bool) {
	rating := ts.cfg.UI.Headers[models.IdxRating]
	completed := ts.cfg.UI.Headers[models.IdxCompletedAt]

	rating.IsHidden = t
	completed.IsHidden = t

	ts.cfg.UI.Headers[models.IdxRating] = rating
	ts.cfg.UI.Headers[models.IdxCompletedAt] = completed

	ts.NotifyOnHidden()
}

func (ts *Settings) ToggleHiddenRead() bool {
	rating := ts.cfg.UI.Headers[models.IdxRating]
	completed := ts.cfg.UI.Headers[models.IdxCompletedAt]
	t := !(completed.IsHidden || rating.IsHidden)
	ts.SetReadHidden(t)
	return t
}

func (ts *Settings) SetWidth(label string, width float32) {
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
	idx := slices.Index(models.BookEntryFields(), label)
	if idx == -1 {
		log.Printf("Error: invalid header label '%s'", label)
		return ts.cfg.UI.TableMinWidth
	}
	h := ts.cfg.UI.Headers[idx]
	width := h.Width
	return width
}

// AddOnHidden for changes to hidden columns.
func (ts *Settings) AddOnHidden(fn func()) {
	if ts.onHidden == nil {
		ts.onHidden = make([]func(), 0)
	}
	ts.onHidden = append(ts.onHidden, fn)
}

func (ts *Settings) NotifyOnHidden() {
	for _, fn := range ts.onHidden {
		fn()
	}
}

//
// Searching
//

const columnAll = -1

type Searchable struct {
	table *Table

	column int

	l []func()
}

func NewSearchable(t *Table) *Searchable {
	s := &Searchable{
		table: t,
	}

	t.settings.AddOnHidden(func() {
		s.notify()
	})
	return s
}

func (s *Searchable) GetOptions() []string {
	o := []string{
		"All",
	}
	o = append(o, s.table.sheet.Header()...)
	return o
}

func (s *Searchable) SetBy(c string) {
	if c == "All" {
		s.column = columnAll
		return
	}
	s.column = slices.Index(s.table.sheet.Header(), c)
}

func (s *Searchable) AddListener(fn func()) {
	if s.l == nil {
		s.l = make([]func(), 0)
	}
	s.l = append(s.l, fn)
}

func (s *Searchable) notify() {
	for _, fn := range s.l {
		fn()
	}
}

type Searching struct {
	table *Table

	Searchable *Searchable

	// The row that is selected from the scored matches.
	row    int
	scored []Point

	l []func()
}

func NewSearching(t *Table) *Searching {
	sr := &Searching{
		table:      t,
		Searchable: NewSearchable(t),
	}
	return sr
}

func (s *Searching) Selected() Point {
	return s.scored[s.row]
}

func (s *Searching) Has() bool {
	return len(s.scored) != 0
}

func (s *Searching) Prev() {
	s.row -= 1
	if s.row < 0 {
		s.row = len(s.scored) - 1
	}
	s.notify()
}

func (s *Searching) Next() {
	s.row += 1
	if s.row == len(s.scored) {
		s.row = 0
	}
	s.notify()
}

func (s *Searching) AddListener(fn func()) {
	if s.l == nil {
		s.l = make([]func(), 0)
	}
	s.l = append(s.l, fn)
}

func (s *Searching) notify() {
	for _, fn := range s.l {
		fn()
	}
}

func (s *Searching) search(pattern string) {
	data := s.table.sheet.data

	var trv search.Traverser
	if s.Searchable.column == columnAll {
		trv = (&search.TableTraverse{}).Set(data)
	} else {
		trv = (&search.ColumnTraverse{}).Set(data, s.Searchable.column)
	}
	srch := (&search.Searcher{}).Set(trv, pattern)
	s.searchToScored(srch)
	s.notify()
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
	point Point
	l     []func(has bool, p Point)
}

func newSelected() *Selected {
	es := &Selected{
		point: Point{
			Row: -1,
			Col: -1,
		},
	}
	return es
}

func (es *Selected) Select(p Point) {
	es.point = p
	es.notify()
}
func (es *Selected) Unselect() {
	es.point.Row = -1
	es.point.Col = -1
	es.notify()
}

func (es *Selected) get() Point {
	return es.point
}

func (es *Selected) has() bool {
	return es.point.IsSelectable()
}

func (es *Selected) AddListener(fn func(has bool, p Point)) {
	if es.l == nil {
		es.l = make([]func(has bool, p Point), 0)
	}
	es.l = append(es.l, fn)
}

func (es *Selected) notify() {
	for _, fn := range es.l {
		fn(es.has(), es.get())
	}
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

