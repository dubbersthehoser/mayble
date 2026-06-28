package viewmodel

import (
	"log"
	"slices"
	"fmt"

	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/models"
)

type SortingTable struct {
	cfg *config.Config
	
	l []func()
}

func newSortingTable(cfg *config.Config) *SortingTable {
	s := &SortingTable{
		cfg: cfg,
	}
	return s
}

func (s *SortingTable) SetOrderBy(l string) {
	idx := slices.Index(models.BookEntryFields(), l)
	if idx == -1 {
		log.Printf("Error: invalid header lable '%s'", l)
		return
	}
	s.cfg.UI.TableSortBy = idx
}

func (s *SortingTable) SetAscending(t bool) {
	s.cfg.UI.TableAscending = t
}

func (s *SortingTable) GetOrderBy() string {
	return models.BookEntryFields()[s.cfg.UI.TableSortBy]
}

func (s *SortingTable) GetAscending() bool {
	return s.cfg.UI.TableAscending
}

func (s *SortingTable) Sort() {
	s.notify()
}

func (s *SortingTable) AddListener(fn func()) {
	if s.l == nil {
		s.l = make([]func(), 0)
	}
	s.l = append(s.l, fn)
}

func (s *SortingTable) notify() {
	for _, fn := range s.l {
		fn()
	}
}

type ColumnSettings struct {
	cfg *config.Config
	l []func()
}

func newColumnSettings(cfg *config.Config) *ColumnSettings {
	cs := &ColumnSettings{
		cfg: cfg,
	}
	return cs
}

func (ts *ColumnSettings) MinWidth() float32 {
	return ts.cfg.UI.TableMinWidth
}

func (ts *ColumnSettings) Headers() []string {
	headers := models.BookEntryFields()
	removeIdxs := removeHiddenColumns(ts.cfg)
	for _, idx := range removeIdxs {
		headers = slices.Delete(headers, idx, idx+1)
	}
	return headers
}

func (ts *ColumnSettings) IsLoanHidden() bool {
	return isLoanHidden(ts.cfg)
}

func (ts *ColumnSettings) IsReadHidden() bool {
	return isReadHidden(ts.cfg)
}

func (ts *ColumnSettings) IsIDHidden() bool {
	return isIDHidden(ts.cfg)
}

func (ts *ColumnSettings) SetIDHidden(t bool) {
	header := ts.cfg.UI.Headers[models.IdxID]
	header.IsHidden = t
	ts.cfg.UI.Headers[models.IdxID] = header
	ts.notify()
}

func (ts *ColumnSettings) SetLoanHidden(t bool) {
	
	loaned := ts.cfg.UI.Headers[models.IdxLoanedAt]
	borrower := ts.cfg.UI.Headers[models.IdxBorrower]

	loaned.IsHidden = t
	borrower.IsHidden = t

	ts.cfg.UI.Headers[models.IdxLoanedAt] = loaned
	ts.cfg.UI.Headers[models.IdxBorrower] = borrower

	ts.notify()
}

func (ts *ColumnSettings) SetReadHidden(t bool) {
	rating := ts.cfg.UI.Headers[models.IdxRating]
	completed := ts.cfg.UI.Headers[models.IdxCompletedAt]

	rating.IsHidden = t
	completed.IsHidden = t

	ts.cfg.UI.Headers[models.IdxRating] = rating
	ts.cfg.UI.Headers[models.IdxCompletedAt] = completed

	ts.notify()
}

func (ts *ColumnSettings) SetWidth(label string, width float32) {
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

func (ts *ColumnSettings) GetWidth(label string) float32 {
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

// AddListener listen for changes to hidden columns.
func (ts *ColumnSettings) AddListener(fn func()) {
	if ts.l == nil {
		ts.l = make([]func(), 0)
	}
	ts.l = append(ts.l, fn)
}

func (ts *ColumnSettings) notify() {
	for _, fn := range ts.l {
		fn()
	}
}

type EntrySelected struct {
	row int
	col int
	l []func()
}

func newEntrySelected() *EntrySelected {
	es := &EntrySelected{
		row: -1,
		col: -1,
	}
	return es
}

func (es *EntrySelected) Select(row, col int) {
	es.row = row
	es.col = col
	es.notify()
}
func (es *EntrySelected) Unselect() {
	es.row = -1
	es.notify()
}

func (es *EntrySelected) Get() (int, int) {
	return es.row, es.col
}

func (es *EntrySelected) Has() bool {
	return es.row != -1
}

func (es *EntrySelected) AddListener(fn func()) {
	if es.l == nil {
		es.l = make([]func(), 0)
	}
	es.l = append(es.l, fn)
}

func (es *EntrySelected) notify() {
	for _, fn := range es.l {
		fn()
	}
}

type DataTable struct {
	service *app.Service
	data    [][]string
	rowToID map[int]int64

	cfg     *config.Config

	l []func()
}

func newDataTable(cfg *config.Config, s *app.Service) *DataTable {
	dt := &DataTable{
		service: s,
		cfg: cfg,
		rowToID: make(map[int]int64),
	}
	dt.service.AddListener(func() {
		dt.load()
	})
	dt.load()
	return dt
}

func (dt *DataTable) Size() (length, width int) {
	fields := models.BookEntryFields()
	for _, idx := range removeHiddenColumns(dt.cfg) {
		fields = slices.Delete(fields, idx, idx+1)
	}
	length = len(dt.data)
	width = len(fields)
	return length, width
}

func (dt *DataTable) Get(row, col int) string {
	return dt.data[row][col]
}

// AddListener listen for (re)loads.
func (dt *DataTable) AddListener(fn func()) {
	if dt.l == nil {
		dt.l = make([]func(), 0)
	}
	dt.l = append(dt.l, fn)
}

func (dt *DataTable) notify() {
	for _, fn := range dt.l {
		fn()
	}
}

func (dt *DataTable) load() {

	by := dt.cfg.UI.TableSortBy
	asc := dt.cfg.UI.TableAscending

	books, err := dt.service.GetAllBooks()
	if err != nil {
		log.Println("Error:", err)
		return
	}

	if err := app.SortBooks(books, by, asc); err != nil {
		log.Println("Error:", err)
		return
	}

	dt.data = dt.data[:0]
	clear(dt.rowToID)
	
	for row, book := range books {
		dt.rowToID[row] = book.ID
		values := entryValues(&book)
		for _, idx := range removeHiddenColumns(dt.cfg) {
			values = slices.Delete(values, idx, idx+1)
		}
		dt.data = append(dt.data, values)
	}

	dt.notify()
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

func entryValues(e *models.BookEntry) []string {
	
	headers := models.BookEntryFields()
	values  := make([]string, len(headers))

	for i, header := range headers {
		switch i {
		case models.IdxID:
			values[i] = fmt.Sprintf("%d", e.ID)
		case models.IdxTitle:
			values[i] = e.Title
		case models.IdxAuthor:
			values[i] = e.Author
		case models.IdxGenre:
			values[i] = e.Genre
		case models.IdxRating:
			values[i] = formatRating(e.Rating)
		case models.IdxCompletedAt:
			values[i] = formatDate(&e.CompletedAt)
		case models.IdxBorrower:
			values[i] = e.Borrower
		case models.IdxLoanedAt:
			values[i] = formatDate(&e.LoanedAt)
		default:
			panic("unknown field:" + header)
		}
	}
	return values
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

