package table

import (
	"errors"
	"log"
	"slices"
	"cmp"

	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/search"
	"github.com/dubbersthehoser/mayble/internal/viewmodel/display"
)

type Point struct {
	Col int
	Row int
}

func (p Point) IsSelectable() bool {
	return p.Col >= 0 && p.Row >= 0
}

type Table struct {
	Sheet     *Sheet
	Settings  *Settings
	Sorting   *Sorting
	Searching *Searching
	Selected  *Selected
}

func NewTable(cfg *config.Config, srv *app.Service) *Table {
	t := &Table{
		Sheet: NewSheet(cfg, srv),
		Sorting: NewSorting(cfg),
		Settings: NewSettings(cfg),
		Selected: newSelected(),
	}

	t.Searching = NewSearching()

	t.Searching.AddListener(func() {
		if t.Searching.Has() {
			t.Selected.Select(t.Searching.Selected())
		} else {
			t.Selected.Unselect()
		}
	})

	t.Sorting.AddListener(func() {
		t.Sheet.Load()
		t.Selected.Unselect()
	})

	t.Settings.AddListener(func() {
		t.Sheet.Load()
		t.Selected.Unselect()
	})

	return t
}

func (t *Table) Search(s string) {
	if s == "" {
		t.Selected.Unselect()
		return
	}
	t.Searching.Search(t.Sheet.data, s)
}


//
// Sheet
//

type Sheet struct {
	srv     *app.Service
	cfg     *config.Config
	data    [][]string

	rowToID map[int]int64

	l []func()
}

func NewSheet(cfg *config.Config, srv *app.Service) *Sheet {
	s := &Sheet{
		srv: srv,
		cfg: cfg,
		data: make([][]string, 0),
		rowToID: make(map[int]int64),
	}
	return s
}

func (s *Sheet) RowToID(row int) (int64, bool) {
	v, ok := s.rowToID[row]
	return v, ok
}

func (s *Sheet) Get(p Point) (string, error) {
	if p.Row >= len(s.data) || p.Row < 0 {
		return "", errors.New("row point out of range")
	}
	if p.Col >= len(s.data[p.Row]) || p.Col < 0 {
		return "", errors.New("column point out of range")
	}
	return s.data[p.Row][p.Col], nil
}

func (s *Sheet) Size() (length, width int) {
	
	headers := models.BookEntryFields()

	for _, idx := range removeHiddenColumns(s.cfg) {
		headers = slices.Delete(headers, idx, idx+1)
	}
	
	width = len(headers)
	length = len(s.data)
	return length, width
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

	books, err := s.srv.GetAllBooks()
	if err != nil {
		return err
	}

	if err := app.SortBooks(books, by, asc); err != nil {
		return err
	}


	s.data = s.data[:0]
	clear(s.rowToID)

	for row, book := range books {
		s.rowToID[row] = book.ID
		values := display.EntryValues(&book)
		for _, idx := range removeHiddenColumns(s.cfg) {
			values = slices.Delete(values, idx, idx+1)
		}
		s.data = append(s.data, values)
	}
	
	s.notify()
	return nil
}

func (s *Sheet) AddListener(fn func() ) {
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
	cfg *config.Config
	l   []func()
}

func NewSettings(cfg *config.Config) *Settings {
	cs := &Settings{
		cfg: cfg,
	}
	return cs
}

func (ts *Settings) MinWidth() float32 {
	return ts.cfg.UI.TableMinWidth
}

func (ts *Settings) Headers() []string {
	headers := models.BookEntryFields()
	removeIdxs := removeHiddenColumns(ts.cfg)
	for _, idx := range removeIdxs {
		headers = slices.Delete(headers, idx, idx+1)
	}
	return headers
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
	ts.notify()
}

func (ts *Settings) SetLoanHidden(t bool) {

	loaned := ts.cfg.UI.Headers[models.IdxLoanedAt]
	borrower := ts.cfg.UI.Headers[models.IdxBorrower]

	loaned.IsHidden = t
	borrower.IsHidden = t

	ts.cfg.UI.Headers[models.IdxLoanedAt] = loaned
	ts.cfg.UI.Headers[models.IdxBorrower] = borrower

	ts.notify()
}

