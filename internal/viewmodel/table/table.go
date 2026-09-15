package table

import (
	"log"
	"fmt"
	"slices"
	"time"

	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/event"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/search"
	"github.com/dubbersthehoser/mayble/internal/snapshot"
	"github.com/dubbersthehoser/mayble/internal/command"
	"github.com/dubbersthehoser/mayble/internal/worker"
)

const ColumnAll = "All"

type Table struct {
	worker *worker.Worker
	eb     *event.EventBus

	Searching       *Searching
	Selected        *Selected
	Sorting         *Sorting
	Sheet           *Sheet
	Searchable      *Searchable
	SearchSelection *SearchSelection
	Settings        *Settings
}

func NewTable(cfg *config.Config, w *worker.Worker, cb *command.CommandBus, eb *event.EventBus) *Table {
	t := &Table{
		worker:     w,
		Sheet:      newSheet(eb, getShownHeader(cfg)),
		Searchable: newSearchable(getShownHeader(cfg)),
		Searching:  newSearching(ColumnAll, cb),
		Sorting:    newSorting(cb, cfg.UI.TableSortBy, cfg.UI.TableAscending),
		Settings:   newSettings(eb, cfg),
		Selected:   newSelected(eb, cb),
		SearchSelection: newSearchSelection(eb, cb),
	}
	setupCommands(t, eb, cb)
	t.Searchable.onChangedSearchBy = t.Searching.setSearchColumn
	eb.Subscribe(event.StoredSnapshot{}, func(v event.Event) {
		e := v.(event.StoredSnapshot)
		if e.Failed{
			return
		}
		t.Sorting.Sort()
	})
	return t
}

func setupCommands(t *Table, eb *event.EventBus, cb *command.CommandBus) {

	// CommandCellSelect
	cb.Register(command.CellSelect{}, func(v command.Command) error{

		e := v.(command.CellSelect)

		println("command.cell_select:", e.Version)

		pp := snapshot.Current.Load()
		if pp.Version() != e.Version {
			log.Printf("Warning: snapshot select: de-synced versions: %d != %d", pp.Version(), e.Version)
			return nil
		}

		if !e.Has {
			eb.Notify(event.CellSelected{
				Has: e.Has,
			})
			return nil
		}

		p := models.Cell{
			ID:  e.Point.ID,
			Column: e.Point.Column,
		}

		eb.Notify(event.CellSelected {
			Has: e.Has,
			Point: p,
			Version: e.Version,
		})
		return nil
	})

	// CommandSearch
	cb.Register(command.TableSearch{}, func(v command.Command) error {
		e := v.(command.TableSearch)
		pattern := e.Pattern
		column := e.Column
		t.worker.Jobs <- NewJobSearchTable(t.worker, pattern, column) 
		return nil
	})

	// CommandSort
	cb.Register(command.TableSort{}, func(v command.Command) error {
		e := v.(command.TableSort)
		column := e.Column
		asc := e.Asc
		t.worker.Jobs <- NewJobSortTable(t.worker, column, asc)
		return nil
	})
}

//
// Sheet
//

// Sheet a refrence view for table. 
// Methods should only be called by UI thread.
type Sheet struct {
	ssVersion       int64
	header          []string
	sorted          []int64
	idToRow         map[int64]int
	OnHeaderChanged func()
	OnSorted        func()
}

func newSheet(eb *event.EventBus, header []string) *Sheet {
	s := &Sheet{
		header: header,
		sorted: make([]int64, 0),
		idToRow: make(map[int64]int),
		OnSorted: func() {},
		OnHeaderChanged: func() {},
	}
	eb.Subscribe(event.HiddenColumn{}, func(v event.Event) {
		e := v.(event.HiddenColumn)
		header := make([]string, 0)
		for i, label := range models.BookEntryFields() {
			if !e.Hidden[i] {
				header = append(header, label)
			}
		}
		s.header = header
		s.OnHeaderChanged()
	})
	eb.Subscribe(event.TableSorted{}, func(v event.Event){
		ss := snapshot.Current.Load()
		e := v.(event.TableSorted)
		if ss.Version() == e.Version {
			s.ssVersion = e.Version
			s.sorted = e.Sorted
			clear(s.idToRow)
			for row, id := range s.sorted {
				s.idToRow[id] = row
			}
			s.OnSorted()
		} else {
			log.Printf("sheet.table_sorted: de-syned versions: %d != %d", ss.Version(), e.Version)
		}
	})
	return s
}

func (s *Sheet) RowToID(row int) (int64, error) {
	if len(s.sorted) <= row || row < 0 {
		return 0, fmt.Errorf("row %d to id: index out of range", row)
	}
	id := s.sorted[row]
	return id, nil
}

