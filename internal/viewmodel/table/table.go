package table

import (
	"errors"
	"log"
	"slices"

	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/config"

)

type Point struct {
	Col int
	Row int
}

type Sheet struct {
	srv  *app.Service
	cfg  *config.Config
	data [][]string
}

func (s *Sheet) Get(p Point) (string, error) {
	if p.Row >= len(s.data) || p.Row > 0 {
		return "", errors.New("row point out of range")
	}
	if p.Col >= len(s.data[p.Row]) || p.Col < 0 {
		return "", errors.New("column point out of range")
	}
	return s.data[p.Row][p.Col], nil
}

func (s *Sheet) Size() (width, length int) {
	width = len(s.data)
	if width > 0 {
		length = len(s.data[0])
	}
	return width, length
}

func (s *Sheet) Header() []string {
	return nil
}

func (s *Sheet) Load() error {
	
	return nil
}


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
		log.Printf("Error: invalid header lable '%s'", l)
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

type ColumnSettings struct {
	cfg *config.Config
	l   []func()
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