func (ts *Settings) SetReadHidden(t bool) {
	rating := ts.cfg.UI.Headers[models.IdxRating]
	completed := ts.cfg.UI.Headers[models.IdxCompletedAt]

	rating.IsHidden = t
	completed.IsHidden = t

	ts.cfg.UI.Headers[models.IdxRating] = rating
	ts.cfg.UI.Headers[models.IdxCompletedAt] = completed

	ts.notify()
}

func (ts *Settings) SetWidth(label string, width float32) {
	if width <= ts.cfg.UI.TableMinWidth {
		width = ts.cfg.UI.TableMinWidth
	}
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
	if width <= ts.cfg.UI.TableMinWidth {
		width = ts.cfg.UI.TableMinWidth
	}
	return width
}

// AddListener for changes to hidden columns.
func (ts *Settings) AddListener(fn func()) {
	if ts.l == nil {
		ts.l = make([]func(), 0)
	}
	ts.l = append(ts.l, fn)
}

func (ts *Settings) notify() {
	for _, fn := range ts.l {
		fn()
	}
}

//
// Searching
//

const columnAll = -1

type Searching struct {

	cellSearch  search.CellSearch
	tableSearch search.TableSearch

	// The column to be searched.
	column int

	// The row that is selected out of scored.
	row    int
	scored []Point

	l []func()
}

func NewSearching() *Searching {
	sr := &Searching{}
	return sr
}

func (s *Searching) GetOptions() []string {
	// TODO: Update this to only allow columns that are shown.
	return []string{
		"All",
		models.BookEntryFields()[models.IdxTitle],
		models.BookEntryFields()[models.IdxAuthor],
		models.BookEntryFields()[models.IdxGenre],
		models.BookEntryFields()[models.IdxBorrower],
		models.BookEntryFields()[models.IdxLoanedAt],
		models.BookEntryFields()[models.IdxCompletedAt],
	}
}

func (s *Searching) SetBy(c string) {
	if c == "All" {
		s.column = columnAll
		return
	}

	s.column = slices.Index(models.BookEntryFields(), c)
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

func (s *Searching) Search(data [][]string, search string) {

	if s.column == columnAll {
		s.searchAll(data, search)
	} else {
		s.searchColumn(data, search)
	}
	s.notify()
}

func (s *Searching) searchColumn(data [][]string, search string) {
	dataCol := make([]string, 0)
	for _, row := range data {
		dataCol = append(dataCol, row[s.column])
	}
	type result struct {
		row, score int
	}
	results := make([]result, 0)
	s.cellSearch.Set(dataCol, search)
	for s.cellSearch.Next() {
		row := s.cellSearch.Pos()
		score := s.cellSearch.Score()
		if score == -1 {
			continue
		}
		r := result{
			row:   row,
			score: score,
		}
		results = append(results, r)
	}

	if len(results) == 0 {
		s.scored = s.scored[:0]
		s.row = 0
		return
	}
	slices.SortFunc(results, func(a, b result) int {
		r := cmp.Compare(a.score, b.score)
		if r == 0 {
			return cmp.Compare(a.row, b.row)
		}
		return cmp.Compare(a.score, b.score)
	})

	s.row = 0
	s.scored = s.scored[:0]
	for _, r := range results {
		p := Point{
			Row: r.row,
			Col: s.column,
		}
		s.scored = append(s.scored, p)
	}
}

func (s *Searching) searchAll(data [][]string, search string) {
	type result struct {
		row, col, score int
	}
	results := make([]result, 0)
	s.tableSearch.Set(data, search)
	for s.tableSearch.Next() {
		row, col := s.tableSearch.Pos()
		score := s.tableSearch.Score()
		if score == -1 {
			continue
		}
		r := result{
			row:   row,
			col:   col,
			score: score,
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

type Selected struct {
	point Point
	l   []func()
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

func (es *Selected) Get() Point {
	return es.point
}

func (es *Selected) Has() bool {
	return es.point.IsSelectable()
}

func (es *Selected) AddListener(fn func()) {
	if es.l == nil {
		es.l = make([]func(), 0)
	}
	es.l = append(es.l, fn)
}

func (es *Selected) notify() {
	for _, fn := range es.l {
		fn()
	}
}

func AllowedSearchOptions(options, headers []string) []string {
	o := make([]string, 0)
	for _, option := range options {
		if slices.Contains(headers, option) || option == "All" {
			o = append(o, option)
		}
	}
	return o
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