func (s *Sheet) IDToRow(id int64) (int, error) {
	row, ok := s.idToRow[id]
	if !ok {
		return 0, fmt.Errorf("id %d to row: id not found", id)
	}
	return row, nil
}

func (s *Sheet) Get(c models.Cell) (string, error) {
	ss := snapshot.Current.Load()
	return ss.Get(c)
}

func (s *Sheet) CordsToCell(row, col int) (models.Cell, error) {
	id, err := s.RowToID(row)
	if err != nil {
		return models.Cell{}, err
	}
	coll := s.Header()[col]
	return models.Cell{
		ID: id,
		Column: coll,
	}, nil
}

func (s *Sheet) PointToCords(p models.Cell) (row, col int, err error) {
	row, err = s.IDToRow(p.ID)
	if err != nil {
		return 0, 0, fmt.Errorf("point %d to cord: %w", p.ID, err)
	}
	col = slices.Index(s.Header(), p.Column)
	if col == -1 {
		return 0, 0, fmt.Errorf("point %s to cord: invalid column label", p.Column)
	}
	return row, col, nil
}

func (s *Sheet) Size() (rows, cols int) {
	ss := snapshot.Current.Load()
	rows, _  = ss.Size()
	cols = len(s.header)
	return
}

func (s *Sheet) Header() []string {
	return s.header
}


//
// Sorting
//

type Sorting struct {
	cb        *command.CommandBus
	Column    string
	Ascending bool
}

func newSorting(cb *command.CommandBus, column string, asc bool) *Sorting {
	s := &Sorting{
		cb: cb,
		Column: column,
		Ascending: asc,
	}
	return s
}

func (s *Sorting) Sort() {
	s.cb.Dispatch(command.TableSort{
		Asc: s.Ascending,
		Column: s.Column,
	})
}

//
// Searchable
//

type Searchable struct {
	headers  []string
	OnChangedOptions  func()
	onChangedSearchBy func(s string)
}

func newSearchable(headers []string) *Searchable {
	s := &Searchable{
		headers: headers,
		OnChangedOptions: func() {},
		onChangedSearchBy: func(_ string) {},

	}
	return s
}

func (s *Searchable) SetSearchBy(h string) {
	s.onChangedSearchBy(h)
}

func (s *Searchable) SetSelectable(headers []string) {
	s.headers = headers
	s.OnChangedOptions()
}

func (s *Searchable) Options() []string {
	o := []string{
		ColumnAll,
	}
	o = append(o, s.headers...)
	return o
}

//
// Searching
//

type Searching struct {
	cb       *command.CommandBus
	debounce func(func())
	column   string
}

func newSearching(column string, cb *command.CommandBus) *Searching {
	sr := &Searching{
		column: column,
		cb: cb,
		debounce: worker.Debounce(time.Millisecond * 300),
	}
	return sr
}

func (s *Searching) setSearchColumn(h string) {
	s.column = h
}

func (s *Searching) Search(pattern string) {
	s.debounce(func() {
		s.cb.Dispatch(command.TableSearch{Pattern: pattern, Column: s.column})
	})
}

//
// Selected
//

// Selected 
type Selected struct {
	cb         *command.CommandBus
	selected   models.Cell
	has        bool
	OnSelected func(models.Cell, bool)
}


func newSelected(eb *event.EventBus, cb *command.CommandBus) *Selected {
	es := &Selected{
		cb: cb,
		OnSelected: func(_ models.Cell, _ bool) {},
	}
	eb.Subscribe(event.CellSelected{}, func(v event.Event) {
		e := v.(event.CellSelected)
		es.selected = e.Point
		es.has = e.Has
		es.OnSelected(e.Point, e.Has)
	})
	return es
}

func (s *Selected) Get() (cell models.Cell, has bool) {
	has = s.has
	cell = s.selected
	return
}

func (es *Selected) Set(version int64, p models.Cell, ok bool) {
	c := models.Cell{
		Column: p.Column,
		ID: p.ID,
	}
	es.cb.Dispatch(command.CellSelect{Point: c, Has: ok, Version: version})
}

//
// Search Selection
//

type SearchSelection struct {
	eb        *event.EventBus
	cb        *command.CommandBus
	ssVersion int64
	selection []models.Cell
	position  int
}

