package viewmodel

import (
	"log"
	"slices"

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

func (ts *ColumnSettings) Headers() []string {
	fullFields := models.BookEntryFields()
	if ts.IsLoanHidden() {
		fullFields[models.IdxBorrower] = ""
		fullFields[models.IdxLoanedAt] = ""
	}
	if ts.IsReadHidden() {
		fullFields[models.IdxCompletedAt] = ""
		fullFields[models.IdxRating] = ""
	}
	headers := make([]string, 0)
	for _, h := range fullFields {
		if h != "" {
			headers = append(headers, h)
		}
	}
	return headers
}

func (ts *ColumnSettings) IsLoanHidden() bool {
	return isLoanHidden(ts.cfg)
}

func (ts *ColumnSettings) IsReadHidden() bool {
	return isReadHidden(ts.cfg)
}

func (ts *ColumnSettings) SetLoanHidden(t bool) {
	for idx, h := range ts.cfg.UI.Headers {
		switch idx {
		case models.IdxLoanedAt, models.IdxBorrower:
			h.IsHidden = t
			ts.cfg.UI.Headers[idx] = h
		}
	}
	ts.notify()
}

func (ts *ColumnSettings) SetReadHidden(t bool) {
	for idx, h := range ts.cfg.UI.Headers {
		switch idx {
		case models.IdxCompletedAt, models.IdxRating:
			h.IsHidden = t
			ts.cfg.UI.Headers[idx] = h
		}
	}
	ts.notify()
}

func (ts *ColumnSettings) SetWidth(label string, width float32) {
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
		return 200.0 // Todo: add constant or use config
	}
	h := ts.cfg.UI.Headers[idx]
	return h.Width
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

func (es *EntrySelected) Select(row, col int) {
	es.row = row
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
	}
	dt.service.AddListener(func() {
		dt.load()
	})
	dt.load()
	return dt
}

func (dt *DataTable) Size() (length, width int) {
	length = len(dt.data)
	width = len(models.BookEntryFields())
	if isLoanHidden(dt.cfg) {
		width -= 2
	} 
	if isReadHidden(dt.cfg) {
		width -= 2
	}
	return length, width
}

func (dt *DataTable) Get(row, col int) string {
	col = calcOffset(
		isLoanHidden(dt.cfg),
		isReadHidden(dt.cfg),
		col,
	)
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

	dt.data = dt.data[0:]
	clear(dt.rowToID)
	
	for row, book := range books {
		dt.rowToID[row] = book.ID
		values := entryValues(&book)
		dt.data = append(dt.data, values)
	}

	dt.notify()
}

func calcOffset(loanHidden, readHidden bool, col int) int {
	var (
		hasLoan bool = !loanHidden
		hasRead bool = !readHidden
	)
	if col <= models.IdxGenre {
		return col
	}

	if hasLoan && !hasRead {
		return col + 2
	}
	return col
}

func isLoanHidden(cfg *config.Config) bool {
	var (
		loaned bool
		borrower bool
	)
	for idx, h := range cfg.UI.Headers {
		switch idx {
		case models.IdxLoanedAt:
			loaned = h.IsHidden
		case models.IdxBorrower:
			borrower = h.IsHidden
		}
	}
	return loaned && borrower
}

func isReadHidden(cfg *config.Config) bool {
	var (
		rating bool
		completed bool
	)
	for idx, h := range cfg.UI.Headers {
		switch idx {
		case models.IdxRating:
			rating = h.IsHidden
		case models.IdxCompletedAt:
			completed = h.IsHidden
		}
	}
	return rating && completed
}

func entryValues(e *models.BookEntry) []string {
	
	headers := models.BookEntryFields()
	values  := make([]string, len(headers))

	for i, header := range headers {
		switch i {
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