func newSearchSelection(eb *event.EventBus, cb *command.CommandBus) *SearchSelection {
	sc := &SearchSelection{
		eb: eb,
		cb: cb,
		position: -1,
	}
	eb.Subscribe(event.TableSearched{}, func(v event.Event){
		e := v.(event.TableSearched)
		sc.ssVersion = e.Version
		sc.selection = e.Points
		sc.position = 0
		println("event.table_search: listener:", e.Version)

		if len(e.Points) != 0 && e.Pattern != "" {
			sc.selected()
		} else {
			sc.unselected()
		}
	})
	return sc
}

func (es *SearchSelection) Next() {
	if len(es.selection) == 0 {
		return
	}
	es.position += 1
	if es.position >= len(es.selection) {
		es.position = 0
	}
	es.selected()
}

func (es *SearchSelection) Prev() {
	if len(es.selection) == 0 {
		return
	}
	es.position -= 1
	if es.position < 0 {
		es.position = len(es.selection)-1
	}
	es.selected()
}

func (es *SearchSelection) selected() {
	p := es.selection[es.position]
	println("dispatch.cell_search:", es.ssVersion)
	es.cb.Dispatch(command.CellSelect{
		
		Version: es.ssVersion,
		Point: p,
		Has: true,
	})
}

func (es *SearchSelection) unselected() {
	es.selection = es.selection[:0]
	es.position = -1
	es.cb.Dispatch(command.CellSelect{
		Has: false,
	})
}


//
// Settings
//

type Settings struct {
	eb       *event.EventBus
	cfg      *config.Config

	OnColumnHidden func()
}

func newSettings(eb *event.EventBus, cfg *config.Config) *Settings {
	cs := &Settings{
		cfg: cfg,
		eb: eb,
	}
	cs.eb.Subscribe(event.TableSorted{}, func(v event.Event) {
		e := v.(event.TableSorted)
		cfg.UI.TableSortBy = e.Column
		cfg.UI.TableAscending = e.Asc
	})
	cs.eb.Subscribe(event.HiddenColumn{}, func(_ event.Event) {
		cs.OnColumnHidden()
	})
	return cs
}

func (ts *Settings) HeaderMinWidth() float32 {
	return ts.cfg.UI.TableMinWidth
}

func (ts *Settings) HeaderHeight() float32 {
	return ts.cfg.UI.TableHeaderHeight
}

func (ts *Settings) HeaderGetWidth(label string) float32 {
	idx := slices.Index(models.BookEntryFields(), label)
	return ts.cfg.UI.Headers[idx].Width
}

func (ts *Settings) HeaderSetWidth(label string, width float32) {
	idx := slices.Index(models.BookEntryFields(), label)
	h := ts.cfg.UI.Headers[idx]
	h.Width = width
	ts.cfg.UI.Headers[idx] = h
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
	ts.notifyHidden()
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
	ts.notifyHidden()
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
	ts.notifyHidden()
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

func (ts *Settings) notifyHidden() {
	hidden := make([]bool, len(ts.cfg.UI.Headers))
	for idx := range models.BookEntryFields() {
		hidden[idx] = ts.cfg.UI.Headers[idx].IsHidden
	}
	ts.eb.Notify(event.HiddenColumn{
		Hidden: hidden,
	})
}

//
// Functions and Helpers
//

func getSnapshotTraverser(ss *snapshot.Snapshot, by string) (search.Traverser, error) {
	var trv search.Traverser
	if by == ColumnAll {
		trv = newTableTraverse(ss)
	} else {
		idx := slices.Index(models.BookEntryFields(), by)
		if idx == -1 {
			return nil, fmt.Errorf("search %s: invalid column label", by)
		}
		trv = newColumnTraverse(ss, idx)
	}
	return trv, nil
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

func getShownHeader(cfg *config.Config) []string {	
	keys := make([]int, 0)
	for k, h := range cfg.UI.Headers {
		if !h.IsHidden {
			keys = append(keys, k)
		}
	}
	slices.Sort(keys)
	set := make([]string, len(keys))
	for i, k := range keys {
		set[i] = cfg.UI.Headers[k].Name
	}
	return set
}

func snapshotSort(ss *snapshot.Snapshot, column string, asc bool) ([]int64, error) {
	idx := slices.Index(models.BookEntryFields(), column)
	if idx == -1 {
		return nil, fmt.Errorf("sort %s: invalid column", column)
	}
	comp, err := app.CompareBookEntry(idx, asc)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, ss.Length())
	for i := range ss.Length() {
		id, _ := ss.RowToID(i)
		ids[i] = id
	}
	slices.SortFunc(ids, func(a, b int64) int {
		ba, _ := ss.GetBookEntryByID(a)
		bb, _ := ss.GetBookEntryByID(b)
		return comp(*ba, *bb)
	})
	return  ids, nil
}
